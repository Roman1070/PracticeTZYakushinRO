package workers

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Job func(id string, jobData []byte, simulatedError error) error

func (w *Workers) Worker(ctx context.Context, id int, jobs chan JobData) {
	for job := range jobs {
		fmt.Printf("Worker %d processing task %s, \n", id, job.Name)

		w.JobsStatuses.Store(job.ID, StatusRunning)

		err := doJobWithTimeout(ctx, time.Duration(w.Timeout*float32(time.Second)), PerformJob, job)

		if err != nil {
			handleError(job, jobs, w, err)
		} else {
			w.JobsStatuses.Store(job.ID, StatusFinished)
		}
	}
}

func handleError(job JobData, jobs chan JobData, w *Workers, err error) {
	if job.RetriesCount > 1 {
		job.RetriesCount--

		//simulating external job error fix
		if rand.Int31n(2) == 0 {
			job.SimulatedError = nil
		}

		fmt.Printf("Retrying task %s, err: %v, retries left: %d\n", job.Name, err.Error(), job.RetriesCount)
		jobs <- job
	} else {
		fmt.Printf("Task %s failed after retries, err: %v\n", job.Name, err.Error())
		w.JobsStatuses.Store(job.ID, StatusFinishedWithError)
	}
}

func doJobWithTimeout(ctx context.Context, timeout time.Duration, job Job, jobData JobData) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	done := make(chan struct{}, 1)

	var err error
	go func() {
		err = job(jobData.ID, jobData.JobData, jobData.SimulatedError)
		done <- struct{}{}
	}()

	select {
	case <-timeoutCtx.Done():
		return ErrorOutOfTime
	case <-done:
		if err != nil {
			fmt.Printf("Job %v finished with error %v\n", jobData.ID, err.Error())
		} else {
			fmt.Printf("Job %v finished successfully!\n", jobData.ID)
		}
		return err
	}
}
