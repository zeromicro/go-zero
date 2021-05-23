package mr

import (
	"errors"
	"fmt"
	"sync"

	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/syncx"
	"github.com/tal-tech/go-zero/core/threading"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

var (
	// ErrCancelWithNil is an error that mapreduce was cancelled with nil.
	ErrCancelWithNil = errors.New("mapreduce cancelled with nil")
	// ErrReduceNoOutput is an error that reduce did not output a value.
	ErrReduceNoOutput = errors.New("reduce not writing value")
)

type (
	// GenerateFunc is used to let callers send elements into source.
	GenerateFunc func(source chan<- interface{})
	// MapFunc is used to do element processing and write the output to writer.
	MapFunc func(item interface{}, writer Writer)
	// VoidMapFunc is used to do element processing, but no output.
	VoidMapFunc func(item interface{})
	// MapperFunc is used to do element processing and write the output to writer,
	// use cancel func to cancel the processing.
	MapperFunc func(item interface{}, writer Writer, cancel func(error))
	// ReducerFunc is used to reduce all the mapping output and write to writer,
	// use cancel func to cancel the processing.
	ReducerFunc func(pipe <-chan interface{}, writer Writer, cancel func(error))
	// VoidReducerFunc is used to reduce all the mapping output, but no output.
	// Use cancel func to cancel the processing.
	VoidReducerFunc func(pipe <-chan interface{}, cancel func(error))
	// Option defines the method to customize the mapreduce.
	Option func(opts *mapReduceOptions)

	mapReduceOptions struct {
		workers int
	}

	// Writer interface wraps Write method.
	Writer interface {
		Write(v interface{})
	}
)

// Finish runs fns parallelly, cancelled on any error.
func Finish(fns ...func() error) error {
	if len(fns) == 0 {
		return nil
	}

	return MapReduceVoid(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}, writer Writer, cancel func(error)) {
		fn := item.(func() error)
		if err := fn(); err != nil {
			cancel(err)
		}
	}, func(pipe <-chan interface{}, cancel func(error)) {
		drain(pipe)
	}, WithWorkers(len(fns)))
}

// FinishVoid runs fns parallelly.
func FinishVoid(fns ...func()) {
	if len(fns) == 0 {
		return
	}

	MapVoid(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}) {
		fn := item.(func())
		fn()
	}, WithWorkers(len(fns)))
}

// Map maps all elements generated from given generate func, and returns an output channel.
func Map(generate GenerateFunc, mapper MapFunc, opts ...Option) chan interface{} {
	options := buildOptions(opts...)
	source := buildSource(generate)
	collector := make(chan interface{}, options.workers)
	done := syncx.NewDoneChan()

	go executeMappers(mapper, source, collector, done.Done(), options.workers)

	return collector
}

// MapReduce maps all elements generated from given generate func,
// and reduces the output elements with given reducer.
func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc, opts ...Option) (interface{}, error) {
	source := buildSource(generate)
	return MapReduceWithSource(source, mapper, reducer, opts...)
}

// MapReduceWithSource maps all elements from source, and reduce the output elements with given reducer.
func MapReduceWithSource(source <-chan interface{}, mapper MapperFunc, reducer ReducerFunc,
	opts ...Option) (interface{}, error) {
	options := buildOptions(opts...)
	output := make(chan interface{})
	collector := make(chan interface{}, options.workers)
	done := syncx.NewDoneChan()
	writer := newGuardedWriter(output, done.Done())
	var closeOnce sync.Once
	var retErr errorx.AtomicError
	finish := func() {
		closeOnce.Do(func() {
			done.Close()
			close(output)
		})
	}
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}

		drain(source)
		finish()
	})

	go func() {
		defer func() {
			drain(collector)

			if r := recover(); r != nil {
				cancel(fmt.Errorf("%v", r))
			} else {
				finish()
			}
		}()

		reducer(collector, writer, cancel)
	}()

	go executeMappers(func(item interface{}, w Writer) {
		mapper(item, w, cancel)
	}, source, collector, done.Done(), options.workers)

	value, ok := <-output
	if err := retErr.Load(); err != nil {
		return nil, err
	} else if ok {
		return value, nil
	} else {
		return nil, ErrReduceNoOutput
	}
}

// MapReduceVoid maps all elements generated from given generate,
// and reduce the output elements with given reducer.
func MapReduceVoid(generate GenerateFunc, mapper MapperFunc, reducer VoidReducerFunc, opts ...Option) error {
	_, err := MapReduce(generate, mapper, func(input <-chan interface{}, writer Writer, cancel func(error)) {
		reducer(input, cancel)
		// We need to write a placeholder to let MapReduce to continue on reducer done,
		// otherwise, all goroutines are waiting. The placeholder will be discarded by MapReduce.
		writer.Write(lang.Placeholder)
	}, opts...)
	return err
}

// MapVoid maps all elements from given generate but no output.
func MapVoid(generate GenerateFunc, mapper VoidMapFunc, opts ...Option) {
	drain(Map(generate, func(item interface{}, writer Writer) {
		mapper(item)
	}, opts...))
}

// WithWorkers customizes a mapreduce processing with given workers.
func WithWorkers(workers int) Option {
	return func(opts *mapReduceOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

func buildOptions(opts ...Option) *mapReduceOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func buildSource(generate GenerateFunc) chan interface{} {
	source := make(chan interface{})
	threading.GoSafe(func() {
		defer close(source)
		generate(source)
	})

	return source
}

// drain drains the channel.
func drain(channel <-chan interface{}) {
	// drain the channel
	for range channel {
	}
}

func executeMappers(mapper MapFunc, input <-chan interface{}, collector chan<- interface{},
	done <-chan lang.PlaceholderType, workers int) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(collector)
	}()

	pool := make(chan lang.PlaceholderType, workers)
	writer := newGuardedWriter(collector, done)
	for {
		select {
		case <-done:
			return
		case pool <- lang.Placeholder:
			item, ok := <-input
			if !ok {
				<-pool
				return
			}

			wg.Add(1)
			// better to safely run caller defined method
			threading.GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				mapper(item, writer)
			})
		}
	}
}

func newOptions() *mapReduceOptions {
	return &mapReduceOptions{
		workers: defaultWorkers,
	}
}

func once(fn func(error)) func(error) {
	once := new(sync.Once)
	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}

type guardedWriter struct {
	channel chan<- interface{}
	done    <-chan lang.PlaceholderType
}

func newGuardedWriter(channel chan<- interface{}, done <-chan lang.PlaceholderType) guardedWriter {
	return guardedWriter{
		channel: channel,
		done:    done,
	}
}

func (gw guardedWriter) Write(v interface{}) {
	select {
	case <-gw.done:
		return
	default:
		gw.channel <- v
	}
}
