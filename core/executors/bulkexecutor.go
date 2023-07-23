package executors

import "time"

const defaultBulkTasks = 1000

type (
	// BulkOption defines the method to customize a BulkExecutor.
	BulkOption func(options *bulkOptions)

	// A BulkExecutor is an executor that can execute tasks on either requirement meets:
	// 1. up to given size of tasks
	// 2. flush interval time elapsed
	BulkExecutor struct {
		executor  *PeriodicalExecutor
		container *bulkContainer
	}

	bulkOptions struct {
		cachedTasks   int
		flushInterval time.Duration
	}
)

// NewBulkExecutor returns a BulkExecutor.
func NewBulkExecutor(execute Execute, opts ...BulkOption) *BulkExecutor {
	options := newBulkOptions()
	for _, opt := range opts {
		opt(&options)
	}

	container := &bulkContainer{
		execute:  execute,
		maxTasks: options.cachedTasks,
	}
	executor := &BulkExecutor{
		executor:  NewPeriodicalExecutor(options.flushInterval, container),
		container: container,
	}

	return executor
}

// Add adds task into be.
func (be *BulkExecutor) Add(task any) error {
	be.executor.Add(task)
	return nil
}

// Flush forces be to flush and execute tasks.
func (be *BulkExecutor) Flush() {
	be.executor.Flush()
}

// Wait waits be to done with the task execution.
func (be *BulkExecutor) Wait() {
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

type bulkContainer struct {
	tasks    []any
	execute  Execute
	maxTasks int
}

func (bc *bulkContainer) AddTask(task any) bool {
	bc.tasks = append(bc.tasks, task)
	return len(bc.tasks) >= bc.maxTasks
}

func (bc *bulkContainer) Execute(tasks any) {
	vals := tasks.([]any)
	bc.execute(vals)
}

func (bc *bulkContainer) RemoveAll() any {
	tasks := bc.tasks
	bc.tasks = nil
	return tasks
}
