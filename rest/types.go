package rest

import "net/http"

type (
	// Middleware defines the middleware method.
	Middleware func(next http.HandlerFunc) http.HandlerFunc

	// A Route is a http route.
	Route struct {
		Method  string
		Path    string
		Handler http.HandlerFunc
	}

	// RouteOption defines the method to customize a featured route.
	RouteOption func(r *featuredRoutes)

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
)
