package state

import (
	"sync"
)

var (
	// Jobs - Holds pointers to all the current jobs
	Jobs = &jobs{
		// ID -> *Job
		active: &sync.Map{},
	}
	jobID = 0
)

// Job - Manages background jobs
type Job struct {
	ID           int
	Name         string
	Description  string
	Protocol     string
	Port         uint16
	Domains      []string
	JobCtrl      chan bool
	PersistentID string
}

// jobs - Holds refs to all active jobs
type jobs struct {
	active *sync.Map
}

// All - Return a list of all jobs
func (j *jobs) All() []*Job {
	all := []*Job{}
	j.active.Range(func(key, value interface{}) bool {
		all = append(all, value.(*Job))
		return true
	})
	return all
}

// Add - Add a job to the hive (atomically)
func (j *jobs) Add(job *Job) {
	j.active.Store(job.ID, job)
}

// Remove - Remove a job
func (j *jobs) Remove(job *Job) {
	j.active.LoadAndDelete(job.ID)
}

// Get - Get a Job
func (j *jobs) Get(jobID int) *Job {
	if jobID <= 0 {
		return nil
	}
	val, ok := j.active.Load(jobID)
	if ok {
		return val.(*Job)
	}
	return nil
}

// NextJobID - Returns an incremental nonce as an id
func NextJobID() int {
	newID := jobID + 1
	jobID++
	return newID
}
