package workers

import "context"

func (w *Workers) Run(ctx context.Context) {
	for i := 0; i <= int(w.WorkersCount); i++ {
		go w.Worker(ctx, i, w.CurrentJobs)
	}
}
