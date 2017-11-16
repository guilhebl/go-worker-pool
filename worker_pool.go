package worker

import (
	"log"
)

// Data a Job receives as Param
type Payload struct {
	JobType string
	Params  map[string]string
}

func NewPayload(jobType string, params map[string]string) Payload {
	return Payload{
		JobType: jobType,
		Params:  params,
	}
}

// Job represents the model to be run
type Job struct {
	Payload       Payload
	ReturnChannel chan interface{}
}

func NewJob(jobType string, params map[string]string, returnChannel chan interface{}) Job {
	return Job{
		Payload:       NewPayload(jobType, params),
		ReturnChannel: returnChannel,
	}
}

type JobRunner interface {
	Run(payload Payload) (interface{}, error)
}

// Worker represents the worker that executes the model
type Worker struct {
	JobRunner  JobRunner
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job, runner JobRunner) Worker {
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
				log.Printf("job req" )

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
