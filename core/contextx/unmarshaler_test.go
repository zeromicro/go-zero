package contextx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalContext(t *testing.T) {

	type Person struct {
		Name string `ctx:"name"`
		Age  int    `ctx:"age"`
	}

	var PersonNameKey = "name"
	var PersonAgeKey = "age"

	ctx := context.Background()
	ctx = context.WithValue(ctx, PersonNameKey, "kevin")
	ctx = context.WithValue(ctx, PersonAgeKey, 20)

	var person Person
	err := For(ctx, &person)

	assert.Nil(t, err)
	assert.Equal(t, "kevin", person.Name)
	assert.Equal(t, 20, person.Age)
}

func TestUnmarshalContextWithOptional(t *testing.T) {
	type Person struct {
		Name string `ctx:"name"`
		Age  int    `ctx:"age,optional"`
	}
	var PersonNameKey = "name"

	ctx := context.Background()
	ctx = context.WithValue(ctx, PersonNameKey, "kevin")

	var person Person
	err := For(ctx, &person)

	assert.Nil(t, err)
	assert.Equal(t, "kevin", person.Name)
	assert.Equal(t, 0, person.Age)
}

func TestUnmarshalContextWithMissing(t *testing.T) {
	type Person struct {
		Name string `ctx:"name"`
		Age  int    `ctx:"age"`
	}
	type name string
	const PersonNameKey name = "name"

	ctx := context.Background()
	ctx = context.WithValue(ctx, PersonNameKey, "kevin")

	var person Person
	err := For(ctx, &person)

	assert.NotNil(t, err)
}
