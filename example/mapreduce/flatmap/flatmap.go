package main

import (
	"fmt"

	"zero/core/mapreduce"
)

var (
	persons = []string{"john", "mary", "alice", "bob"}
	friends = map[string][]string{
		"john":  {"harry", "hermione", "ron"},
		"mary":  {"sam", "frodo"},
		"alice": {},
		"bob":   {"jamie", "tyrion", "cersei"},
	}
)

func main() {
	var allFriends []string
	for v := range mapreduce.Map(func(source chan<- interface{}) {
		for _, each := range persons {
			source <- each
		}
	}, func(item interface{}, writer mapreduce.Writer) {
		writer.Write(friends[item.(string)])
	}, mapreduce.WithWorkers(100)) {
		allFriends = append(allFriends, v.([]string)...)
	}
	fmt.Println(allFriends)
}
