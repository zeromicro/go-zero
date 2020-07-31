package internal

import (
	"context"
	"net/http"

	"zero/core/proc"
)

func StartServer(srv *http.Server) error {
	proc.AddWrapUpListener(func() {
		srv.Shutdown(context.Background())
	})

	return srv.ListenAndServe()
}
