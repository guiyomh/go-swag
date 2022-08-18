package swagger_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guiyomh/swagger/pkg/router"
	"github.com/guiyomh/swagger/pkg/swagger"
)

func ExampleNew() {

	MyHandler := func() {
		// nothing for example
	}

	type MyResponse struct {
		Content string `json:"content"`
	}

	routers := []*router.Router{
		router.New(
			"/product/:id",
			http.MethodGet,
			MyHandler,
			router.Description("my description"),
			router.Responses(router.ResponseMap{
				"200": &router.Response{Description: "Success", Model: new(MyResponse)},
				"404": &router.Response{Description: "Not found", Model: new(MyResponse)},
			}),
		),
	}

	swag, err := swagger.New(
		"my api",
		"this this a example of swagger",
		"2.0.1",
		routers,
	)

	if err != nil {
		panic(err)
	}

	buf, err := json.MarshalIndent(swag, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))

	// Output:
	// {
	//   "components": {},
	//   "info": {
	//     "description": "this this a example of swagger",
	//     "title": "my api",
	//     "version": "2.0.1"
	//   },
	//   "openapi": "3.0.0",
	//   "paths": {
	//     "/product/{id}": {
	//       "get": {
	//         "description": "my description",
	//         "responses": {
	//           "200": {
	//             "content": {
	//               "application/json": {
	//                 "schema": {
	//                   "properties": {
	//                     "content": {
	//                       "type": "string"
	//                     }
	//                   },
	//                   "type": "object"
	//                 }
	//               }
	//             },
	//             "description": "Success"
	//           },
	//           "404": {
	//             "content": {
	//               "application/json": {
	//                 "schema": {
	//                   "properties": {
	//                     "content": {
	//                       "type": "string"
	//                     }
	//                   },
	//                   "type": "object"
	//                 }
	//               }
	//             },
	//             "description": "Not found"
	//           }
	//         }
	//       }
	//     }
	//   }
	// }
}
