package logger

import "github.com/pkg/errors"

var (
	ErrNilWriter = errors.New("the writer is nil")
)
