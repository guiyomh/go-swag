package router

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func FromComments(directory string) ([]*Router, error) {
	routers := make([]*Router, 0)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, directory, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, astpkg := range pkgs {
		for file, astfile := range astpkg.Files {
			if rts, err := fromASTFile(astfile); err == nil {
				routers = append(routers, rts...)
			} else {
				return nil, errors.Wrap(ErrParseFileComment, file)
			}
		}
	}

	return routers, nil
}

func fromASTFile(astfile *ast.File) ([]*Router, error) {
	routers := make([]*Router, 0)
	for _, group := range astfile.Comments {
		router, err := fromASTCommentGroup(group)
		if err != nil {
			return nil, err
		}
		routers = append(routers, router)
	}

	return routers, nil
}

type Attribute string

const (
	DescriptionAttr             Attribute = "@description"
	SummaryAttr                 Attribute = "@summary"
	RouterAttr                  Attribute = "@router"
	ResponseAttr                Attribute = "@response"
	nbPartOfRouter                        = 2
	nbPartOfResponse                      = 3
	nbPartOfResponseWithoutDesc           = 2
)

func fromASTCommentGroup(group *ast.CommentGroup) (*Router, error) {
	rte := &Router{
		Responses: make(ResponseMap),
	}
	for _, comment := range group.List {
		rte.ParseComments(comment.Text)
	}

	return rte, nil
}

func (router *Router) ParseComments(comment string) error {

	commentLine := strings.TrimSpace(strings.TrimLeft(comment, "/"))
	if len(commentLine) == 0 {
		return nil
	}

	attribute := Attribute(strings.ToLower(strings.Fields(commentLine)[0]))
	lineRemainder := strings.TrimSpace(commentLine[len(attribute):])

	switch attribute {
	case DescriptionAttr:
		router.parseDescription(lineRemainder)
	case SummaryAttr:
		router.Summary = lineRemainder
	case ResponseAttr:
		return router.parseResponse(lineRemainder)
	case RouterAttr:
		return router.parseRouter(lineRemainder)
	}

	return nil
}

func (router *Router) parseDescription(comment string) {
	if router.Description == "" {
		router.Description = comment

		return
	}

	router.Description += "\n" + comment
}

func (router *Router) parseRouter(comment string) error {
	fields := strings.Fields(comment)
	if len(fields) != nbPartOfRouter {
		return ErrParseRouterComment
	}

	router.Path = fields[0]
	router.Method = strings.ToUpper(strings.Trim(fields[1], "[]"))

	return nil
}

func (router *Router) parseResponse(comment string) error {
	fields := strings.SplitN(comment, " ", nbPartOfResponse)

	if len(fields) < nbPartOfResponseWithoutDesc {
		return ErrParseResponseComment
	}

	if router.Responses == nil {
		router.Responses = make(ResponseMap)
	}

	response := &Response{}

	if len(fields) == nbPartOfResponse {
		response.Description = fields[2]
	}

	router.Responses[fields[0]] = response

	return nil
}

// https://stackoverflow.com/questions/23030884/is-there-a-way-to-create-an-instance-of-a-struct-from-a-string
var typeRegistry = make(map[string]reflect.Type)
