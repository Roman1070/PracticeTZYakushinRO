package workers

import (
	"errors"
	"sync"
)

type Workers struct {
	Timeout      float32
	RetriesCount uint16
	WorkersCount uint16
	CurrentJobs  chan JobData
	JobsStatuses JobsMap
}

type JobData struct {
	ID             string
	Name           string
	JobData        []byte
	Priority       uint16
	RetriesCount   uint16
	SimulatedError error
}

func New(timeout float32, retriesCount, workersCount uint16) *Workers {
	return &Workers{
		Timeout:      timeout,
		RetriesCount: retriesCount,
		WorkersCount: workersCount,
		CurrentJobs:  make(chan JobData, 100),
		JobsStatuses: JobsMap{
			mx: sync.RWMutex{},
			m:  map[string]string{},
		},
	}
}

var (
	ErrorOutOfTime = errors.New("timeout exceeded")
)
