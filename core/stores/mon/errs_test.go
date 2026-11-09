package mon

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestAcceptable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"NilError", nil, true},
		{"NoDocuments", mongo.ErrNoDocuments, true},
		{"NilValue", mongo.ErrNilValue, true},
		{"NilDocument", mongo.ErrNilDocument, true},
		{"NilCursor", mongo.ErrNilCursor, true},
		{"EmptySlice", mongo.ErrEmptySlice, true},
		{"DuplicateKeyError", mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: duplicateKeyCode}}}, true},
		{"OtherError", errors.New("other error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, acceptable(tt.err))
		})
	}
}

func TestIsDupKeyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"NilError", nil, false},
		{"NonDupKeyError", errors.New("some other error"), false},
		{"DupKeyError", mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: duplicateKeyCode}}}, true},
		{"OtherMongoError", mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 12345}}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isDupKeyError(tt.err))
		})
	}
}
