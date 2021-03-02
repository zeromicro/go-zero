package utils

import "github.com/google/uuid"

// NewUuid returns a uuid string.
func NewUuid() string {
	return uuid.New().String()
}
