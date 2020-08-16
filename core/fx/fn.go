package fx

import (
	"sort"
	"sync"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/threading"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

type (
	rxOptions struct {
		unlimitedWorkers bool
		workers          int
	}

	FilterFunc   func(item interface{}) bool
	ForAllFunc   func(pipe <-chan interface{})
	ForEachFunc  func(item interface{})
	GenerateFunc func(source chan<- interface{})
	KeyFunc      func(item interface{}) interface{}
	LessFunc     func(a, b interface{}) bool
	MapFunc      func(item interface{}) interface{}
	Option       func(opts *rxOptions)
	ParallelFunc func(item interface{})
	ReduceFunc   func(pipe <-chan interface{}) (interface{}, error)
	WalkFunc     func(item interface{}, pipe chan<- interface{})

	Stream struct {
		source <-chan interface{}
	}
)

// From constructs a Stream from the given GenerateFunc.
func From(generate GenerateFunc) Stream {
	source := make(chan interface{})

	threading.GoSafe(func() {
		defer close(source)
		generate(source)
	})

	return Range(source)
}

// Just converts the given arbitrary items to a Stream.
func Just(items ...interface{}) Stream {
	source := make(chan interface{}, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)

	return Range(source)
}

// Range converts the given channel to a Stream.
func Range(source <-chan interface{}) Stream {
	return Stream{
		source: source,
	}
}

// Buffer buffers the items into a queue with size n.
func (p Stream) Buffer(n int) Stream {
	if n < 0 {
		n = 0
	}

	source := make(chan interface{}, n)
	go func() {
		for item := range p.source {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

// Distinct removes the duplicated items base on the given KeyFunc.
func (p Stream) Distinct(fn KeyFunc) Stream {
	source := make(chan interface{})

	threading.GoSafe(func() {
		defer close(source)

		keys := make(map[interface{}]lang.PlaceholderType)
		for item := range p.source {
			key := fn(item)
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = lang.Placeholder
			}
		}
	})

	return Range(source)
}

// Done waits all upstreaming operations to be done.
func (p Stream) Done() {
	for range p.source {
	}
}

// Filter filters the items by the given FilterFunc.
func (p Stream) Filter(fn FilterFunc, opts ...Option) Stream {
	return p.Walk(func(item interface{}, pipe chan<- interface{}) {
		if fn(item) {
			pipe <- item
		}
	}, opts...)
}

// ForAll handles the streaming elements from the source and no later streams.
func (p Stream) ForAll(fn ForAllFunc) {
	fn(p.source)
}

// ForEach seals the Stream with the ForEachFunc on each item, no successive operations.
func (p Stream) ForEach(fn ForEachFunc) {
	for item := range p.source {
		fn(item)
	}
}

// Group groups the elements into different groups based on their keys.
func (p Stream) Group(fn KeyFunc) Stream {
	groups := make(map[interface{}][]interface{})
	for item := range p.source {
		key := fn(item)
		groups[key] = append(groups[key], item)
	}

	source := make(chan interface{})
	go func() {
		for _, group := range groups {
			source <- group
		}
		close(source)
	}()

	return Range(source)
}

func (p Stream) Head(n int64) Stream {
	source := make(chan interface{})

	go func() {
		for item := range p.source {
			n--
			if n >= 0 {
				source <- item
			}
			if n == 0 {
				// let successive method go ASAP even we have more items to skip
				// why we don't just break the loop, because if break,
				// this former goroutine will block forever, which will cause goroutine leak.
				close(source)
			}
		}
		if n > 0 {
			close(source)
		}
	}()

	return Range(source)
}

// Maps converts each item to another corresponding item, which means it's a 1:1 model.
func (p Stream) Map(fn MapFunc, opts ...Option) Stream {
	return p.Walk(func(item interface{}, pipe chan<- interface{}) {
		pipe <- fn(item)
	}, opts...)
}

// Merge merges all the items into a slice and generates a new stream.
func (p Stream) Merge() Stream {
	var items []interface{}
	for item := range p.source {
		items = append(items, item)
	}

	source := make(chan interface{}, 1)
	source <- items
	close(source)

	return Range(source)
}

// Parallel applies the given ParallelFunc to each item concurrently with given number of workers.
func (p Stream) Parallel(fn ParallelFunc, opts ...Option) {
	p.Walk(func(item interface{}, pipe chan<- interface{}) {
		fn(item)
	}, opts...).Done()
}

// Reduce is a utility method to let the caller deal with the underlying channel.
func (p Stream) Reduce(fn ReduceFunc) (interface{}, error) {
	return fn(p.source)
}

// Reverse reverses the elements in the stream.
func (p Stream) Reverse() Stream {
	var items []interface{}
	for item := range p.source {
		items = append(items, item)
	}
	// reverse, official method
	for i := len(items)/2 - 1; i >= 0; i-- {
		opp := len(items) - 1 - i
		items[i], items[opp] = items[opp], items[i]
	}

	return Just(items...)
}

// Sort sorts the items from the underlying source.
func (p Stream) Sort(less LessFunc) Stream {
	var items []interface{}
	for item := range p.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})

	return Just(items...)
}

func (p Stream) Tail(n int64) Stream {
	source := make(chan interface{})

	go func() {
		ring := collection.NewRing(int(n))
		for item := range p.source {
			ring.Add(item)
		}
		for _, item := range ring.Take() {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

// Walk lets the callers handle each item, the caller may write zero, one or more items base on the given item.
func (p Stream) Walk(fn WalkFunc, opts ...Option) Stream {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return p.walkUnlimited(fn, option)
	} else {
		return p.walkLimited(fn, option)
	}
}

func (p Stream) walkLimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan interface{}, option.workers)

	go func() {
		var wg sync.WaitGroup
		pool := make(chan lang.PlaceholderType, option.workers)

		for {
			pool <- lang.Placeholder
			item, ok := <-p.source
			if !ok {
				<-pool
				break
			}

			wg.Add(1)
			// better to safely run caller defined method
			threading.GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				fn(item, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

func (p Stream) walkUnlimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan interface{}, defaultWorkers)

	go func() {
		var wg sync.WaitGroup

		for {
			item, ok := <-p.source
			if !ok {
				break
			}

			wg.Add(1)
			// better to safely run caller defined method
			threading.GoSafe(func() {
				defer wg.Done()
				fn(item, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

// UnlimitedWorkers lets the caller to use as many workers as the tasks.
func UnlimitedWorkers() Option {
	return func(opts *rxOptions) {
		opts.unlimitedWorkers = true
	}
}

// WithWorkers lets the caller to customize the concurrent workers.
func WithWorkers(workers int) Option {
	return func(opts *rxOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func buildOptions(opts ...Option) *rxOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func newOptions() *rxOptions {
	return &rxOptions{
		workers: defaultWorkers,
	}
}
