package handler

import (
	"net/http"

	"zero/rest"
)

func RegisterHandlers(engine *rest.Engine) {
	engine.AddRoutes([]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: GreetHandler,
		},
	})
}
