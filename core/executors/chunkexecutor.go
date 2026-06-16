package executors

import "time"

const defaultChunkSize = 1024 * 1024 // 1M

type (
	// ChunkOption defines the method to customize a ChunkExecutor.
	ChunkOption func(options *chunkOptions)

	// A ChunkExecutor is an executor to execute tasks when either requirement meets:
	// 1. up to given chunk size
	// 2. flush interval elapsed
	ChunkExecutor[T any] struct {
		executor  *PeriodicalExecutor[chunk[T]]
		container *chunkContainer[T]
	}

	chunkOptions struct {
		chunkSize     int
		flushInterval time.Duration
	}
)

// NewChunkExecutor returns a ChunkExecutor.
func NewChunkExecutor[T any](execute Execute[T], opts ...ChunkOption) *ChunkExecutor[T] {
	options := newChunkOptions()
	for _, opt := range opts {
		opt(&options)
	}

	container := &chunkContainer[T]{
		execute:      execute,
		maxChunkSize: options.chunkSize,
	}
	executor := &ChunkExecutor[T]{
		executor:  NewPeriodicalExecutor[chunk[T]](options.flushInterval, container),
		container: container,
	}

	return executor
}

// Add adds task with given chunk size into ce.
func (ce *ChunkExecutor[T]) Add(task T, size int) error {
	ce.executor.Add(chunk[T]{
		val:  task,
		size: size,
	})
	return nil
}

// Flush forces ce to flush and execute tasks.
func (ce *ChunkExecutor[T]) Flush() {
	ce.executor.Flush()
}

// Wait waits the execution to be done.
func (ce *ChunkExecutor[T]) Wait() {
	ce.executor.Wait()
}

// WithChunkBytes customizes a ChunkExecutor with the given chunk size.
func WithChunkBytes(size int) ChunkOption {
	return func(options *chunkOptions) {
		options.chunkSize = size
	}
}

// WithFlushInterval customizes a ChunkExecutor with the given flush interval.
func WithFlushInterval(duration time.Duration) ChunkOption {
	return func(options *chunkOptions) {
		options.flushInterval = duration
	}
}

func newChunkOptions() chunkOptions {
	return chunkOptions{
		chunkSize:     defaultChunkSize,
		flushInterval: defaultFlushInterval,
	}
}

type chunkContainer[T any] struct {
	tasks        []chunk[T]
	execute      Execute[T]
	size         int
	maxChunkSize int
}

func (bc *chunkContainer[T]) AddTask(task chunk[T]) bool {
	bc.tasks = append(bc.tasks, task)
	bc.size += task.size
	return bc.size >= bc.maxChunkSize
}

func (bc *chunkContainer[T]) Execute(tasks []chunk[T]) {
	vals := make([]T, 0, len(tasks))
	for _, elem := range tasks {
		vals = append(vals, elem.val)
	}
	bc.execute(vals)
}

func (bc *chunkContainer[T]) RemoveAll() []chunk[T] {
	tasks := bc.tasks
	bc.tasks = nil
	bc.size = 0
	return tasks
}

type chunk[T any] struct {
	val  T
	size int
}
