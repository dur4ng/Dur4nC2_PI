package jobs

import (
	"Dur4nC2/server/console"
	"Dur4nC2/server/console/command/settings"
	"Dur4nC2/server/state"
	"errors"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/jedib0t/go-pretty/v6/table"
)

func JobsCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	if ctx.Flags.Int("kill") != -1 {
		err := KillJob(uint32(ctx.Flags.Int("kill")), con)
		return err
	} else if ctx.Flags.Bool("kill-all") {
		err := KillAllJobs(con)
		return err
	} else {
		err := GetAllJobs(con)
		return err
	}
}
func GetAllJobs(con *console.ServerConsoleClient) error {
	jobs := state.Jobs.All()
	if len(jobs) == 0 {
		con.PrintInfof("No jobs 🙁\n")
		return nil
	}
	PrintJobs(jobs, con)
	return nil
}
func KillJob(jobID uint32, con *console.ServerConsoleClient) error {
	con.PrintInfof("Killing job #%d ...\n", jobID)
	job := state.Jobs.Get(int(jobID))
	if job == nil {
		return errors.New("incorrect job id")
	} else {
		job.JobCtrl <- true
		con.PrintSuccessf("Successfully killed job #%d\n", job.ID)
		return nil
	}
}
func KillAllJobs(con *console.ServerConsoleClient) error {
	jobs := state.Jobs.All()
	if len(jobs) == 0 {
		con.PrintInfof("No jobs 🙁\n")
		return nil
	}
	for _, job := range jobs {
		if job == nil {
			return errors.New("incorrect job id")
		} else {
			job.JobCtrl <- true
			con.PrintSuccessf("Successfully killed job #%d\n", job.ID)
		}
	}
	return nil
}

// PrintJobs - Prints a list of active jobs
func PrintJobs(jobs []*state.Job, con *console.ServerConsoleClient) {

	tw := table.NewWriter()
	tw.SetStyle(settings.GetTableStyle(con))
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Protocol",
		"Port",
	})

	for _, job := range jobs {
		tw.AppendRow(table.Row{
			fmt.Sprintf("%d", job.ID),
			job.Name,
			job.Protocol,
			fmt.Sprintf("%d", job.Port),
		})
	}
	con.Printf("%s\n", tw.Render())
}
