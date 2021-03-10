package httpx

/*
	copy gin redirect part function
	tks gin
*/
import (
	"fmt"
	"net/http"
)

type redirect struct {
	Code     int
	Request  *http.Request
	Location string
}

type Render interface {
	// Render writes data with custom ContentType.
	Render(http.ResponseWriter) error
	// WriteContentType writes custom ContentType.
	WriteContentType(w http.ResponseWriter)
}

// Render (Redirect) redirects the http request to new location and writes redirect response.
func (r redirect) Render(w http.ResponseWriter) error {
	if (r.Code < http.StatusMultipleChoices || r.Code > http.StatusPermanentRedirect) && r.Code != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", r.Code))
	}
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}

// WriteContentType (Redirect) don't write any ContentType.
func (r redirect) WriteContentType(http.ResponseWriter) {}

func Redirect(code int, location string, r *http.Request, w http.ResponseWriter) {
	render(-1, redirect{Code: code, Location: location, Request: r}, w)
}

func render(code int, r Render, w http.ResponseWriter) {

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(w)
		return
	}

	if err := r.Render(w); err != nil {
		panic(err)
	}
}

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
