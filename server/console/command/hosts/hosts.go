package hosts

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/server/console"
	"Dur4nC2/server/console/command/settings"
	hostUsecase "Dur4nC2/server/domain/host/usecase"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/term"
	"regexp"
	"strings"
	"time"
)

func HostCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	hosts, err := hostUsecase.NewHostUsecase().List()
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
	PrintHosts(hosts, filter, filterRegex, con)
	return nil
}
func PrintHosts(hosts *clientpb.Hosts, filter string, filterRegex *regexp.Regexp, con *console.ServerConsoleClient) {
	if len(hosts.Hosts) == 0 {
		con.PrintInfof("No hosts 🙁\n")
		return
	}
	tw := renderHosts(hosts, filter, filterRegex, con)
	con.Printf("%s\n", tw.Render())
}
func renderHosts(hosts *clientpb.Hosts, filter string, filterRegex *regexp.Regexp, con *console.ServerConsoleClient) table.Writer {
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
			"CreatedAt",
			"Hostname",
			"OSVersion",
			"Locale",
		})
	} else {
		tw.AppendHeader(table.Row{
			"ID",
			"CreatedAt",
			"Hostname",
			"OSVersion",
			"Locale",
		})
	}

	for _, host := range hosts.Hosts {
		color := console.Normal

		// We need a slice of strings so we can apply filters
		var rowEntries []string

		if wideTermWidth {
			rowEntries = []string{
				fmt.Sprintf(color+"%s"+console.Normal, host.HostUUID),
				con.FormatDateDelta(time.Unix(host.FirstContact, 0), wideTermWidth, false),
				fmt.Sprintf(color+"%s"+console.Normal, host.Hostname),
				fmt.Sprintf(color+"%s"+console.Normal, host.OSVersion),
				fmt.Sprintf(color+"%s"+console.Normal, host.Locale),
			}
		} else {
			rowEntries = []string{
				fmt.Sprintf(color+"%s"+console.Normal, host.HostUUID),
				con.FormatDateDelta(time.Unix(host.FirstContact, 0), wideTermWidth, false),
				fmt.Sprintf(color+"%s"+console.Normal, host.Hostname),
				fmt.Sprintf(color+"%s"+console.Normal, host.OSVersion),
				fmt.Sprintf(color+"%s"+console.Normal, host.Locale),
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
