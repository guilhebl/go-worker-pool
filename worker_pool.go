package job

import "sync"

// manages a pool of goroutines.
// it uses a buffered pool of workers in a Job/Worker pattern
type WorkerPool struct {
	// A pool of worker channels that are registered with the pool
	WorkerPool chan chan Job
	Workers []Worker
	maxWorkers int
}

func NewWorkerPool(maxWorkers int) WorkerPool {
	pool := make(chan chan Job, maxWorkers)
	workers := make([]Worker, 0)
	return WorkerPool{
		WorkerPool: pool,
		Workers: workers,
		maxWorkers: maxWorkers}
}

// Starts the WorkerPool
func (p *WorkerPool) Run(queue chan Job) {
	// starting n number of workers
	for i := 0; i < p.maxWorkers; i++ {
		worker := NewWorker(p.WorkerPool)
		p.Workers = append(p.Workers, worker)
		worker.Start()
	}

	go p.dispatch(queue)
}

// stops the Pool
func (p *WorkerPool) Stop() {
	// stops all workers
	for _, worker := range p.Workers {
		worker.Stop()
	}

	// close the Job pool chan
	close(p.WorkerPool)
}

// dispatches a job to be handled by an idle Worker of the pool
func (p *WorkerPool) dispatch(jobQueue chan Job) {
	for {
		select {
		case job := <-jobQueue:
			// a model request has been received
			go func(job Job) {
				// try to obtain a worker model channel that is available.
				// this will block until a worker is idle
				jobChannel := <-p.WorkerPool

				// dispatch the model to the worker model channel
				jobChannel <- job
			}(job)
		}
	}
}

// Utility function to merge multiple JobResult output channels into one
func Merge(cs ...<-chan JobResult) <-chan JobResult {
	var wg sync.WaitGroup
	out := make(chan JobResult)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan JobResult) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
