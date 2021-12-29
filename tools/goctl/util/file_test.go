package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGitHome(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	actual, err := GetGitHome()
	if err != nil {
		return
	}

	expected := filepath.Join(homeDir, goctlDir, gitDir)
	assert.Equal(t, expected, actual)
}
