package rest

import "net/http"

type (
	Route struct {
		Method  string
		Path    string
		Handler http.HandlerFunc
	}

	jwtSetting struct {
		enabled    bool
		secret     string
		prevSecret string
	}

	signatureSetting struct {
		SignatureConf
		enabled bool
	}

	featuredRoutes struct {
		priority  bool
		jwt       jwtSetting
		signature signatureSetting
		routes    []Route
	}

	RouteOption func(r *featuredRoutes)
)
