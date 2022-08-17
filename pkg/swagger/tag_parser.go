package swagger

import (
	"github.com/fatih/structtag"
	"github.com/getkin/kin-openapi/openapi3"
)

func parseParameter(tags *structtag.Tags, parameter *openapi3.Parameter, tagName, in string) {
	tag, err := tags.Get(tagName)
	if err == nil {
		parameter.In = in
		parameter.Name = tag.Name
	}
}

func parseTagQuery(tags *structtag.Tags, parameter *openapi3.Parameter) {
	parseParameter(tags, parameter, QUERY, openapi3.ParameterInQuery)
}

func parseTagURI(tags *structtag.Tags, parameter *openapi3.Parameter) {
	parseParameter(tags, parameter, URI, openapi3.ParameterInPath)
}

func parseTagHeader(tags *structtag.Tags, parameter *openapi3.Parameter) {
	parseParameter(tags, parameter, HEADER, openapi3.ParameterInHeader)
}

func parseTagCookie(tags *structtag.Tags, parameter *openapi3.Parameter) {
	parseParameter(tags, parameter, COOKIE, openapi3.ParameterInCookie)
}

func parseTagDescription(tags *structtag.Tags, parameter *openapi3.Parameter) {
	tag, err := tags.Get(DESCRIPTION)
	if err == nil {
		parameter.WithDescription(tag.Name)
	}
}

func parseTags(tagName string, tags *structtag.Tags, schema, fieldSchema *openapi3.Schema) {
	validateTag, err := tags.Get(VALIDATE)
	if err == nil && validateTag.Name == REQUIRED {
		schema.Required = append(schema.Required, tagName)
	}
	descriptionTag, err := tags.Get(DESCRIPTION)
	if err == nil {
		fieldSchema.Description = descriptionTag.Name
	}
	defaultTag, err := tags.Get(DEFAULT)
	if err == nil {
		fieldSchema.Default = defaultTag.Name
	}
	exampleTag, err := tags.Get(EXAMPLE)
	if err == nil {
		fieldSchema.Example = exampleTag.Name
	}
}
