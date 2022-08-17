// Package router contains the repreentation of a request routing and its response
package router

type Handler interface{}

const (
	MIMEApplicationJSON = "application/json"
)

type Router struct {
	Path                string
	Method              string
	Handler             Handler
	Summary             string
	Description         string
	Deprecated          bool
	RequestContentType  string
	ResponseContentType string
	Tags                []string
	Model               any
	OperationID         string
	Responses           map[string]*Response
}

func New(path, method string, handler Handler, options ...Option) *Router {
	//nolint:nolintlint,exhaustruct
	router := &Router{
		Path:                path,
		Method:              method,
		Handler:             handler,
		RequestContentType:  MIMEApplicationJSON,
		ResponseContentType: MIMEApplicationJSON,
	}

	for _, opt := range options {
		opt(router)
	}

	return router
}
