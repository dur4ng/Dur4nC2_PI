package beacons

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/server/console"
	"Dur4nC2/server/console/command/settings"
	"Dur4nC2/server/db"
	"Dur4nC2/server/domain/beacon/repository/postgres"
	beaconRepository "Dur4nC2/server/domain/beacon/repository/postgres"
	"Dur4nC2/server/domain/beacon/usecase"
	hostRepository "Dur4nC2/server/domain/host/repository/postgres"
	"Dur4nC2/server/domain/models"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/gofrs/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/term"
	"regexp"
	"strings"
	"time"
)

func BeaconsCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	killFlag := ctx.Flags.String("kill")
	killAll := ctx.Flags.Bool("kill-all")

	// Handle kill
	if killFlag != "" {
		beaconUUID, err := uuid.FromString(killFlag)
		if err != nil {
			return err
		}
		err = beaconRepository.NewPostgresBeaconRepository(db.Session()).Delete(beaconUUID)
		if err != nil {
			return err
		}
		if con.ActiveTarget.GetBeacon() != nil && killFlag == con.ActiveTarget.GetBeacon().ID {
			con.ActiveTarget.Background()
		}

		con.PrintSuccessf("Successfully killed beacon %s\n", killFlag)
		con.Println()

	}
	if killAll {
		beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
		err := beaconRepo.DeleteAll()
		if err != nil {
			return err
		}
		con.ActiveTarget.Background()
		con.PrintSuccessf("all beacons were killed")
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

	beaconRepo := postgres.NewPostgresBeaconRepository(db.Session())
	beaconUsecase := usecase.NewBeaconUsecase(beaconRepo)
	beacons, err := beaconUsecase.List()
	if err != nil {
		return err
	}
	PrintBeacons(beacons.Beacons, filter, filterRegex, con)
	return nil
}
func PrintBeacons(beacons []*clientpb.Beacon, filter string, filterRegex *regexp.Regexp, con *console.ServerConsoleClient) {
	if len(beacons) == 0 {
		con.PrintInfof("No beacons 🙁\n")
		return
	}
	tw := renderBeacons(beacons, filter, filterRegex, con)
	con.Printf("%s\n", tw.Render())
}
func renderBeacons(beacons []*clientpb.Beacon, filter string, filterRegex *regexp.Regexp, con *console.ServerConsoleClient) table.Writer {
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
			"Name",
			"Tasks",
			"Transport",
			"Remote Address",
			"Hostname",
			"Username",
			"Operating System",
			"Locale",
			"Last Check-in",
			"Next Check-in",
		})
	} else {
		tw.AppendHeader(table.Row{
			"ID",
			"Name",
			"Transport",
			"Username",
			"Operating System",
			"Last Check-in",
			"Next Check-in",
		})
	}

	for _, beacon := range beacons {
		color := console.Normal
		activeBeacon := con.ActiveTarget.GetBeacon()
		if activeBeacon != nil && activeBeacon.ID == beacon.ID {
			color = console.Green
		}
		hostRepo := hostRepository.NewPostgresHostRepository(db.Session())
		host, err := hostRepo.Read(uuid.FromStringOrNil(beacon.HostID))
		if err != nil {
			host = models.Host{ID: uuid.FromStringOrNil("..."), Hostname: "..."}
		}
		// We need a slice of strings so we can apply filters
		var rowEntries []string

		if wideTermWidth {
			rowEntries = []string{
				//fmt.Sprintf(color+"%s"+console.Normal, strings.Split(beacon.ID, "-")[0]),
				fmt.Sprintf(color+"%s"+console.Normal, beacon.ID),
				fmt.Sprintf(color+"%s"+console.Normal, beacon.Name),
				fmt.Sprintf(color+"%d/%d"+console.Normal, beacon.TasksCountCompleted, beacon.TasksCount),
				fmt.Sprintf(color+"%s"+console.Normal, beacon.Transport),
				fmt.Sprintf(color+"%s"+console.Normal, beacon.RemoteAddress),
				//fmt.Sprintf(color+"%s"+console.Normal, beacon.Hostname),
				fmt.Sprintf(color+"%s"+console.Normal, host.Hostname),
				//fmt.Sprintf(color+"%s"+console.Normal, strings.TrimPrefix(beacon.Username, beacon.Hostname+"\\")),
				fmt.Sprintf(color+"%s"+console.Normal, strings.TrimPrefix(beacon.Username, host.Hostname+"\\")),
				fmt.Sprintf(color+"%s/%s"+console.Normal, beacon.OS, beacon.Arch),
				fmt.Sprintf(color+"%s"+console.Normal, beacon.Locale),
				con.FormatDateDelta(time.Unix(beacon.LastCheckin, 0), wideTermWidth, false),
				con.FormatDateDelta(time.Unix(beacon.NextCheckin, 0), wideTermWidth, true),
			}
		} else {
			rowEntries = []string{
				//fmt.Sprintf(color+"%s"+console.Normal, strings.Split(beacon.ID, "-")[0]),
				fmt.Sprintf(color+"%s"+console.Normal, beacon.ID),
				fmt.Sprintf(color+"%s"+console.Normal, beacon.Name),
				fmt.Sprintf(color+"%s"+console.Normal, beacon.Transport),
				fmt.Sprintf(color+"%s"+console.Normal, strings.TrimPrefix(beacon.Username, host.Hostname+"\\")),
				fmt.Sprintf(color+"%s/%s"+console.Normal, beacon.OS, beacon.Arch),
				con.FormatDateDelta(time.Unix(beacon.LastCheckin, 0), wideTermWidth, false),
				con.FormatDateDelta(time.Unix(beacon.NextCheckin, 0), wideTermWidth, true),
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
