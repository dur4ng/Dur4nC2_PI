package docker_test

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	_beaconRepo "Dur4nC2/server/domain/beacon/repository/postgres"
	_hostRepo "Dur4nC2/server/domain/host/repository/postgres"
	_iocRepo "Dur4nC2/server/domain/ioc/repository/postgres"
	_lootRepository "Dur4nC2/server/domain/loot/repository/postgres"
	"Dur4nC2/server/domain/models"
	_operatorRepo "Dur4nC2/server/domain/operator/repository/docker"
	_taskRepo "Dur4nC2/server/domain/task/repository/postgres"
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("gRPC", func() {
	BeforeEach(func() {
		// *** DB ***
		err := Db.AutoMigrate(
			&models.Operator{},
			&models.Host{},
			&models.Beacon{},
			&models.BeaconTask{},
			&models.IOC{},
			&models.Loot{},
		)
		Ω(err).To(Succeed())

		// *** Operator ***
		o := _operatorRepo.NewDockerOpertatorRepository(Db)
		err = o.Create(&models.Operator{Username: "Juan", Token: "1234"})
		Ω(err).To(Succeed())
		var operators []models.Operator
		operators, err = o.List()
		Ω(operators).To(HaveLen(1))

		// *** Host ***
		h := _hostRepo.NewPostgresHostRepository(Db)
		err = h.Create(&models.Host{
			Hostname:  "DC-01",
			OSVersion: "Windows Server 2019",
			Locale:    "USA",
		})
		Ω(err).To(Succeed())
		hosts, err := h.List()
		Ω(hosts).To(HaveLen(1))

		// *** Beacon ***
		b := _beaconRepo.NewPostgresBeaconRepository(Db)
		err = b.Create(
			models.Beacon{
				Name:              "beacon01",
				Username:          "Administrator",
				UID:               "asdf",
				GID:               "asdf",
				OS:                hosts[0].OSVersion,
				Arch:              "x64",
				Transport:         "transport",
				RemoteAddress:     "192.168.0.1",
				PID:               1234,
				LastCheckin:       time.Now(),
				ReconnectInterval: 3000,
				ActiveC2:          "true",
				Locale:            hosts[0].Locale,
				Interval:          3,
				Jitter:            5,
				NextCheckin:       3,
			},
			hosts[0])
		Ω(err).To(Succeed())
		beacons, err := b.List()
		Ω(beacons).To(HaveLen(1))

		// *** Task ***
		t := _taskRepo.NewPostgresBeaconTaskRepository(Db)
		err = t.Create(
			beacons[0],
			models.BeaconTask{
				BeaconID:    beacons[0].ID,
				State:       models.PENDING,
				SentAt:      time.Now(),
				Description: "test",
			})
		Ω(err).To(Succeed())
		tasks, err := t.List()
		Ω(err).To(Succeed())
		Ω(tasks).To(HaveLen(1))

		// *** IOC ***

		// *** Loot ***

	})
	Context("Operator", func() {
		It("List", func() {
			operators_pb, err := RPCClient.ListOperators(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(operators_pb.Operators).To(HaveLen(1))
		})
	})
	Context("Beacon", func() {
		It("ListBeacons", func() {
			beacons_pb, err := RPCClient.ListBeacons(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(beacons_pb.Beacons).To(HaveLen(1))
		})
		It("GetBeacon", func() {
			beacons_pb, err := RPCClient.ListBeacons(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(beacons_pb.Beacons).To(HaveLen(1))
			beacon_pb, err := RPCClient.GetBeacon(context.Background(), beacons_pb.Beacons[0])
			Ω(err).To(Succeed())
			Ω(beacon_pb.ID).Should(Equal(beacons_pb.Beacons[0].ID))
		})
		It("DeleteBeacon", func() {
			beacons_pb, err := RPCClient.ListBeacons(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(beacons_pb.Beacons).To(HaveLen(1))
			_, err = RPCClient.DeleteBeacon(context.Background(), beacons_pb.Beacons[0])
			Ω(err).To(Succeed())
			beacons_pb, err = RPCClient.ListBeacons(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(beacons_pb.Beacons).To(HaveLen(0))
		})
		It("GetBeaconTasks", func() {
			beacons_pb, err := RPCClient.ListBeacons(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(beacons_pb.Beacons).To(HaveLen(1))
			tasks_pb, err := RPCClient.GetBeaconTasks(context.Background(), beacons_pb.Beacons[0])
			Ω(err).To(Succeed())
			Ω(tasks_pb.Tasks).To(HaveLen(1))
		})
	})
	Context("Host", func() {
		It("ListHosts", func() {
			hosts_pb, err := RPCClient.ListHosts(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(hosts_pb.Hosts).To(HaveLen(1))
		})
		It("GetHost", func() {
			hosts_pb, err := RPCClient.ListHosts(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(hosts_pb.Hosts).To(HaveLen(1))
			host_pb, err := RPCClient.GetHost(context.Background(), hosts_pb.Hosts[0])
			Ω(err).To(Succeed())
			Ω(host_pb.Hostname).To(Equal(hosts_pb.Hosts[0].Hostname))
		})
		It("DeleteHost", func() {
			hosts_pb, err := RPCClient.ListHosts(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(hosts_pb.Hosts).To(HaveLen(1))
			_, err = RPCClient.DeleteHost(context.Background(), hosts_pb.Hosts[0])
			Ω(err).To(Succeed())
			hosts_pb, err = RPCClient.ListHosts(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(hosts_pb.Hosts).To(HaveLen(0))
		})
	})
	Context("IOC", func() {
		It("AddHostIOC", func() {
			hosts_pb, err := RPCClient.ListHosts(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(hosts_pb.Hosts).To(HaveLen(1))
			_, err = RPCClient.AddHostIOC(
				context.Background(),
				&clientpb.IOC{
					Path:     "C:\\Temp\\implant.exe",
					FileHash: "asdfasdfdf23q4",
					HostID:   hosts_pb.Hosts[0].HostUUID,
				})
			Ω(err).To(Succeed())
			i := _iocRepo.NewPostgresIOCRepository(Db)
			iocs, err := i.List()
			Ω(err).To(Succeed())
			Ω(iocs).To(HaveLen(1))
		})
		It("DeleteHostIOC", func() {

		})
	})
	Context("Loot", func() {
		It("AddLoot", func() {
			hosts_pb, err := RPCClient.ListHosts(context.Background(), &commonpb.Empty{})
			Ω(err).To(Succeed())
			Ω(hosts_pb.Hosts).To(HaveLen(1))
			_, err = RPCClient.AddLoot(
				context.Background(),
				&clientpb.Loot{
					HostID:     hosts_pb.Hosts[0].HostUUID,
					Name:       "Test loot",
					Type:       1,
					Credential: &clientpb.Credential{User: "dur4n", Password: "1234"},
					FileType:   0,
					File:       &commonpb.File{Name: "", Data: nil},
				})
			Ω(err).To(Succeed())
			i := _lootRepository.NewPostgresLootRepository(Db)
			loots, err := i.List()
			Ω(err).To(Succeed())
			Ω(loots).To(HaveLen(1))
		})
	})
})
