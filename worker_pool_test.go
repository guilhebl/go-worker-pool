package job

import (
	"testing"
	"log"
)

func TestCreatePool(t *testing.T) {

	p := NewWorkerPool(4)

	if &p == nil {
		t.Error("Error while creating Pool")
	}

	if p.WorkerPool == nil {
		t.Error("Error while creating Pool")
	}

	if p.maxWorkers != 4 {
		t.Error("Error while creating Pool")
	}

	if p.Workers == nil {
		t.Error("Error while creating Pool")
	}

	p.Stop()
}

// Anonymous Task
type TestJob struct{}

func (e *TestJob) Run(payload Payload) (interface{}, error) {
	log.Printf("Run %s", payload.JobType)
	return "job done!", nil
}

func NewTestJob() TestJob {
	return TestJob{}
}

func TestRunJobPool(t *testing.T) {

	// create pool
	p := NewWorkerPool(4)
	jobQueue := make(chan Job)
	runner := NewTestJob()
	p.Run(jobQueue, &runner)

	// check running pool num workers created
	if len(p.Workers) != 4 {
		t.Error("Wrong number of Workers")
	}

	// let's create a test job
	ret := make(chan interface{})
	defer close(ret)
	m := make(map[string]string)
	work := NewJob("test", m, ret)

	// Push the job onto the queue.
	jobQueue <- work

	// wait for response from Job
	resp := <-ret

	if resp.(string) != "job done!" {
		t.Error("Error while getting result of Job")
	}

	// try to close pool
	p.Stop()
}
