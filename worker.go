package job

import "log"

// Worker represents the worker that executes the model
type Worker struct {
	JobRunner  Task
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job, runner Task) Worker {
	return Worker{
		JobRunner:  runner,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				r, err := w.JobRunner.Run(job.Payload)
				if err != nil {
					log.Printf("Error running Job: %v", job.Payload)
					job.ReturnChannel <- err
				} else {
					job.ReturnChannel <- r
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
