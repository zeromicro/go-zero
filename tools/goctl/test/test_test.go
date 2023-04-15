package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutor_Run(t *testing.T) {
	executor := NewExecutor[string, string]()
	executor.Add([]Data[string, string]{
		{
			Name: "empty",
		},
		{
			Name:  "snake_case",
			Input: "A_B_C",
			Want:  "a_b_c",
		},
		{
			Name:  "camel_case",
			Input: "AaBbCc",
			Want:  "aabbcc",
		},
	}...)
	executor.Run(t, func(s string) string {
		return strings.ToLower(s)
	})
}

func TestExecutor_RunE(t *testing.T) {
	var dummyError = errors.New("dummy error")
	executor := NewExecutor[string, string]()
	executor.Add([]Data[string, string]{
		{
			Name: "empty",
		},
		{
			Name:  "snake_case",
			Input: "A_B_C",
			Want:  "a_b_c",
		},
		{
			Name:  "camel_case",
			Input: "AaBbCc",
			Want:  "aabbcc",
		},
		{
			Name:  "invalid_input",
			Input: "😄",
			E:     dummyError,
		},
	}...)
	executor.RunE(t, func(s string) (string, error) {
		for _, r := range s {
			if r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
				continue
			}
			return "", dummyError
		}
		return strings.ToLower(s), nil
	})
}

func TestWithComparison(t *testing.T) {
	var dummyError = errors.New("dummy error")
	executor := NewExecutor[string, string](WithComparison[string, string](func(t *testing.T, expected, actual string) bool {
		return assert.Equal(t, expected, actual)
	}))
	executor.Add([]Data[string, string]{
		{
			Name: "empty",
		},
		{
			Name:  "snake_case",
			Input: "A_B_C",
			Want:  "a_b_c",
		},
		{
			Name:  "camel_case",
			Input: "AaBbCc",
			Want:  "aabbcc",
		},
		{
			Name:  "invalid_input",
			Input: "😄",
			E:     dummyError,
		},
	}...)
	executor.RunE(t, func(s string) (string, error) {
		for _, r := range s {
			if r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
				continue
			}
			return "", dummyError
		}
		return strings.ToLower(s), nil
	})
}
