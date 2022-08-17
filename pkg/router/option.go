package router

type Option func(router *Router)

func Tags(tags ...string) Option {
	return func(router *Router) {
		if router.Tags == nil {
			router.Tags = tags
		} else {
			router.Tags = append(router.Tags, tags...)
		}
	}
}

func Summary(summary string) Option {
	return func(router *Router) {
		router.Summary = summary
	}
}

func Description(description string) Option {
	return func(router *Router) {
		router.Description = description
	}
}

// Deprecated mark api is deprecated
func Deprecated() Option {
	return func(router *Router) {
		router.Deprecated = true
	}
}

func OperationID(id string) Option {
	return func(router *Router) {
		router.OperationID = id
	}
}

func Responses(responses map[string]*Response) Option {
	return func(router *Router) {
		router.Responses = responses
	}
}

func Model(model any) Option {
	return func(router *Router) {
		router.Model = model
	}
}
