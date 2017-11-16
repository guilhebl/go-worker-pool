package worker

// responsible for controlling the runtime worker pool.
// it uses a buffered pool of workers in a Job/Worker pattern
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	maxWorkers int
}

func NewDispatcher(maxWorkers int) Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers}
}

func (d *Dispatcher) Run(queue chan Job, jobRunner JobRunner) {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, jobRunner)
		worker.Start()
	}

	go d.dispatch(queue)
}

func (d *Dispatcher) dispatch(jobQueue chan Job) {
	for {
		select {
		case job := <-jobQueue:
			// a model request has been received
			go func(job Job) {
				// try to obtain a worker model channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the model to the worker model channel
				jobChannel <- job
			}(job)
		}
	}
}
