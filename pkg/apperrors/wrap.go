package apperrors

import "fmt"

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

func Wrapf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}

	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	return fmt.Errorf("%s: %w", msg, err)
}

func New(msg string, args ...any) error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return fmt.Errorf(msg)
}
