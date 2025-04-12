package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
	}{
		{"string", "test"},
		{"int", 42},
		{"bool", true},
		{"float", 3.14},
		{"struct", struct{ Name string }{"test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ptr(tt.v)
			assert.NotNil(t, got, "ptr() should not return nil")
			assert.Equal(t, tt.v, *got, "dereferenced pointer should equal input value")
		})
	}
}

type Event struct {
	Type string
	Data map[string]any
}

func parseEvent(input string) (*Event, error) {
	var evt Event
	var dataStr string

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "event:") {
			evt.Type = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "data:") {
			dataStr = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(dataStr) > 0 {
		if err := json.Unmarshal([]byte(dataStr), &evt.Data); err != nil {
			return nil, fmt.Errorf("failed to parse data: %w", err)
		}
	}

	return &evt, nil
}
