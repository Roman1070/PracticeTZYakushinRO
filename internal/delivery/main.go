package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	config "workers/internal"
	workers "workers/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	AppPortEnv = "APP_PORT"
)

type Workers interface {
	Queue(ctx context.Context, job workers.JobData)
	Status(ctx context.Context, id string) workers.JobStatus
	StatusAll(ctx context.Context) []workers.JobStatus
}

type Client struct {
	workers      Workers
	retriesCount uint16
}

func main() {
	cfg := config.MustLoad()

	workers := workers.New(cfg.Timeout, cfg.RetriesCount, cfg.WorkersCount)

	workers.Run(context.Background())

	client := Client{workers: workers, retriesCount: cfg.RetriesCount}

	router := mux.NewRouter()

	router.HandleFunc("/job/{job_id}", client.Queue).Methods(http.MethodPut)
	router.HandleFunc("/job", client.Queue).Methods(http.MethodPost)
	router.HandleFunc("/job/{job_id}", client.Status).Methods(http.MethodGet)
	router.HandleFunc("/jobs", client.StatusAll).Methods(http.MethodGet)

	http.ListenAndServe(os.Getenv(AppPortEnv), router)
}

func (c *Client) Queue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string
	id, ok := vars["job_id"]
	if !ok {
		id = uuid.NewString()
	}

	job, err := c.CreateJob(r, id)
	if err != nil {
		WriteError(w, err.Error())
		fmt.Printf("[Queue] error: %v", err.Error())
		return
	}

	c.workers.Queue(r.Context(), job)
	w.WriteHeader(http.StatusOK)
}

func (c *Client) Status(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["job_id"]

	status := c.workers.Status(r.Context(), id)

	result, err := json.Marshal(status)
	if err != nil {
		WriteError(w, err.Error())
		fmt.Printf("[Status] error: %v", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Client) StatusAll(w http.ResponseWriter, r *http.Request) {
	status := c.workers.StatusAll(r.Context())

	result, err := json.Marshal(status)
	if err != nil {
		WriteError(w, err.Error())
		fmt.Printf("[StatusAll] error: %v", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (c *Client) CreateJob(r *http.Request, jobID string) (workers.JobData, error) {
	type request struct {
		JobData        []byte `json:"jobData,omitempty"`
		Name           string `json:"name,omitempty"`
		Priority       uint16 `json:"priority,omitempty"`
		SimulatedError string `json:"error,omitempty"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if (err != nil) && (err != io.EOF) {
		fmt.Printf("[CreateJob] error: %v", err.Error())
		return workers.JobData{}, fmt.Errorf("[CreateJob] error: %v", err.Error())
	}

	var job workers.JobData

	if len(req.JobData) == 0 {
		job.JobData = []byte(fmt.Sprintf("job data for id %v", jobID))
	} else {
		job.JobData = req.JobData
	}

	if len(req.Name) == 0 {
		job.Name = fmt.Sprintf("job name for id %v", jobID)
	} else {
		job.Name = req.Name
	}

	if len(req.SimulatedError) > 0 {
		job.SimulatedError = errors.New(req.SimulatedError)
	}

	job.ID = jobID
	job.Priority = req.Priority
	job.RetriesCount = c.retriesCount

	return job, nil
}

type ErrorWrapper struct {
	Err string `json:"err"`
}

func WriteError(w http.ResponseWriter, err string) {
	errWrapper := ErrorWrapper{Err: err}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json, _ := json.Marshal(errWrapper)
	w.Write(json)
}
