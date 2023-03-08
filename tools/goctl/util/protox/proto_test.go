package protox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindEndOfService(t *testing.T) {
	serviceStr := "service Example {\n}\n\nservice School {\n  \n}"
	exampleBeginIndex, exampleMidIndex, exampleEndIndex := FindBeginEndOfService(serviceStr, "Example")
	schoolBeginIndex, schoolMidIndex, schoolEndIndex := FindBeginEndOfService(serviceStr, "School")
	assert.Equal(t, 0, exampleBeginIndex)
	assert.Equal(t, 18, exampleEndIndex)
	assert.Equal(t, 16, exampleMidIndex)
	assert.Equal(t, 21, schoolBeginIndex)
	assert.Equal(t, 41, schoolEndIndex)
	assert.Equal(t, 36, schoolMidIndex)
}
