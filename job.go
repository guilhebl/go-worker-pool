package job

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