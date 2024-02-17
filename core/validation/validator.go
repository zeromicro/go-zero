package validation

// Validator represents a validator.
type Validator interface {
	// Validate validates the value.
	Validate() error
}
