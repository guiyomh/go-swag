package router

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRouter_parseRouter(t *testing.T) {

	tests := []struct {
		comment    string
		wantMethod string
		wantPath   string
	}{
		{
			comment:    "/bar [get]",
			wantMethod: "GET",
			wantPath:   "/bar",
		},
		{
			comment:    "/bar post",
			wantMethod: "POST",
			wantPath:   "/bar",
		},
		{
			comment:    "/bar/foo/baz GET",
			wantMethod: "GET",
			wantPath:   "/bar/foo/baz",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Should parse router comment > %s ", tt.comment), func(t *testing.T) {
			rte := &Router{}
			err := rte.parseRouter(tt.comment)
			require.NoError(t, err)
			require.Equal(t, tt.wantMethod, rte.Method)
			require.Equal(t, tt.wantPath, rte.Path)
		})
	}
}

func TestRouter_parseDescription(t *testing.T) {

	tests := []struct {
		rte     *Router
		comment string
		want    string
	}{
		{
			rte:     &Router{},
			comment: "biloute bar baz",
			want:    "biloute bar baz",
		},
		{
			rte:     &Router{Description: "foo"},
			comment: "biloute bar baz",
			want:    "foo\nbiloute bar baz",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Should parse description > %s ", tt.comment), func(t *testing.T) {
			tt.rte.parseDescription(tt.comment)
			require.Equal(t, tt.want, tt.rte.Description)
		})
	}
}

func TestRouter_parseResponse(t *testing.T) {
	tests := []struct {
		comment          string
		wantStatusCode   string
		wantResponseDesc string
	}{
		{
			comment:          "200 nil my description",
			wantStatusCode:   "200",
			wantResponseDesc: "my description",
		},
		{
			comment:          "404 nil not found resut",
			wantStatusCode:   "404",
			wantResponseDesc: "not found resut",
		},
		{
			comment:          "500 nil",
			wantStatusCode:   "500",
			wantResponseDesc: "",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Should parse response > %s ", tt.comment), func(t *testing.T) {
			rte := &Router{}
			err := rte.parseResponse(tt.comment)
			require.NoError(t, err)
			if response, ok := rte.Responses[tt.wantStatusCode]; ok {
				require.Equal(t, tt.wantResponseDesc, response.Description)

				return
			}
			t.Error(fmt.Sprintf("StatusCode '%s' doesn't exist in response map", tt.wantStatusCode))
		})
	}
}
