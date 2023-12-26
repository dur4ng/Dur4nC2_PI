package postgresql

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"context"
)

func (s *RPCServer) GetBeaconTaskContent(ctx context.Context, in *clientpb.BeaconTask) (*clientpb.BeaconTask, error) {
	task, err := s.taskUsecase.Read(in)
	if err != nil {
		return nil, err
	}
	return task, nil
}
func (s *RPCServer) CancelBeaconTask(ctx context.Context, in *clientpb.BeaconTask) (*clientpb.BeaconTask, error) {
	canceledTask, err := s.taskUsecase.Cancel(in)
	if err != nil {
		return nil, err
	}
	return canceledTask, nil
}
