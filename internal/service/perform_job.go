package workers

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

func PerformJob(id string, jobData []byte, simulatedError error) error {
	duration := time.Duration(rand.Int63n(int64(3 * time.Second)))
	fmt.Printf("Started job %s at %v, duration: %v\n", id, time.Now(), fmt.Sprint(duration))
	// do some work. Random sleep time; max = 3s
	<-time.After(duration)

	if simulatedError != nil {
		return simulatedError
	}

	return nil
}
