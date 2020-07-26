package main

import (
	"fmt"
	"strconv"
	"time"

	"zero/dq"
)

func main() {
	producer := dq.NewProducer([]dq.Beanstalk{
		{
			Endpoint: "localhost:11300",
			Tube:     "tube",
		},
		{
			Endpoint: "localhost:11301",
			Tube:     "tube",
		},
		{
			Endpoint: "localhost:11302",
			Tube:     "tube",
		},
		{
			Endpoint: "localhost:11303",
			Tube:     "tube",
		},
		{
			Endpoint: "localhost:11304",
			Tube:     "tube",
		},
	})
	for i := 0; i < 5; i++ {
		_, err := producer.At([]byte(strconv.Itoa(i)), time.Now().Add(time.Second*10))
		if err != nil {
			fmt.Println(err)
		}
	}
}
