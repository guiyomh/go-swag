package swagger

import (
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

// validate option
const (
	MinOption  = "min="
	MaxOption  = "max="
	LenOption  = "len="
	EnumOption = "enum="
)

type validateOption = func(*openapi3.SchemaRef, string) error

func validateLenOption(schema *openapi3.SchemaRef, option string) error {
	if strings.HasPrefix(option, LenOption) {
		value, err := strconv.ParseInt(option[len(LenOption):], BASEINT, BITSIZE)
		if err != nil {
			return errors.Wrap(err, ErrParseLenOption.Error())
		}
		schema.Value.WithLength(value)
	}

	return nil
}

func validateMaxOption(schema *openapi3.SchemaRef, option string) error {
	if strings.HasPrefix(option, MaxOption) {
		value, err := strconv.ParseFloat(option[len(MaxOption):], BITSIZE)
		if err != nil {
			return errors.Wrap(err, ErrParseMaxOption.Error())
		}
		schema.Value.WithMax(value)
	}

	return nil
}

func validateMinOption(schema *openapi3.SchemaRef, option string) error {
	if strings.HasPrefix(option, MinOption) {
		value, err := strconv.ParseFloat(option[len(MinOption):], BITSIZE)
		if err != nil {
			return errors.Wrap(err, ErrParseMinOption.Error())
		}
		schema.Value.WithMin(value)
	}

	return nil
}

func validateEnumOption(schema *openapi3.SchemaRef, option string) error {
	if strings.HasPrefix(option, EnumOption) {
		optionItems := strings.Split(option[len(EnumOption):], ",")
		enums := make([]interface{}, len(optionItems))
		for i, optionItem := range optionItems {
			enums[i] = optionItem
		}
		schema.Value.WithEnum(enums...)
	}

	return nil
}
