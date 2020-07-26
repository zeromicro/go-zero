package main

import (
	"encoding/json"
	"fmt"
	"log"

	jsonx "github.com/segmentio/encoding/json"
)

type A struct {
	AA string `json:"aa,omitempty"`
}

type B struct {
	*A
	BB string `json:"bb,omitempty"`
}

func main() {
	var b B
	b.BB = "b"
	b.A = new(A)
	b.A.AA = ""

	fmt.Println("github.com/segmentio/encoding/json")
	data, err := jsonx.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
	fmt.Println()

	fmt.Println("encoding/json")
	data, err = json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}
