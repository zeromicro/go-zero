package handler

import (
	"net/http"

	"zero/rest/httpx"
)

type (
	request struct {
		User string `form:"user,optional"`
	}

	response struct {
		Code  int    `json:"code"`
		Greet string `json:"greet"`
		From  string `json:"from,omitempty"`
	}
)

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	var req request
	err := httpx.Parse(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	httpx.OkJson(w, response{
		Code:  0,
		Greet: "hello",
	})
}
