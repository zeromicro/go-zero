<img align="right" width="150px" src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/go-zero.png">

# mapreduce

English | [简体中文](readme-cn.md)

## Why MapReduce is needed

In practical business scenarios we often need to get the corresponding properties from different rpc services to assemble complex objects.

For example, to query product details.

1. product service - query product attributes
2. inventory service - query inventory properties
3. price service - query price attributes
4. marketing service - query marketing properties

If it is a serial call, the response time will increase linearly with the number of rpc calls, so we will generally change serial to parallel to optimize response time.

Simple scenarios using `WaitGroup` can also meet the needs, but what if we need to check the data returned by the rpc call, data processing, data aggregation? The official go library does not have such a tool (CompleteFuture is provided in java), so we implemented an in-process data batching MapReduce concurrent tool based on the MapReduce architecture.

## Design ideas

Let's try to put ourselves in the author's shoes and sort out the possible business scenarios for the concurrency tool:

1. querying product details: supporting concurrent calls to multiple services to combine product attributes, and supporting call errors that can be ended immediately.
2. automatic recommendation of user card coupons on product details page: support concurrently verifying card coupons, automatically rejecting them if they fail, and returning all of them.

The above is actually processing the input data and finally outputting the cleaned data. There is a very classic asynchronous pattern for data processing: the producer-consumer pattern. So we can abstract the life cycle of data batch processing, which can be roughly divided into three phases.

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/mapreduce-serial-en.png" width="500">

1. data production generate
2. data processing mapper
3. data aggregation reducer

Data producing is an indispensable stage, data processing and data aggregation are optional stages, data producing and processing support concurrent calls, data aggregation is basically a pure memory operation, so a single concurrent process can do it.

Since different stages of data processing are performed by different goroutines, it is natural to consider the use of channel to achieve communication between goroutines.

<img src="https://raw.githubusercontent.com/zeromicro/zero-doc/main/doc/images/mapreduce-en.png" width="500">

How can I terminate the process at any time?

It's simple, just receive from a  channel or the given context in the goroutine.

## A simple example

Calculate the sum of squares, simulating the concurrency.

```go
package main

import (
    "fmt"
    "log"

    "github.com/zeromicro/go-zero/core/mr"
)

func main() {
    val, err := mr.MapReduce(func(source chan<- int) {
        // generator
        for i := 0; i < 10; i++ {
            source <- i
        }
    }, func(i int, writer mr.Writer[int], cancel func(error)) {
        // mapper
        writer.Write(i * i)
    }, func(pipe <-chan int, writer mr.Writer[int], cancel func(error)) {
        // reducer
        var sum int
        for i := range pipe {
            sum += i
        }
        writer.Write(sum)
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("result:", val)
}
```

More examples: [https://github.com/zeromicro/zero-examples/tree/main/mapreduce](https://github.com/zeromicro/zero-examples/tree/main/mapreduce)

## Give a Star! ⭐

If you like or are using this project to learn or start your solution, please give it a star. Thanks!
