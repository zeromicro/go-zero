package utils

import (
	"context"
	"fmt"
	"testing"
)

func TestGO(t *testing.T) {
	GO(context.Background(), func(ctx context.Context) {
		fmt.Println("go")
		panic("1234")
	})
}
