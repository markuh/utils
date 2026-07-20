package app_runner

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"golang.org/x/sync/errgroup"
)

var (
	ErrInterruptedBySignal = errors.New("interrupted by signal")
)

// App interface for applications
type App interface {
	Start(ctx context.Context) error
	Shutdown() error
}

// AppsRunner interface for running applications
type AppsRunner interface {
	Register(name string, start starter, shutdown shutdowner)
	RegisterApp(name string, instance App)
	Run(ctx context.Context) error
}

// appInstance internal structure
type appInstance struct {
	Name     string
	Start    starter
	Shutdown shutdowner
}

type starter func(ctx context.Context) error
type shutdowner func() error

// appsRunner service for running applications in graceful shutdown mode
type appsRunner struct {
	apps   []appInstance
	logger *slog.Logger
}

// NewAppsRunner constructor for the service
func NewAppsRunner(logger *slog.Logger) AppsRunner {
	if logger == nil {
		logger = slog.Default()
	}

	return &appsRunner{
		apps:   make([]appInstance, 0),
		logger: logger,
	}
}

// Register function for registration
func (a *appsRunner) Register(name string, start starter, shutdown shutdowner) {
	a.apps = append(a.apps, appInstance{
		Name:     name,
		Start:    start,
		Shutdown: shutdown,
	})
}

// RegisterApp function for registration of both methods
func (a *appsRunner) RegisterApp(name string, instance App) {
	if instance == nil {
		return
	}
	a.apps = append(a.apps, appInstance{
		Name:     name,
		Start:    instance.Start,
		Shutdown: instance.Shutdown,
	})
}

func (a *appsRunner) Run(ctx context.Context) error {
	// create context with cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// start all applications in errgroup
	eg, ctx := errgroup.WithContext(ctx)

	// start all runners
	for _, app := range a.apps {
		a.logger.Debug("Run service", "name", app.Name)
		eg.Go(func() error {
			return app.Start(ctx)
		})
	}

	// graceful shutdown of all applications
	eg.Go(func() error {
		<-ctx.Done()

		var errorsList []error
		for _, app := range slices.Backward(a.apps) {
			a.logger.Debug("Stop service", "name", app.Name)
			if app.Shutdown == nil {
				a.logger.Debug("No shutdown function for service", "name", app.Name)
				continue
			}

			err := app.Shutdown()
			a.logger.Debug("Stop service result", "name", app.Name, "error", err)
			if err != nil {
				errorsList = append(errorsList, err)
				a.logger.Error("shutdown error", "name", app.Name, "error", err)
			}
		}

		return errors.Join(errorsList...)
	})

	eg.Go(func() error {
		sig := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
		shutdownCh := make(chan os.Signal, len(sig))
		signal.Notify(shutdownCh, sig...)
		defer signal.Stop(shutdownCh)

		select {
		case <-shutdownCh:
			cancel()
			return ErrInterruptedBySignal
		case <-ctx.Done():
			return nil
		}
	})

	if err := eg.Wait(); err != nil {
		if errors.Is(err, ErrInterruptedBySignal) {
			a.logger.Warn("Server shutting down", "error", err)
		} else {
			a.logger.Error("Server terminating with error", "error", err)
			return err
		}
	}

	a.logger.Info("Application was stopped")
	return nil
}
