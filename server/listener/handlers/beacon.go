package handlers

import (
	"Dur4nC2/misc/protobuf/implantpb"
	"Dur4nC2/server/db"
	beaconRepository "Dur4nC2/server/domain/beacon/repository/postgres"
	hostRepository "Dur4nC2/server/domain/host/repository/postgres"
	"Dur4nC2/server/domain/models"
	taskRepository "Dur4nC2/server/domain/task/repository/postgres"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"io/ioutil"
	"time"
	"unsafe"
)

func beaconRegisterHandler(data []byte) *implantpb.Envelope {
	beaconReg := &implantpb.BeaconRegister{}
	err := proto.Unmarshal(data, beaconReg)
	if err != nil {
		fmt.Printf("Error decoding beacon registration message: %s", err)
		return nil
	}

	beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
	hostRepo := hostRepository.NewPostgresHostRepository(db.Session())
	reponseEnvelope := &implantpb.Envelope{}
	//check if hosts exists
	host, err := hostRepo.ReadByHostname(beaconReg.Register.Hostname)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		//fmt.Printf("Host not found %s", err)
		newHost := &models.Host{
			Hostname:  beaconReg.Register.Hostname,
			OSVersion: beaconReg.Register.Os,
			Locale:    beaconReg.Register.Locale,
		}
		err = hostRepo.Create(newHost)
		host, err = hostRepo.ReadByHostname(beaconReg.Register.Hostname)
		//host = &newhost
		if err != nil {
			fmt.Println("Database write %s", err)
			return nil
		}
	}
	//check if beacon exists
	beacon, err := beaconRepo.Read(uuid.FromStringOrNil(beaconReg.ID))
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		beacon = models.Beacon{
			ID: uuid.FromStringOrNil(beaconReg.ID),
		}
		beacon.Name = beaconReg.Register.Name
		beacon.HostID = host.ID
		//beacon.Hostname = beaconReg.Register.Hostname
		//beacon.UUID = uuid.FromStringOrNil(beaconReg.Register.Uuid)
		beacon.Username = beaconReg.Register.Username
		beacon.OS = beaconReg.Register.Os
		beacon.Arch = beaconReg.Register.Arch
		beacon.PID = beaconReg.Register.Pid
		beacon.ReconnectInterval = beaconReg.Register.ReconnectInterval
		beacon.ActiveC2 = beaconReg.Register.ActiveC2
		// beacon.ConfigID = uuid.FromStringOrNil(beaconReg.Register.ConfigID)
		beacon.Locale = beaconReg.Register.Locale

		beacon.Interval = beaconReg.Interval
		beacon.Jitter = beaconReg.Jitter
		beacon.NextCheckin = time.Now().Unix() + beaconReg.NextCheckin

		_, err := beaconRepo.Create(beacon, host)
		beacons, err := beaconRepo.List()
		//TODO refactor this, race condition
		index := len(beacons) - 1
		reponseEnvelope.Data = []byte(beacons[index].ID.String())
		reponseEnvelope.Type = implantpb.MsgBeaconID

		if err != nil {
			fmt.Println("Database write %s", err)
		}
	} else {
		reponseEnvelope.Data = []byte(beacon.ID.String())
		reponseEnvelope.Type = implantpb.MsgBeaconID
	}

	return reponseEnvelope
}
func beaconTasksHandler(data []byte) *implantpb.Envelope {
	beaconTasks := &implantpb.BeaconTasks{}
	err := proto.Unmarshal(data, beaconTasks)
	if err != nil {
		fmt.Println("Error decoding beacon tasks message: %s", err)
		return nil
	}
	taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())

	// If the message contains tasks then process it as results
	// otherwise send the beacon any pending tasks. Currently, we
	// don't receive results and send pending tasks at the same
	// time. We only send pending tasks if the request is empty.
	// If we send the Beacon 0 tasks it should not respond at all.
	if 0 < len(beaconTasks.Tasks) {
		//fmt.Printf("[team-server] Beacon %s returned %d task result(s)\n", beaconTasks.ID, len(beaconTasks.Tasks))
		go beaconTaskResults(beaconTasks.ID, beaconTasks.Tasks)
		return nil
	}
	//fmt.Printf("[team-server] Beacon %s requested pending task(s)\n", beaconTasks.ID)
	// Retrieve beacon pending tasks ordered by their creation time.
	pendingTasks, err := taskRepo.ListPendingTasks(uuid.FromStringOrNil(beaconTasks.ID))
	//pendingTasks, err := beaconRepo.ListTasks()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("Beacon task database error: %s\n", err)
		return nil
	}
	tasks := []*implantpb.Envelope{}
	for _, pendingTask := range pendingTasks {
		envelope := &implantpb.Envelope{}
		err = proto.Unmarshal(pendingTask.Request, envelope)
		if err != nil {
			fmt.Println("Error decoding pending task: %s", err)
			continue
		}
		envelope.ID = pendingTask.EnvelopeID
		tasks = append(tasks, envelope)
		pendingTask.State = models.SENT
		pendingTask.SentAt = time.Now()
		//TODO create a update in the tasks repository
		err = db.Session().Model(&models.BeaconTask{}).Where(&models.BeaconTask{
			ID: pendingTask.ID,
		}).Updates(pendingTask).Error
		if err != nil {
			fmt.Println("Database error: %s", err)
		}
	}
	taskData, err := proto.Marshal(&implantpb.BeaconTasks{Tasks: tasks})
	if err != nil {
		fmt.Println("Error marshaling beacon tasks message: %s", err)
		return nil
	}
	//fmt.Printf("[team-server] Sending %d task(s) to beacon %s\n", len(pendingTasks), beaconTasks.ID)
	return &implantpb.Envelope{
		Type: implantpb.MsgBeaconTasks,
		Data: taskData,
	}
}
func beaconTaskResults(beaconID string, taskEnvelopes []*implantpb.Envelope) *implantpb.Envelope {
	taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())
	for _, envelope := range taskEnvelopes {
		//dbTask, err := taskRepo.ListPendingTasks(uuid.FromStringOrNil(whoamiResp.Response.BeaconID))
		dbTask, err := taskRepo.ListTaskByEnvelopeID(envelope.ID)
		if err != nil {
			fmt.Printf("Error finding db task: %s\n", err)
			continue
		}

		//TODO provisional way to show task results
		switch envelope.Type {
		case implantpb.MsgWhoamiResp:
			whoamiResponsepb := &implantpb.WhoamiResp{}
			proto.Unmarshal(envelope.Data, whoamiResponsepb)
			dbTask.Response = unsafe.Slice(unsafe.StringData(whoamiResponsepb.Whoami.Username), len(whoamiResponsepb.Whoami.Username))
			fmt.Printf("\n[team-server] Task %s result:\n%s\n", dbTask.ID, whoamiResponsepb.Whoami.Username)
		case implantpb.MsgRegisterExtensionResp:
			registerExtensionResponsepb := &implantpb.RegisterExtensionResp{}
			proto.Unmarshal(envelope.Data, registerExtensionResponsepb)
			if registerExtensionResponsepb.Response.Err != "" {
				fmt.Printf("\n[team-server] register extension response: %s", registerExtensionResponsepb.Response.Err)
			} else {
				fmt.Printf("\n[team-server] register extension response: successful")
			}
		case implantpb.MsgListExtensionResp:
			listExtensionResponsepb := &implantpb.ListExtensionsResp{}
			proto.Unmarshal(envelope.Data, listExtensionResponsepb)
			if len(listExtensionResponsepb.Names) == 0 {
				fmt.Printf("\n[team-server] list extension response: there are no extensions register")
			} else {
				fmt.Printf("\n[team-server] list extension response: %s", listExtensionResponsepb.Names)
			}
		case implantpb.MsgCallExtensionResp:
			callingExtensionResponsepb := &implantpb.CallExtensionResp{}
			proto.Unmarshal(envelope.Data, callingExtensionResponsepb)
			dbTask.Response = callingExtensionResponsepb.Output
			fmt.Printf("\n[team-server] Task %s result:\n%s\n", dbTask.ID, callingExtensionResponsepb.Output)
		case implantpb.MsgExecuteAssemblyResp:
			execAssemblyResponsepb := &implantpb.ExecuteAssemblyResp{}
			proto.Unmarshal(envelope.Data, execAssemblyResponsepb)
			dbTask.Response = execAssemblyResponsepb.Output
			fmt.Printf("\n[team-server] Task %s result:\n%s\n", dbTask.ID, execAssemblyResponsepb.Output)
		case implantpb.MsgExecuteShellcodeResp:
			execShellcodeResponsepb := &implantpb.ExecuteShellcodeResp{}
			proto.Unmarshal(envelope.Data, execShellcodeResponsepb)
			dbTask.Response = execShellcodeResponsepb.Output
			fmt.Printf("\n[team-server] Task %s result:\n%s\n", dbTask.ID, execShellcodeResponsepb.Output)
		case implantpb.MsgDownloadResp:
			downloadResponsepb := &implantpb.DownloadResp{}
			proto.Unmarshal(envelope.Data, downloadResponsepb)
			dbTask.Response = downloadResponsepb.File
			if downloadResponsepb.Response.Err != "" {
				fmt.Printf("\n[team-server] Task %s result:\nFile download error: \n", dbTask.ID)
				break
			}
			err := ioutil.WriteFile(downloadResponsepb.LocalPath, downloadResponsepb.File, 0644)
			if err != nil {
				fmt.Printf("\n[team-server] Task %s result:\nFile download error: \n", dbTask.ID)
			} else {
				fmt.Printf("\n[team-server] Task %s result:\nFile downloaded into: %s\n", dbTask.ID, downloadResponsepb.LocalPath)
			}
		case implantpb.MsgUploadResp:
			uploadResponsepb := &implantpb.UploadResp{}
			proto.Unmarshal(envelope.Data, uploadResponsepb)
			dbTask.Response = []byte(uploadResponsepb.Response.Err)
			if uploadResponsepb.Response.Err != "" {
				fmt.Println("Error writing file:", err)
				fmt.Printf("\n[team-server] Task %s result:\nFile upload error: %s\n", dbTask.ID, uploadResponsepb.Response.Err)
			} else {
				fmt.Printf("\n[team-server] Task %s result:\nFile upload into: %s\n", dbTask.ID)
			}
		}
		dbTask.State = models.COMPLETED
		dbTask.CompletedAt = time.Now()

		//TODO create update in task repository
		err = db.Session().Model(&models.BeaconTask{}).Where(&models.BeaconTask{
			ID: dbTask.ID,
		}).Updates(dbTask).Error
		if err != nil {
			fmt.Printf("Error updating db task: %s\n", err)
			continue
		}

	}
	return nil
}
