package main

import (
	"log"
	"strconv"

	"github.com/tal-tech/go-zero/core/mr"
)

type User struct {
	Uid  int
	Name string
}

func main() {
	uids := []int{111, 222, 333}
	res, err := mr.MapReduce(func(source chan<- interface{}) {
		for _, uid := range uids {
			source <- uid
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		uid := item.(int)
		user := &User{
			Uid:  uid,
			Name: strconv.Itoa(uid),
		}
		writer.Write(user)
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		var users []*User
		for p := range pipe {
			users = append(users, p.(*User))
		}
		// missing writer.Write(...), should not panic
	})
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(len(res.([]*User)))
}
