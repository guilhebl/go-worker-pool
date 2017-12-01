package job

import (
	"strconv"
	"testing"
)

// Anonymous Task
type TestTask struct{}

func (e *TestTask) Run(payload Payload) JobResult {
	return NewJobResult("job done!", nil)
}

func NewTestTask() TestTask {
	return TestTask{}
}

// Anonymous Sum Task
type TestSumTask struct{}

func (e *TestSumTask) Run(payload Payload) JobResult {
	x, _ := strconv.ParseInt(payload.Params["x"], 10, 0)
	y, _ := strconv.ParseInt(payload.Params["y"], 10, 0)
	return NewJobResult(x+y, nil)
}

func NewTestSumTask() TestSumTask {
	return TestSumTask{}
}

func testClosePool(t *testing.T, p *WorkerPool) {
	// try to close pool
	ok := p.Stop()

	if ok {
		t.Error("WorkerPool not closed")
	}
}

// just creates a pool but doesn't run it tests open and close
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

	testClosePool(t, &p)
}

// tries to run a single worker job pool
func TestRunSingleJob(t *testing.T) {

	// create pool
	p := NewWorkerPool(1)
	jobQueue := make(chan Job)
	p.Run(jobQueue)

	// check running pool num workers created
	if len(p.Workers) != 1 {
		t.Error("Wrong number of Workers")
	}

	// let's create a test job
	ret := NewJobResultChannel()
	//defer close(ret)
	m := make(map[string]string)
	task := NewTestTask()
	work := NewJob(&task, m, ret)

	// Push the job onto the queue.
	jobQueue <- work

	// wait for response from Job
	resp := <-ret

	if resp.Error != nil || resp.Value.(string) != "job done!" {
		t.Error("Error while getting result of Job")
	}

	// try to close pool
	testClosePool(t, &p)
}

// tries to run multiple jobs
func TestRunMultipleJob(t *testing.T) {

	// create pool
	p := NewWorkerPool(4)
	jobQueue := make(chan Job)
	p.Run(jobQueue)

	// check running pool num workers created
	if len(p.Workers) != 4 {
		t.Error("Wrong number of Workers")
	}

	// let's create a test job
	ret := NewJobResultChannel()
	m := make(map[string]string)
	x, y := 1, 2
	m["x"] = strconv.Itoa(x)
	m["y"] = strconv.Itoa(y)
	task := NewTestSumTask()
	work := NewJob(&task, m, ret)

	// let's create a test job 2
	ret2 := NewJobResultChannel()
	m2 := make(map[string]string)
	x2, y2 := 3, 4
	m2["x"] = strconv.Itoa(x2)
	m2["y"] = strconv.Itoa(y2)
	task2 := NewTestSumTask()
	work2 := NewJob(&task2, m2, ret2)

	// let's create a test job 3
	ret3 := NewJobResultChannel()
	m3 := make(map[string]string)
	x3, y3 := 5, 6
	m3["x"] = strconv.Itoa(x3)
	m3["y"] = strconv.Itoa(y3)
	task3 := NewTestSumTask()
	work3 := NewJob(&task3, m3, ret3)

	// let's create a test job 4
	ret4 := NewJobResultChannel()
	m4 := make(map[string]string)
	x4, y4 := 5, 6
	m4["x"] = strconv.Itoa(x4)
	m4["y"] = strconv.Itoa(y4)
	task4 := NewTestSumTask()
	work4 := NewJob(&task4, m4, ret4)

	// Push each job onto the queue.
	jobQueue <- work
	jobQueue <- work2
	jobQueue <- work3
	jobQueue <- work4

	// Consume the merged output from all jobs and check matching sum result
	sum := int64(0)
	for n := range Merge(ret, ret2, ret3, ret4) {
		result := n.Value
		sum += result.(int64)
	}

	if sum != 32 {
		t.Error("error while summing all results from merged Jobs")
	}

	// try to close pool
	testClosePool(t, &p)
}
