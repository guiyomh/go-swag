// Package swagger contains all the mecanisms to generate a swagger description
package swagger

import (
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"time"

	"github.com/fatih/structtag"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/guiyomh/swagger/pkg/router"
)

var (
	fixPathRe = regexp.MustCompile(`/:(\w+)`)
)

// tag attribute
const (
	DEFAULT     = "default"
	VALIDATE    = "validate"
	DESCRIPTION = "description"
	EMBED       = "embed"
	EXAMPLE     = "example"
)

// binding attributes
const (
	QUERY    = "query"
	FORM     = "form"
	URI      = "uri"
	HEADER   = "header"
	COOKIE   = "cookie"
	JSON     = "json"
	REQUIRED = "required"
)

const (
	BITSIZE = 64
	BASEINT = 10
)

type Swagger struct {
	Title           string
	Description     string
	Version         string
	DocsURL         string
	RedocURL        string
	OpenAPIURL      string
	Routers         []*router.Router
	Servers         openapi3.Servers
	TermsOfService  string
	Contact         *openapi3.Contact
	License         *openapi3.License
	OpenAPI         *openapi3.T
	SwaggerOptions  map[string]interface{}
	RedocOptions    map[string]interface{}
	validateOptions []validateOption
}

func New(title, description, version string, routers []*router.Router) *Swagger {
	//nolint:exhaustruct,nolintlint
	swagger := &Swagger{
		Title:       title,
		Description: description,
		Version:     version,
		DocsURL:     "/docs",
		RedocURL:    "/redoc",
		OpenAPIURL:  "/openapi.json",
		Routers:     routers,
		validateOptions: []validateOption{
			validateLenOption,
			validateEnumOption,
			validateMaxOption,
			validateMinOption,
		},
	}
	swagger.buildOpenAPI()

	return swagger
}

func (swagger *Swagger) buildOpenAPI() {

	components := openapi3.NewComponents()
	components.SecuritySchemes = openapi3.SecuritySchemes{}
	//nolint:exhaustruct,nolintlint
	swagger.OpenAPI = &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:          swagger.Title,
			Description:    swagger.Description,
			TermsOfService: swagger.TermsOfService,
			Contact:        swagger.Contact,
			License:        swagger.License,
			Version:        swagger.Version,
		},
		Servers:    swagger.Servers,
		Components: components,
	}
	swagger.OpenAPI.Paths = swagger.paths()
}

func (swagger *Swagger) paths() openapi3.Paths {
	paths := make(openapi3.Paths)
	var ok bool
	for _, router := range swagger.Routers {
		path := swagger.sanitizePath(router.Path)
		if _, ok = paths[path]; !ok {
			paths[path] = &openapi3.PathItem{} //nolint:exhaustruct,nolintlint
		}
		parameters, err := swagger.parametersFromModel(router.Model)
		if err != nil {
			continue
		}
		//nolint:exhaustruct,nolintlint
		operation := &openapi3.Operation{
			Tags:        router.Tags,
			OperationID: router.OperationID,
			Summary:     router.Summary,
			Description: router.Description,
			Deprecated:  router.Deprecated,
			Responses:   swagger.responses(router.Responses, router.ResponseContentType),
			Parameters:  parameters,
		}
		swagger.addPath(paths, router.Method, path, operation)
	}

	return paths
}

func (swagger *Swagger) addPath(paths openapi3.Paths, method, path string, operation *openapi3.Operation) {
	switch method {
	case http.MethodGet:
		paths[path].Get = operation
	case http.MethodPost:
		paths[path].Post = operation
	case http.MethodDelete:
		paths[path].Delete = operation
	case http.MethodPut:
		paths[path].Put = operation
	case http.MethodPatch:
		paths[path].Patch = operation
	case http.MethodHead:
		paths[path].Head = operation
	case http.MethodOptions:
		paths[path].Options = operation
	case http.MethodConnect:
		paths[path].Connect = operation
	case http.MethodTrace:
		paths[path].Trace = operation
	}
}

func (swagger *Swagger) parametersFromModel(model interface{}) (openapi3.Parameters, error) {
	parameters := openapi3.NewParameters()
	if model == nil {
		return parameters, nil
	}
	modelType, modelValue := swagger.typeAndValue(model)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		value := modelValue.Field(i)
		tags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			return openapi3.Parameters{}, err
		}
		_, err = tags.Get(EMBED)
		if err == nil {
			if embedParameters, err := swagger.parametersFromModel(value.Interface()); err != nil {
				parameters = append(parameters, embedParameters...)
			}
		}
		//nolint:exhaustruct,nolintlint
		parameter := &openapi3.Parameter{
			Schema: openapi3.NewSchemaRef("", swagger.schemaFromType(value.Interface())),
		}

		if params, err := swagger.parseQueryFromTags(tags, parameter, value, parameters); err == nil {
			parameters = params
		}
	}

	return parameters, nil
}

func (swagger *Swagger) parseQueryFromTags(
	tags *structtag.Tags,
	parameter *openapi3.Parameter,
	value reflect.Value,
	parameters openapi3.Parameters,
) (openapi3.Parameters, error) {
	parseTagQuery(tags, parameter)
	parseTagURI(tags, parameter)
	parseTagHeader(tags, parameter)
	parseTagCookie(tags, parameter)

	if parameter.In == "" {
		return openapi3.Parameters{}, ErrNoInParameter
	}
	parseTagDescription(tags, parameter)
	validateTag, err := tags.Get(VALIDATE)
	if err == nil {
		parameter.WithRequired(validateTag.Name == REQUIRED)
		options := validateTag.Options
		if len(options) > 0 {
			if schema, err := swagger.validateSchema(value.Interface(), options); err == nil {
				parameter.Schema = schema
			}
		}
	}
	defaultTag, err := tags.Get(DEFAULT)
	if err == nil {
		parameter.Schema.Value.WithDefault(defaultTag.Name)
	}
	exampleTag, err := tags.Get(EXAMPLE)
	if err == nil {
		parameter.Schema.Value.Example = exampleTag.Name
	}

	//nolint:nolintlint,exhaustruct
	parameters = append(parameters, &openapi3.ParameterRef{
		Value: parameter,
	})

	return parameters, nil
}

func (*Swagger) typeAndValue(model interface{}) (reflect.Type, reflect.Value) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	return modelType, modelValue
}

func (swagger *Swagger) validateSchema(value interface{}, options []string) (*openapi3.SchemaRef, error) {
	schema := openapi3.NewSchemaRef("", swagger.schemaFromType(value))
	for _, option := range options {

		for _, validateFunc := range swagger.validateOptions {
			if err := validateFunc(schema, option); err != nil {
				return nil, err
			}
		}
	}

	return schema, nil
}

func (swagger *Swagger) responses(responses map[string]*router.Response, contentType string) openapi3.Responses {
	resp := make(openapi3.Responses)
	for statusCode, response := range responses {
		schema := swagger.schemaFromModel(response.Model)
		description := response.Description
		//nolint:exhaustruct,nolintlint
		resp[statusCode] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: &description,
				Headers:     response.Headers,
				Content:     openapi3.NewContentWithSchema(schema, []string{contentType}),
			},
		}
	}

	return resp
}

func (swagger *Swagger) schemaFromType(model any) *openapi3.Schema {
	var schema *openapi3.Schema
	var min float64 = 0

	switch model.(type) {
	case int, int8, int16:
		schema = openapi3.NewIntegerSchema()
	case uint, uint8, uint16:
		schema = openapi3.NewIntegerSchema()
		schema.Min = &min
	case int32:
		schema = openapi3.NewInt32Schema()
	case uint32:
		schema = openapi3.NewInt32Schema()
		schema.Min = &min
	case int64:
		schema = openapi3.NewInt64Schema()
	case uint64:
		schema = openapi3.NewInt64Schema()
		schema.Min = &min
	case string:
		schema = openapi3.NewStringSchema()
	case time.Time:
		schema = openapi3.NewDateTimeSchema()
	case float32, float64:
		schema = openapi3.NewFloat64Schema()
	case bool:
		schema = openapi3.NewBoolSchema()
	case []byte:
		schema = openapi3.NewBytesSchema()
	case *multipart.FileHeader:
		schema = openapi3.NewStringSchema()
		schema.Format = "binary"
	case []*multipart.FileHeader:
		schema = openapi3.NewArraySchema()
		//nolint:exhaustruct,nolintlint
		schema.Items = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:   "string",
				Format: "binary",
			},
		}
	default:
		schema = swagger.schemaFromModel(model)
	}

	return schema
}

func (swagger *Swagger) schemaFromModel(model any) *openapi3.Schema {
	schema := openapi3.NewObjectSchema()
	if model == nil {
		return schema
	}

	modelType, modelValue := swagger.typeAndValue(model)

	//nolint:exhaustive,nolintlint
	switch modelType.Kind() {
	case reflect.Struct:
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			value := modelValue.Field(i)
			if err := swagger.schemaFromReflectStruct(value, field, schema); err != nil {
				continue
			}
		}
	case reflect.Slice:
		schema = openapi3.NewArraySchema()

		//nolint:exhaustruct,nolintlint
		schema.Items = &openapi3.SchemaRef{
			Value: swagger.schemaFromModel(reflect.New(modelType.Elem()).Elem().Interface()),
		}
	case reflect.Map:
		schema = openapi3.NewObjectSchema()
	default:
		schema = swagger.schemaFromType(model)
	}

	return schema
}

func (swagger *Swagger) schemaFromReflectStruct(
	value reflect.Value,
	field reflect.StructField,
	schema *openapi3.Schema,
) error {
	fieldSchema := swagger.schemaFromType(value.Interface())
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return err
	}
	_, err = tags.Get(EMBED)
	if err == nil {
		embedSchema := swagger.schemaFromModel(value.Interface())
		for key, embedProperty := range embedSchema.Properties {
			schema.Properties[key] = embedProperty
		}
		schema.Required = append(schema.Required, embedSchema.Required...)
	}
	tag, err := tags.Get(JSON)
	if err != nil {
		return err
	}
	parseTags(tag.Name, tags, schema, fieldSchema)
	schema.Properties[tag.Name] = openapi3.NewSchemaRef("", fieldSchema)

	return nil
}

func (swagger *Swagger) sanitizePath(path string) string {
	return fixPathRe.ReplaceAllString(path, "/{${1}}")
}

func (swagger *Swagger) MarshalJSON() ([]byte, error) {
	return swagger.OpenAPI.MarshalJSON()
}
