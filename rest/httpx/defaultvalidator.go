package httpx

import (
	"net/http"

	playgroundvalidator "github.com/go-playground/validator/v10"
)

type defaultValidator struct {
	validate *playgroundvalidator.Validate
}

func newDefaultValidator() *defaultValidator {
	return &defaultValidator{
		validate: playgroundvalidator.New(),
	}
}

// Validate only validates the data, not the http request.
//
// The validation rule for a http request is various,
// if we use a tag to register the validation rule, it will be too complex,
// so currently we do not validate the http request, only validate the data.
func (v *defaultValidator) Validate(r *http.Request, data any) error {

	return v.validate.Struct(data)
}

func (v *defaultValidator) RegisterValidation(tag string, fn playgroundvalidator.Func) error {
	err := v.validate.RegisterValidation(tag, fn)
	return err
}

func init() {
	v := newDefaultValidator()
	SetValidator(v)
}
