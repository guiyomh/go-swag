//nolint:exhaustruct,nolintlint
package router

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	rte := &Router{}

	Tags("foo", "bar")(rte)

	require.Len(t, rte.Tags, 2)
	require.Equal(t, []string{"foo", "bar"}, rte.Tags)
}

func TestSummary(t *testing.T) {
	rte := &Router{}

	Summary("foo bar")(rte)

	require.Equal(t, "foo bar", rte.Summary)
}

func TestDescription(t *testing.T) {
	rte := &Router{}

	Description("foo bar")(rte)

	require.Equal(t, "foo bar", rte.Description)
}

func TestDeprecated(t *testing.T) {
	rte := &Router{}

	require.False(t, rte.Deprecated)
	Deprecated()(rte)

	require.True(t, rte.Deprecated)
}

func TestOperationID(t *testing.T) {
	rte := &Router{}

	OperationID("biloute")(rte)

	require.Equal(t, "biloute", rte.OperationID)
}

func TestResponses(t *testing.T) {
	rte := &Router{}
	responses := map[string]*Response{
		"200": {
			Description: "fake response",
		},
	}

	Responses(responses)(rte)

	require.Len(t, rte.Responses, 1)
	require.Equal(t, "fake response", rte.Responses["200"].Description)
}

func TestModel(t *testing.T) {
	rte := &Router{}

	type FakeModel struct {
		Foo string
		Bar int
	}
	Model(new(FakeModel))(rte)

	require.IsType(t, new(FakeModel), rte.Model)
}
