package msgqueue

import "encoding/json"

type Envelope struct {
	JobId   string          `json:"job_id"`
	JobType string          `json:"job_type"`
	Payload json.RawMessage `json:"payload"`
}

func NewRawEnvelop(JobId, JobType string, Payload json.RawMessage) Envelope {
	return Envelope{
		JobId:   JobId,
		JobType: JobType,
		Payload: Payload,
	}
}
