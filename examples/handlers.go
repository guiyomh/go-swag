package examples

type ModelExample struct {
	Name string `query:"name" validate:"required" json:"name" description:"name of model" default:"test"`
	Age  int
}

// ExampleEndpoint foo
// @Summary foo route
// @Description return a basic message
// @Produce json
// @Response 200 ModelExample reponse description
// @Response 400 ModelExample bad response
// @Router / [get]
func Foo() {}

// ExampleEndpoint bar
// @Summary bar route
// @Description return a basic message
// @Produce json
// @Success 200 {object} ModelExample
// @Router /bar [get]
func Bar() {}
