package workers

import (
	"context"
	"fmt"
)

func (w *Workers) Queue(ctx context.Context, job JobData) {
	fmt.Printf("[Queue] service started for job with id %v\n", job.ID)
	w.JobsStatuses.Store(job.ID, StatusQueued)
	w.CurrentJobs <- job
}
