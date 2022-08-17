package swagger

import (
	"github.com/pkg/errors"
)

var (
	ErrParseMaxOption  = errors.New("Cannot parse max option. the right syntaxe is validate:\"max=12\".")
	ErrParseMinOption  = errors.New("Cannot parse min option. the right syntaxe is validate:\"min=1\".")
	ErrParseLenOption  = errors.New("Cannot parse len option. the right syntaxe is validate:\"len=1\".")
	ErrParseEnumOption = errors.New("Cannot parse enum option. the right syntaxe is validate:\"enum=red,blue,green\".")
	ErrNoInParameter   = errors.New("No In parameters")
)
