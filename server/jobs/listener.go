package jobs

import (
	"Dur4nC2/server/console"
	"Dur4nC2/server/listener"
	"Dur4nC2/server/state"
	"fmt"
	"sync"
)

// StartHTTPListenerJob - Start a HTTP listener as a job
func StartHTTPListenerJob(conf *listener.HTTPListenerConfig, con *console.ServerConsoleClient) (*state.Job, error) {
	server, err := listener.StartHTTPListener(conf)
	if err != nil {
		return nil, err
	}
	name := "http"

	job := &state.Job{
		ID:          state.NextJobID(),
		Name:        name,
		Description: fmt.Sprintf("%s for domain %s", name, conf.Domain),
		Protocol:    "tcp",
		Port:        conf.LPort,
		JobCtrl:     make(chan bool),
		Domains:     []string{conf.Domain},
	}
	state.Jobs.Add(job)

	cleanup := func(err error) {
		server.Cleanup()
		state.Jobs.Remove(job)
	}
	once := &sync.Once{}

	go func() {
		var err error
		err = server.HTTPServer.ListenAndServe()
		if err != nil {
			once.Do(func() { cleanup(err) })
			job.JobCtrl <- true // Cleanup other goroutine
			con.PrintErrorf("Error starting job: %s\n", err.Error())
		}
	}()

	go func() {
		<-job.JobCtrl
		once.Do(func() {
			//fmt.Println("JobCtrl hit")
			cleanup(nil)
			con.PrintErrorf("Cleaned job: %d", job.ID)
		})
	}()

	return job, nil
}
