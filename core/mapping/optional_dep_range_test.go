package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalOptionalDepKeepsRange(t *testing.T) {
	type config struct {
		Dep string `json:"dep,optional"`
		X   int    `json:"x,optional=dep,range=[1:10]"`
	}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "x in range",
			input: `{"dep":"yes","x":5}`,
		},
		{
			name:    "x above range",
			input:   `{"dep":"yes","x":100}`,
			wantErr: true,
		},
		{
			name:    "x below range",
			input:   `{"dep":"yes","x":0}`,
			wantErr: true,
		},
		{
			name:  "both absent",
			input: `{}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v config
			err := UnmarshalJsonBytes([]byte(test.input), &v)
			if test.wantErr {
				assert.ErrorIs(t, err, errNumberRange)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
