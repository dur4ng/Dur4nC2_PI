package extensions

import (
	"Dur4nC2/misc/protobuf/commonpb"
	"Dur4nC2/misc/protobuf/implantpb"
	"Dur4nC2/server/console"
	"Dur4nC2/server/console/command/settings"
	"Dur4nC2/server/db"
	beaconRepository "Dur4nC2/server/domain/beacon/repository/postgres"
	"Dur4nC2/server/domain/models"
	taskRepository "Dur4nC2/server/domain/task/repository/postgres"
	"errors"
	"github.com/desertbit/grumble"
	"github.com/gofrs/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"google.golang.org/protobuf/proto"
	"math/rand"
)

func ExternalExtensionsCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	beacon := con.ActiveTarget.GetBeacon()
	if beacon == nil {
		return errors.New("there is not an active beacon. Please, select one")
	}

	beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
	currentBeacon, err := beaconRepo.Read(uuid.FromStringOrNil(beacon.ID))
	if err != nil {
		return err
	}

	request := &commonpb.Request{BeaconID: beacon.ID}
	listExtensionReq := &implantpb.ListExtensionsReq{
		Request: request,
	}

	v := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
	data, _ := proto.Marshal(listExtensionReq)
	taskReq := &implantpb.Envelope{
		ID:   int64(v),
		Type: implantpb.MsgListExtensionReq,
		Data: data,
	}
	taskReqData, _ := proto.Marshal(taskReq)

	newTask := models.BeaconTask{
		EnvelopeID:  int64(v),
		BeaconID:    uuid.FromStringOrNil(beacon.ID),
		Description: "List register extensions",
		State:       models.PENDING,
		Request:     taskReqData,
	}
	taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())
	err = taskRepo.Create(currentBeacon, newTask)
	if err != nil {
		//con.PrintErrorf("%s\n", err)
		return err
	}
	con.PrintInfof("Task successfully added to the queue, please wait the response...")
	return nil
}
func ListExtensionsCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	if len(installedExtensions) == 0 {
		return errors.New("No installed extensions 🙁")
	}
	PrintExtensions(con)
	return nil
}
func PrintExtensions(con *console.ServerConsoleClient) {
	tw := table.NewWriter()
	tw.SetStyle(settings.GetTableStyle(con))
	tw.AppendHeader(table.Row{
		"Name",
		"Description",
	})
	tw.SortBy([]table.SortBy{
		{Name: "Name", Mode: table.Asc},
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 5, Align: text.AlignCenter},
	})

	for _, extension := range installedExtensions {
		tw.AppendRow(table.Row{
			extension.Name,
			extension.Description,
		})
	}
	con.Println(tw.Render())
}
