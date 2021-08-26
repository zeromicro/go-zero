package threading

// A WorkerGroup is used to run given number of workers to process jobs.
type WorkerGroup struct {
	job     func()
	workers int
}

// NewWorkerGroup returns a WorkerGroup with given job and workers.
func NewWorkerGroup(job func(), workers int) WorkerGroup {
	return WorkerGroup{
		job:     job,
		workers: workers,
	}
}

// Start starts a WorkerGroup.
func (wg WorkerGroup) Start() {
	group := NewRoutineGroup()
	for i := 0; i < wg.workers; i++ {
		group.RunSafe(wg.job)
	}
	group.Wait()
}
