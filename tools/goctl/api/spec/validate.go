package spec

import "errors"

var ErrMissingService = errors.New("missing service")

// Validate validates Validate the integrity of the spec.
func (s *ApiSpec) Validate() error {
	if len(s.Service.Name) == 0 {
		return ErrMissingService
	}
	if len(s.Service.Groups) == 0 {
		return ErrMissingService
	}
	return nil
}
