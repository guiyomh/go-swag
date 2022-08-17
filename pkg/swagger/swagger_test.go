//nolint:exhaustruct, nolintlint
package swagger

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/guiyomh/swagger/pkg/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	swag := New("foo", "my description", "2.2", nil)

	assert.Equal(t, "foo", swag.Title)
	assert.Equal(t, "my description", swag.Description)
	assert.Equal(t, "2.2", swag.Version)
	assert.Len(t, swag.Routers, 0)
}

func TestSwagger_buildOpenApi(t *testing.T) {

	t.Run("Should build the open api without routes", func(t *testing.T) {
		swag := &Swagger{
			Title:       "foo",
			Description: "bar",
			Version:     "2.3",
		}
		swag.buildOpenAPI()
		assert.Equal(t, "foo", swag.OpenAPI.Info.Title)
		assert.Equal(t, "bar", swag.OpenAPI.Info.Description)
		assert.Equal(t, "2.3", swag.OpenAPI.Info.Version)
		assert.Len(t, swag.OpenAPI.Paths, 0)
	})

	t.Run("Should build the open api with routes", func(t *testing.T) {

		swag := &Swagger{
			Title:       "foo",
			Description: "bar",
			Version:     "2.3",
			Routers: []*router.Router{
				router.New("/", http.MethodGet, func() {}),
				router.New("/product", http.MethodDelete, func() {}),
			},
		}
		swag.buildOpenAPI()
		assert.Equal(t, "foo", swag.OpenAPI.Info.Title)
		assert.Equal(t, "bar", swag.OpenAPI.Info.Description)
		assert.Equal(t, "2.3", swag.OpenAPI.Info.Version)
		assert.Len(t, swag.OpenAPI.Paths, 2)
	})
}

func TestSwagger_path(t *testing.T) {
	t.Run("Should build simple GET route", func(t *testing.T) {
		swag := &Swagger{
			Title:       "foo",
			Description: "bar",
			Version:     "2.3",
			Routers: []*router.Router{
				router.New("/", http.MethodGet, func() {}),
			},
		}
		paths := swag.paths()
		require.Len(t, paths, 1)
		require.Nil(t, paths["/"].Post)
		require.Nil(t, paths["/"].Delete)
		require.Nil(t, paths["/"].Put)
		require.Nil(t, paths["/"].Patch)
		require.Nil(t, paths["/"].Head)
		require.Nil(t, paths["/"].Options)
		require.Nil(t, paths["/"].Connect)
		require.Nil(t, paths["/"].Trace)
		require.IsType(t, new(openapi3.Operation), paths["/"].Get)
	})
	t.Run("Should build GET route with responses", func(t *testing.T) {

		type ModelExample struct {
			Name string `json:"name" query:"name"`
			Age  int    `json:"age"`
		}
		type ModelError struct {
			Message string `json:"msg" `
			Code    int    `json:"code"`
		}

		swag := &Swagger{
			Title:       "foo",
			Description: "bar",
			Version:     "2.3",
			Routers: []*router.Router{
				{
					Path:                "/hello/:name",
					Method:              http.MethodGet,
					ResponseContentType: "txt/xml",
					Tags:                []string{"foo", "bar"},
					OperationID:         "foobar",
					Summary:             "foo bar summary",
					Description:         "foo bar description",
					Responses: router.ResponseMap{
						"200": {
							Description: "description success",
							Model:       ModelExample{},
						},
						"400": {
							Description: "description xml",
							Model:       ModelError{},
						},
					},
				},
			},
		}
		paths := swag.paths()
		require.Len(t, paths, 1)
		require.Nil(t, paths["/hello/{name}"].Post)
		require.Nil(t, paths["/hello/{name}"].Delete)
		require.Nil(t, paths["/hello/{name}"].Put)
		require.Nil(t, paths["/hello/{name}"].Patch)
		require.Nil(t, paths["/hello/{name}"].Head)
		require.Nil(t, paths["/hello/{name}"].Options)
		require.Nil(t, paths["/hello/{name}"].Connect)
		require.Nil(t, paths["/hello/{name}"].Trace)
		require.IsType(t, new(openapi3.Operation), paths["/hello/{name}"].Get)
		operation := paths["/hello/{name}"].Get
		require.Equal(t, "foobar", operation.OperationID)
		require.Equal(t, "foo bar summary", operation.Summary)
		require.Equal(t, "foo bar description", operation.Description)
		require.False(t, operation.Deprecated)
		require.Len(t, operation.Responses, 2)
		require.IsType(t, new(openapi3.ResponseRef), operation.Responses["200"])
		require.IsType(t, new(openapi3.ResponseRef), operation.Responses["400"])
		// TODO: Check each response
	})
}

func TestSwagger_schemaFromType(t *testing.T) {
	min := float64(0)
	tests := []struct {
		name       string
		input      any
		wantType   string
		wantFormat string
		wantMin    *float64
	}{
		{
			name:     "Should return a int",
			input:    8,
			wantType: openapi3.TypeInteger,
		},
		{
			name:     "Should return a int8",
			input:    int8(8),
			wantType: openapi3.TypeInteger,
		},
		{
			name:     "Should return a int16",
			input:    int16(8),
			wantType: openapi3.TypeInteger,
		},
		{
			name:       "Should return a int32",
			input:      int32(8),
			wantType:   openapi3.TypeInteger,
			wantFormat: "int32",
		},
		{
			name:       "Should return a int64",
			input:      int64(8),
			wantType:   openapi3.TypeInteger,
			wantFormat: "int64",
		},
		{
			name:     "Should return a uint",
			input:    uint(8),
			wantType: openapi3.TypeInteger,
			wantMin:  &min,
		},
		{
			name:     "Should return a uint8",
			input:    uint8(8),
			wantType: openapi3.TypeInteger,
			wantMin:  &min,
		},
		{
			name:     "Should creates a schema from uint16",
			input:    uint16(8),
			wantType: openapi3.TypeInteger,
			wantMin:  &min,
		},
		{
			name:       "Should creates a schema from uint32",
			input:      uint32(8),
			wantType:   openapi3.TypeInteger,
			wantFormat: "int32",
			wantMin:    &min,
		},
		{
			name:       "Should creates a schema from uint64",
			input:      uint64(8),
			wantType:   openapi3.TypeInteger,
			wantFormat: "int64",
			wantMin:    &min,
		},
		{
			name:       "Should creates a schema from float32",
			input:      float32(8),
			wantType:   openapi3.TypeNumber,
			wantFormat: "",
		},
		{
			name:       "Should creates a schema from float64",
			input:      float64(8),
			wantType:   openapi3.TypeNumber,
			wantFormat: "",
		},
		{
			name:       "Should creates a schema from bool",
			input:      true,
			wantType:   openapi3.TypeBoolean,
			wantFormat: "",
		},
		{
			name:       "Should creates a schema from string of byte",
			input:      []byte{'b', 'a', 'r'},
			wantType:   openapi3.TypeString,
			wantFormat: "byte",
		},
		{
			name:     "Should creates a schema from string",
			input:    "foo",
			wantType: openapi3.TypeString,
		},
	}

	swag := &Swagger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := swag.schemaFromType(tt.input)
			assert.Equal(t, tt.wantType, schema.Type)
			assert.Equal(t, tt.wantFormat, schema.Format)
			assert.Equal(t, tt.wantMin, schema.Min)
		})
	}
}

func TestSwagger_schemaFromModel(t *testing.T) {

	t.Run("Should create a schema from struct", func(t *testing.T) {
		swag := &Swagger{}

		type FakeModel struct {
			Name   string `json:"name" description:"desc of the name" example:"John Doe"`
			Age    uint   `json:"age" default:"21"`
			Active bool   `json:"active" validate:"required"`
		}

		schema := swag.schemaFromModel(new(FakeModel))

		require.Equal(t, openapi3.TypeObject, schema.Type)
		require.Len(t, schema.Properties, 3)
		require.Equal(t, openapi3.TypeString, schema.Properties["name"].Value.Type)
		require.Equal(t, "desc of the name", schema.Properties["name"].Value.Description)
		require.Equal(t, "John Doe", schema.Properties["name"].Value.Example)
		require.Equal(t, openapi3.TypeInteger, schema.Properties["age"].Value.Type)
		require.Equal(t, "21", schema.Properties["age"].Value.Default)
		require.Equal(t, openapi3.TypeBoolean, schema.Properties["active"].Value.Type)
		require.Equal(t, []string{"active"}, schema.Required)
	})

	t.Run("Should create a schema from slice", func(t *testing.T) {
		swag := &Swagger{}

		schema := swag.schemaFromModel([]string{"foo"})
		require.Equal(t, openapi3.TypeArray, schema.Type)
		require.Equal(t, openapi3.TypeString, schema.Items.Value.Type)
	})

	t.Run("Should create a schema from map", func(t *testing.T) {
		swag := &Swagger{}

		schema := swag.schemaFromModel(map[string]string{"foo": "bar"})
		require.Equal(t, openapi3.TypeObject, schema.Type)
	})
}

func TestSwagger_sanitizePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "/hello/:name",
			want:  "/hello/{name}",
		},
		{
			input: "/hello/:name/:id",
			want:  "/hello/{name}/{id}",
		},
		{
			input: "/hello/:name/welcome",
			want:  "/hello/{name}/welcome",
		},
		{
			input: "/:welcome/fr/:name",
			want:  "/{welcome}/fr/{name}",
		},
	}
	swag := &Swagger{}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Should transform %s to %s", tt.input, tt.want), func(t *testing.T) {
			got := swag.sanitizePath(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSwagger_parametersFromModel(t *testing.T) {
	swag := &Swagger{}

	type FakeModel struct {
		Name   string `cookie:"name"`
		Age    uint   `query:"qa" validate:"required"`
		Active bool   `uri:"active"`
		Token  string `header:"Authorization" description:"token used for the request"`
	}

	parameters, err := swag.parametersFromModel(new(FakeModel))
	require.NoError(t, err)
	require.Len(t, parameters, 4)
	assert.Equal(t, "name", parameters[0].Value.Name)
	assert.Equal(t, "cookie", parameters[0].Value.In)

	assert.Equal(t, "qa", parameters[1].Value.Name)
	assert.Equal(t, "query", parameters[1].Value.In)
	assert.True(t, parameters[1].Value.Required)

	assert.Equal(t, "active", parameters[2].Value.Name)
	assert.Equal(t, "path", parameters[2].Value.In)

	assert.Equal(t, "Authorization", parameters[3].Value.Name)
	assert.Equal(t, "header", parameters[3].Value.In)
	assert.Equal(t, "token used for the request", parameters[3].Value.Description)
}

func TestSwagger_validateSchema(t *testing.T) {
	swag := New("foo", "bar", "2.0.0", nil)

	t.Run("Should return a schema with min and max", func(t *testing.T) {
		schema, err := swag.validateSchema(int(8), []string{"min=1", "max=10"})
		require.NoError(t, err)
		assert.Equal(t, 1.0, *schema.Value.Min)
		assert.Equal(t, 10.0, *schema.Value.Max)
	})
	t.Run("Should return a schema with enum", func(t *testing.T) {
		schema, err := swag.validateSchema("red", []string{"enum=red,green,blue"})
		require.NoError(t, err)
		assert.Equal(t, []any{"red", "green", "blue"}, schema.Value.Enum)
	})
}
