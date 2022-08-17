package router

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ResponseMap map[string]*Response

type Response struct {
	Description string
	Model       interface{}
	Headers     openapi3.Headers
}
