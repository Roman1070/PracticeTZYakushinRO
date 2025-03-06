package workers

import (
	"context"
	"fmt"
)

func (w *Workers) StatusAll(ctx context.Context) []JobStatus {
	fmt.Printf("[StatusAll] service started\n")
	w.JobsStatuses.mx.Lock()
	defer w.JobsStatuses.mx.Unlock()

	result := make([]JobStatus, len(w.JobsStatuses.m))
	insertedValues := 0
	for k, v := range w.JobsStatuses.m {
		result[insertedValues] = JobStatus{
			ID:     k,
			Status: v,
		}
		insertedValues++
	}

	return result
}
