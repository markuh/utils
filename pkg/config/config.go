package config

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config interface {
	Validate() error
}

type API struct {
	Address         string        `envconfig:"ADDRESS" default:"localhost:8080"`
	Timeout         time.Duration `envconfig:"TIMEOUT" default:"10s"`
	CORSOrigins     []string      `envconfig:"CORS_ORIGINS" default:"*"`
	LoggerSkipPaths []string      `envconfig:"LOGGER_SKIP_PATHS" default:"/health"`
}

type Postgres struct {
	Host              string        `envconfig:"HOST"`
	Port              string        `envconfig:"PORT"`
	User              string        `envconfig:"USER"`
	Password          string        `envconfig:"PASSWORD"`
	Name              string        `envconfig:"NAME"`
	Schema            string        `envconfig:"SCHEMA"`
	SSLMode           string        `envconfig:"SSL_MODE" default:"require"`
	MaxConns          int           `envconfig:"MAX_CONNS" default:"5"`
	MinConns          int           `envconfig:"MIN_CONNS" default:"0"`
	MaxConnLifetime   time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"30m"`
	MaxConnIdleTime   time.Duration `envconfig:"CONN_MAX_IDLE_TIME" default:"5m"`
	ConnectTimeout    time.Duration `envconfig:"CONNECT_TIMEOUT" default:"5s"`
	HealthCheckPeriod time.Duration `envconfig:"HEALTH_CHECK_PERIOD" default:"30s"`
}

func (p Postgres) Validate() error {
	if p.MaxConns < 1 {
		return fmt.Errorf("postgres MaxConns must be >= 1, got %d", p.MaxConns)
	}
	if p.MinConns > p.MaxConns {
		return fmt.Errorf("postgres MinConns (%d) must be <= MaxConns (%d)", p.MinConns, p.MaxConns)
	}
	if p.MaxConns > 20 {
		log.Printf("Warning: postgres MaxConns (%d) exceeds recommended maximum of 20 for PgBouncer", p.MaxConns)
	}
	return nil
}

type Log struct {
	Level            string `envconfig:"LOG_LEVEL" default:"info"`
	DebugLogFilePath string `envconfig:"LOG_DEBUG_FILE_PATH" default:"logs/debug.log"`
}

func Read(c Config) error {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	if err := envconfig.Process("", &c); err != nil {
		return err
	}
	if err := c.Validate(); err != nil {
		return err
	}
	return nil
}
