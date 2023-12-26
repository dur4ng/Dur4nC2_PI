package docker

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"context"
)

func (s *RPCServer) ListBeacons(ctx context.Context, in *commonpb.Empty) (*clientpb.Beacons, error) {
	beacons, err := s.beaconUsecase.List()
	if err != nil {
		return &clientpb.Beacons{}, err
	}
	return beacons, nil
}
func (s *RPCServer) GetBeacon(ctx context.Context, in *clientpb.Beacon) (*clientpb.Beacon, error) {
	beacon, err := s.beaconUsecase.Read(in)
	if err != nil {
		return nil, err
	}
	return beacon, nil
}
func (s *RPCServer) DeleteBeacon(ctx context.Context, in *clientpb.Beacon) (*commonpb.Empty, error) {
	message, err := s.beaconUsecase.Delete(in)
	if err != nil {
		return message, err
	}
	return message, nil
}

func (s *RPCServer) GetBeaconTasks(ctx context.Context, in *clientpb.Beacon) (*clientpb.BeaconTasks, error) {
	tasks, err := s.beaconUsecase.ListTasks(in)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
