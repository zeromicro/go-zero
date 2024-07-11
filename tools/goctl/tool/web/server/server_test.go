package server

import (
	"log"
	"testing"
)

func TestLocal(t *testing.T) {
	//t.Skip("local testing")
	if err := Run(8080); err != nil {
		log.Fatal(err)
	}
}
