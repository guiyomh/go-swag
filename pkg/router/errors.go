package router

import (
	"github.com/pkg/errors"
)

var (
	ErrParseFileComment     = errors.New("Cannot parse comment of file")
	ErrParseRouterComment   = errors.New("Cannot parse router comment")
	ErrParseResponseComment = errors.New("Cannot parse response comment")
)
