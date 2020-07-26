package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"zero/core/executors"
	"zero/core/httpx"
	"zero/core/logx"
	"zero/example/graceful/dns/api/svc"
	"zero/example/graceful/dns/api/types"
	"zero/example/graceful/dns/rpc/graceful"
)

func gracefulHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	logger := executors.NewLessExecutor(time.Second)
	return func(w http.ResponseWriter, r *http.Request) {
		var resp types.Response

		conn, ok := ctx.Client.Next()
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		host, err := os.Hostname()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
			return
		}

		client := graceful.NewGraceServiceClient(conn)
		rp, err := client.Grace(context.Background(), &graceful.Request{From: host})
		if err != nil {
			logx.Error(err)
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}

		resp.Host = rp.Host
		logger.DoOrDiscard(func() {
			fmt.Printf("%s from host: %s\n", time.Now().Format("15:04:05"), rp.Host)
		})
		httpx.OkJson(w, resp)
	}
}
