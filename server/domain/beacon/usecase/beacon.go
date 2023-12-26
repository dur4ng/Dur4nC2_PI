package usecase

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"Dur4nC2/server/domain/models"
	"github.com/gofrs/uuid"
)

type beaconUsecase struct {
	repository models.BeaconRepository
}

func NewBeaconUsecase(repository models.BeaconRepository) models.BeaconUsecase {
	//b := _beaconRepository.NewMysqlBeaconRepository(db.Session())
	return &beaconUsecase{repository: repository}
}

func (u *beaconUsecase) List() (*clientpb.Beacons, error) {
	beacons, err := u.repository.List()
	if err != nil {
		return &clientpb.Beacons{}, nil
	}

	beacons_pb := &clientpb.Beacons{Beacons: []*clientpb.Beacon{}}
	for _, beacon := range beacons {
		beacons_pb.Beacons = append(beacons_pb.Beacons, &clientpb.Beacon{
			ID:   beacon.ID.String(),
			Name: beacon.Name,
			//Hostname:          beacon.Hostname,
			HostID:            beacon.HostID.String(),
			Username:          beacon.Username,
			OS:                beacon.OS,
			Arch:              beacon.Arch,
			Transport:         beacon.Transport,
			RemoteAddress:     beacon.RemoteAddress,
			PID:               beacon.PID,
			LastCheckin:       beacon.LastCheckin.Unix(),
			ActiveC2:          beacon.ActiveC2,
			ReconnectInterval: beacon.ReconnectInterval,
			Interval:          beacon.Interval,
			Jitter:            beacon.Jitter,
			NextCheckin:       beacon.NextCheckin,
			TasksCount:        int64(len(beacon.Tasks)),
			Locale:            beacon.Locale,
		})
	}
	return beacons_pb, nil
}
func (u *beaconUsecase) Read(beacon_pb *clientpb.Beacon) (*clientpb.Beacon, error) {
	id, _ := uuid.FromString(beacon_pb.ID)
	beacon, err := u.repository.Read(id)
	if err != nil {
		return nil, err
	}
	beacon_pb = &clientpb.Beacon{
		ID:   beacon.ID.String(),
		Name: beacon.Name,
		//Hostname:          beacon.Hostname,
		HostID:            beacon.HostID.String(),
		Username:          beacon.Username,
		OS:                beacon.OS,
		Arch:              beacon.Arch,
		Transport:         beacon.Transport,
		RemoteAddress:     beacon.RemoteAddress,
		PID:               beacon.PID,
		LastCheckin:       beacon.LastCheckin.Unix(),
		ActiveC2:          beacon.ActiveC2,
		ReconnectInterval: beacon.ReconnectInterval,
		Interval:          beacon.Interval,
		Jitter:            beacon.Jitter,
		NextCheckin:       beacon.NextCheckin,
		TasksCount:        int64(len(beacon.Tasks)),
		Locale:            beacon.Locale}
	return beacon_pb, nil
}
func (u *beaconUsecase) Delete(beacon_pb *clientpb.Beacon) (*commonpb.Empty, error) {
	id, _ := uuid.FromString(beacon_pb.ID)
	err := u.repository.Delete(id)
	if err != nil {
		return &commonpb.Empty{}, err
	}
	return &commonpb.Empty{}, nil
}
func (u *beaconUsecase) ListTasks(beacon_pb *clientpb.Beacon) (*clientpb.BeaconTasks, error) {
	id, _ := uuid.FromString(beacon_pb.ID)
	hostID, _ := uuid.FromString(beacon_pb.HostID)
	beacon := &models.Beacon{
		ID:     id,
		HostID: hostID,
	}
	tasks, err := u.repository.ListTasks(*beacon)
	if err != nil {
		return nil, err
	}
	tasks_pb := &clientpb.BeaconTasks{Tasks: []*clientpb.BeaconTask{}}
	for _, task := range tasks {
		tasks_pb.Tasks = append(tasks_pb.Tasks, &clientpb.BeaconTask{
			ID:          task.ID.String(),
			BeaconID:    task.BeaconID.String(),
			CreatedAt:   task.CreatedAt.Unix(),
			State:       task.State,
			SentAt:      task.SentAt.Unix(),
			CompletedAt: task.CompletedAt.Unix(),
			Request:     task.Request,
			Response:    task.Response,
			Description: task.Description,
		})
	}
	return tasks_pb, nil
}
