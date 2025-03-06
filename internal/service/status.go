package workers

import (
	"context"
	"fmt"
)

var (
	StatusQueued            = "Queued"
	StatusRunning           = "Running"
	StatusOutOfTime         = "Out of time"
	StatusNotFound          = "Not found"
	StatusFinished          = "Finished"
	StatusFinishedWithError = "Finished with error"
)

type JobStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (w *Workers) Status(ctx context.Context, id string) JobStatus {
	fmt.Printf("[Status] service started for id %v\n", id)
	status, ok := w.JobsStatuses.Load(id)

	if !ok {
		return JobStatus{
			ID:     id,
			Status: StatusNotFound,
		}
	}
	return JobStatus{
		ID:     id,
		Status: fmt.Sprint(status),
	}
}
