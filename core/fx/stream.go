package fx

import (
	"sort"
	"sync"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/threading"
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

	// FilterFunc defines the method to filter a Stream.
	FilterFunc[T any] func(item T) bool
	// ForAllFunc defines the method to handle all elements in a Stream.
	ForAllFunc[T any] func(pipe <-chan T)
	// ForEachFunc defines the method to handle each element in a Stream.
	ForEachFunc[T any] func(item T)
	// GenerateFunc defines the method to send elements into a Stream.
	GenerateFunc[T any] func(source chan<- T)
	// KeyFunc defines the method to generate keys for the elements in a Stream.
	KeyFunc func(item any) any
	// LessFunc defines the method to compare the elements in a Stream.
	LessFunc func(a, b any) bool
	// MapFunc defines the method to map each element to another object in a Stream.
	MapFunc func(item any) any
	// Option defines the method to customize a Stream.
	Option func(opts *rxOptions)
	// ParallelFunc defines the method to handle elements parallel.
	ParallelFunc func(item any)
	// ReduceFunc defines the method to reduce all the elements in a Stream.
	ReduceFunc[T any] func(pipe <-chan T) (T, error)
	// WalkFunc defines the method to walk through all the elements in a Stream.
	WalkFunc[T any] func(item T, pipe chan<- T)

	// A Stream is a stream that can be used to do stream processing.
	Stream[T any] struct {
		source <-chan T
	}
)

// Concat returns a concatenated Stream.
func Concat[T any](s Stream[T], others ...Stream[T]) Stream[T] {
	return s.Concat(others...)
}

// From constructs a Stream from the given GenerateFunc.
func From[T any](generate GenerateFunc[T]) Stream[T] {
	source := make(chan T)

	threading.GoSafe(func() {
		defer close(source)
		generate(source)
	})

	return Range(source)
}

// Just converts the given arbitrary items to a Stream.
func Just[T any](items ...T) Stream[T] {
	source := make(chan T, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)

	return Range(source)
}

// Range converts the given channel to a Stream.
func Range[T any](source <-chan T) Stream[T] {
	return Stream[T]{
		source: source,
	}
}

// AllMach returns whether all elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s Stream[T]) AllMach(predicate func(item any) bool) bool {
	for item := range s.source {
		if !predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return false
		}
	}

	return true
}

// AnyMach returns whether any elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then false is returned and the predicate is not evaluated.
func (s Stream[T]) AnyMach(predicate func(item any) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return true
		}
	}

	return false
}

// Buffer buffers the items into a queue with size n.
// It can balance the producer and the consumer if their processing throughput don't match.
func (s Stream[T]) Buffer(n int) Stream[T] {
	if n < 0 {
		n = 0
	}

	source := make(chan T, n)
	go func() {
		for item := range s.source {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

// Concat returns a Stream that concatenated other streams
func (s Stream[T]) Concat(others ...Stream[T]) Stream[T] {
	source := make(chan T)

	go func() {
		group := threading.NewRoutineGroup()
		group.Run(func() {
			for item := range s.source {
				source <- item
			}
		})

		for _, each := range others {
			each := each
			group.Run(func() {
				for item := range each.source {
					source <- item
				}
			})
		}

		group.Wait()
		close(source)
	}()

	return Range(source)
}

// Count counts the number of elements in the result.
func (s Stream[T]) Count() (count int) {
	for range s.source {
		count++
	}
	return
}

// Distinct removes the duplicated items base on the given KeyFunc.
func (s Stream[T]) Distinct(fn KeyFunc) Stream[T] {
	source := make(chan T)

	threading.GoSafe(func() {
		defer close(source)

		keys := make(map[any]lang.PlaceholderType)
		for item := range s.source {
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
func (s Stream[T]) Done() {
	drain(s.source)
}

// Filter filters the items by the given FilterFunc.
func (s Stream[T]) Filter(fn FilterFunc[T], opts ...Option) Stream[T] {
	return s.Walk(func(item T, pipe chan<- T) {
		if fn(item) {
			pipe <- item
		}
	}, opts...)
}

// First returns the first item, nil if no items.
func (s Stream[T]) First() any {
	for item := range s.source {
		// make sure the former goroutine not block, and current func returns fast.
		go drain(s.source)
		return item
	}

	return nil
}

// ForAll handles the streaming elements from the source and no later streams.
func (s Stream[T]) ForAll(fn ForAllFunc[T]) {
	fn(s.source)
	// avoid goroutine leak on fn not consuming all items.
	go drain(s.source)
}

// ForEach seals the Stream with the ForEachFunc on each item, no successive operations.
func (s Stream[T]) ForEach(fn ForEachFunc[T]) {
	for item := range s.source {
		fn(item)
	}
}

// Group groups the elements into different groups based on their keys.
func (s Stream[T]) Group(fn KeyFunc) Stream[[]T] {
	groups := make(map[any][]T)
	for item := range s.source {
		key := fn(item)
		groups[key] = append(groups[key], item)
	}

	source := make(chan []T)
	go func() {
		for _, group := range groups {
			source <- group
		}
		close(source)
	}()

	return Range(source)
}

// Head returns the first n elements in p.
func (s Stream[T]) Head(n int64) Stream[T] {
	if n < 1 {
		panic("n must be greater than 0")
	}

	source := make(chan T)

	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				source <- item
			}
			if n == 0 {
				// let successive method go ASAP even we have more items to skip
				close(source)
				// why we don't just break the loop, and drain to consume all items.
				// because if breaks, this former goroutine will block forever,
				// which will cause goroutine leak.
				drain(s.source)
			}
		}
		// not enough items in s.source, but we need to let successive method to go ASAP.
		if n > 0 {
			close(source)
		}
	}()

	return Range(source)
}

// Last returns the last item, or nil if no items.
func (s Stream[T]) Last() (item any) {
	for item = range s.source {
	}
	return
}

// Map converts each item to another corresponding item, which means it's a 1:1 model.
func (s Stream[T]) Map(fn MapFunc, opts ...Option) Stream[T] {
	return s.Walk(func(item T, pipe chan<- T) {
		pipe <- fn(item)
	}, opts...)
}

// Merge merges all the items into a slice and generates a new stream.
func (s Stream[T]) Merge() Stream[[]T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}

	source := make(chan []T, 1)
	source <- items
	close(source)

	return Range(source)
}

// NoneMatch returns whether all elements of this stream don't match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s Stream[T]) NoneMatch(predicate func(item any) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return false
		}
	}

	return true
}

// Parallel applies the given ParallelFunc to each item concurrently with given number of workers.
func (s Stream[T]) Parallel(fn ParallelFunc, opts ...Option) {
	s.Walk(func(item T, pipe chan<- T) {
		fn(item)
	}, opts...).Done()
}

// Reduce is a utility method to let the caller deal with the underlying channel.
func (s Stream[T]) Reduce(fn ReduceFunc[T]) (T, error) {
	return fn(s.source)
}

// Reverse reverses the elements in the stream.
func (s Stream[T]) Reverse() Stream[T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}
	// reverse, official method
	for i := len(items)/2 - 1; i >= 0; i-- {
		opp := len(items) - 1 - i
		items[i], items[opp] = items[opp], items[i]
	}

	return Just(items...)
}

// Skip returns a Stream that skips size elements.
func (s Stream[T]) Skip(n int64) Stream[T] {
	if n < 0 {
		panic("n must not be negative")
	}
	if n == 0 {
		return s
	}

	source := make(chan T)

	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				continue
			} else {
				source <- item
			}
		}
		close(source)
	}()

	return Range(source)
}

// Sort sorts the items from the underlying source.
func (s Stream[T]) Sort(less LessFunc) Stream[T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})

	return Just(items...)
}

// Split splits the elements into chunk with size up to n,
// might be less than n on tailing elements.
func (s Stream[T]) Split(n int) Stream[[]T] {
	if n < 1 {
		panic("n should be greater than 0")
	}

	source := make(chan []T)
	go func() {
		var chunk []T
		for item := range s.source {
			chunk = append(chunk, item)
			if len(chunk) == n {
				source <- chunk
				chunk = nil
			}
		}
		if chunk != nil {
			source <- chunk
		}
		close(source)
	}()

	return Range(source)
}

// Tail returns the last n elements in p.
func (s Stream[T]) Tail(n int64) Stream[T] {
	if n < 1 {
		panic("n should be greater than 0")
	}

	source := make(chan T)

	go func() {
		ring := collection.NewRing(int(n))
		for item := range s.source {
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
func (s Stream[T]) Walk(fn WalkFunc[T], opts ...Option) Stream[T] {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return s.walkUnlimited(fn, option)
	}

	return s.walkLimited(fn, option)
}

func (s Stream[T]) walkLimited(fn WalkFunc[T], option *rxOptions) Stream[T] {
	pipe := make(chan T, option.workers)

	go func() {
		var wg sync.WaitGroup
		pool := make(chan lang.PlaceholderType, option.workers)

		for item := range s.source {
			// important, used in another goroutine
			val := item
			pool <- lang.Placeholder
			wg.Add(1)

			// better to safely run caller defined method
			threading.GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				fn(val, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

func (s Stream[T]) walkUnlimited(fn WalkFunc[T], option *rxOptions) Stream[T] {
	pipe := make(chan T, option.workers)

	go func() {
		var wg sync.WaitGroup

		for item := range s.source {
			// important, used in another goroutine
			val := item
			wg.Add(1)
			// better to safely run caller defined method
			threading.GoSafe(func() {
				defer wg.Done()
				fn(val, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

// UnlimitedWorkers lets the caller use as many workers as the tasks.
func UnlimitedWorkers() Option {
	return func(opts *rxOptions) {
		opts.unlimitedWorkers = true
	}
}

// WithWorkers lets the caller customize the concurrent workers.
func WithWorkers(workers int) Option {
	return func(opts *rxOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

// buildOptions returns a rxOptions with given customizations.
func buildOptions(opts ...Option) *rxOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// drain drains the given channel.
func drain[T any](channel <-chan T) {
	for range channel {
	}
}

// newOptions returns a default rxOptions.
func newOptions() *rxOptions {
	return &rxOptions{
		workers: defaultWorkers,
	}
}
