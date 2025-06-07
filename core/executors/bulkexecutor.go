package executors

import "time"

const defaultBulkTasks = 1000

type (
	// BulkOption defines the method to customize a BulkExecutor.
	BulkOption func(options *bulkOptions)

	// A BulkExecutor is an executor that can execute tasks on either requirement meets:
	// 1. up to given size of tasks
	// 2. flush interval time elapsed
	BulkExecutor[T any] struct {
		executor  *PeriodicalExecutor[T]
		container *bulkContainer[T]
	}

	bulkOptions struct {
		cachedTasks   int
		flushInterval time.Duration
	}
)

// NewBulkExecutor returns a BulkExecutor.
func NewBulkExecutor[T any](execute Execute[T], opts ...BulkOption) *BulkExecutor[T] {
	options := newBulkOptions()
	for _, opt := range opts {
		opt(&options)
	}

	container := &bulkContainer[T]{
		execute:  execute,
		maxTasks: options.cachedTasks,
	}
	executor := &BulkExecutor[T]{
		executor:  NewPeriodicalExecutor[T](options.flushInterval, container),
		container: container,
	}

	return executor
}

// Add adds task into be.
func (be *BulkExecutor[T]) Add(task T) error {
	be.executor.Add(task)
	return nil
}

// Flush forces be to flush and execute tasks.
func (be *BulkExecutor[T]) Flush() {
	be.executor.Flush()
}

// Wait waits be to done with the task execution.
func (be *BulkExecutor[T]) Wait() {
	be.executor.Wait()
}

// WithBulkTasks customizes a BulkExecutor with given tasks limit.
func WithBulkTasks(tasks int) BulkOption {
	return func(options *bulkOptions) {
		options.cachedTasks = tasks
	}
}

// WithBulkInterval customizes a BulkExecutor with given flush interval.
func WithBulkInterval(duration time.Duration) BulkOption {
	return func(options *bulkOptions) {
		options.flushInterval = duration
	}
}

func newBulkOptions() bulkOptions {
	return bulkOptions{
		cachedTasks:   defaultBulkTasks,
		flushInterval: defaultFlushInterval,
	}
}

type bulkContainer[T any] struct {
	tasks    []T
	execute  Execute[T]
	maxTasks int
}

func (bc *bulkContainer[T]) AddTask(task T) bool {
	bc.tasks = append(bc.tasks, task)
	return len(bc.tasks) >= bc.maxTasks
}

func (bc *bulkContainer[T]) Execute(tasks []T) {
	bc.execute(tasks)
}

func (bc *bulkContainer[T]) RemoveAll() []T {
	tasks := bc.tasks
	bc.tasks = nil
	return tasks
}
