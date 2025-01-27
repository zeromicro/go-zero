package conf

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockType int

func (m mockType) Validate() error {
	if m < 10 {
		return errors.New("invalid value")
	}

	return nil
}

type anotherMockType int

func Test_validate(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		wantErr bool
	}{
		{
			name:    "invalid",
			v:       mockType(5),
			wantErr: true,
		},
		{
			name:    "valid",
			v:       mockType(10),
			wantErr: false,
		},
		{
			name:    "not validator",
			v:       anotherMockType(5),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.v)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

type mockVal struct {
}

func (m mockVal) Validate() error {
	return errors.New("invalid value")
}

func Test_validateValPtr(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		wantErr bool
	}{
		{
			name: "invalid",
			v:    mockVal{},
		},
		{
			name: "invalid value",
			v:    &mockVal{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, validate(tt.v))
		})
	}
}
