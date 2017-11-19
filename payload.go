package job

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