package stopwatch

import "errors"

var (
	ErrAlreadyRunning = errors.New("stopwatch is already running")
	ErrNotRunning     = errors.New("stopwatch is not running")
)
