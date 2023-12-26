package tasks

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/server/console"
	"Dur4nC2/server/console/command/settings"
	taskUsecase "Dur4nC2/server/domain/task/usecase"
	"errors"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/term"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

func TasksCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	tasks, err := taskUsecase.NewTaskUsecase().List()
	if err != nil {
		return err
	}
	filter := ctx.Flags.String("filter")
	var filterRegex *regexp.Regexp
	if ctx.Flags.String("filter-re") != "" {
		var err error
		filterRegex, err = regexp.Compile(ctx.Flags.String("filter-re"))
		if err != nil {
			return err
		}
	}
	PrintTasks(tasks, filter, filterRegex, con)
	return nil
}

func ShowTaskCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	taskID := ctx.Args.String("task-id")
	if taskID == "" {
		return errors.New("non completed task")
	}
	task, err := taskUsecase.NewTaskUsecase().Read(&clientpb.BeaconTask{
		ID: taskID,
	})
	if err != nil {
		return err
	}
	if task.State != "completed" {
		return errors.New("non completed task")
	}

	con.Println(string(task.Response))
	filePath := ctx.Args.String("file-path")
	if filePath != "" {
		err := ioutil.WriteFile(filePath, task.Response, 0644)
		if err != nil {
			return errors.New("could not write the task response")
		}
	}
	return nil
}

func PrintTasks(tasks *clientpb.BeaconTasks, filter string, filterRegex *regexp.Regexp, con *console.ServerConsoleClient) {
	if len(tasks.Tasks) == 0 {
		con.PrintInfof("No tasks 🙁\n")
		return
	}
	tw := renderTasks(tasks, filter, filterRegex, con)
	con.Printf("%s\n", tw.Render())
}
func renderTasks(tasks *clientpb.BeaconTasks, filter string, filterRegex *regexp.Regexp, con *console.ServerConsoleClient) table.Writer {
	width, _, err := term.GetSize(0)
	if err != nil {
		width = 999
	}

	tw := table.NewWriter()
	tw.SetStyle(settings.GetTableStyle(con))
	wideTermWidth := 1000 < width
	if wideTermWidth {
		tw.AppendHeader(table.Row{
			"ID",
			"BeaconID",
			"CreatedAt",
			"State",
			"SentAt",
			"CompletedAt",
			"Description",
		})
	} else {
		tw.AppendHeader(table.Row{
			"ID",
			"CreatedAt",
			"State",
			"CompletedAt",
			"Description",
		})
	}

	for _, task := range tasks.Tasks {
		color := console.Normal
		activeBeacon := con.ActiveTarget.GetBeacon()
		if activeBeacon != nil && activeBeacon.ID == task.BeaconID {
			color = console.Green
		}

		// We need a slice of strings so we can apply filters
		var rowEntries []string

		if wideTermWidth {
			rowEntries = []string{
				//fmt.Sprintf(color+"%s"+console.Normal, strings.Split(beacon.ID, "-")[0]),
				fmt.Sprintf(color+"%s"+console.Normal, task.ID),
				fmt.Sprintf(color+"%s"+console.Normal, task.BeaconID),
				con.FormatDateDelta(time.Unix(task.CreatedAt, 0), wideTermWidth, false),
				fmt.Sprintf(color+"%s"+console.Normal, task.State),
				con.FormatDateDelta(time.Unix(task.SentAt, 0), wideTermWidth, false),
				con.FormatDateDelta(time.Unix(task.CompletedAt, 0), wideTermWidth, false),
				fmt.Sprintf(color+"%s"+console.Normal, task.Description),
			}
		} else {
			rowEntries = []string{
				//fmt.Sprintf(color+"%s"+console.Normal, strings.Split(beacon.ID, "-")[0]),
				fmt.Sprintf(color+"%s"+console.Normal, task.ID),
				con.FormatDateDelta(time.Unix(task.CreatedAt, 0), wideTermWidth, false),
				fmt.Sprintf(color+"%s"+console.Normal, task.State),
				con.FormatDateDelta(time.Unix(task.CompletedAt, 0), wideTermWidth, false),
				fmt.Sprintf(color+"%s"+console.Normal, task.Description),
			}
		}
		// Build the row struct
		row := table.Row{}
		for _, entry := range rowEntries {
			row = append(row, entry)
		}
		// Apply filters if any

		if filter == "" && filterRegex == nil {
			tw.AppendRow(row)
		} else {
			for _, rowEntry := range rowEntries {
				if filter != "" {
					if strings.Contains(rowEntry, filter) {
						tw.AppendRow(row)
						break
					}
				}
				if filterRegex != nil {
					if filterRegex.MatchString(rowEntry) {
						tw.AppendRow(row)
						break
					}
				}
			}
		}

	}
	return tw
}
