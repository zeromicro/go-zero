package handler

import (
	"net/http"

	"zero/ngin"
)

func RegisterHandlers(engine *ngin.Engine) {
	engine.AddRoutes([]ngin.Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: GreetHandler,
		},
	})
}
