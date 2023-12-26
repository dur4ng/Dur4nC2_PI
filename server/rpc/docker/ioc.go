package docker

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"context"
)

func (s *RPCServer) AddHostIOC(ctx context.Context, in *clientpb.IOC) (*commonpb.Empty, error) {
	err := s.iocUsecase.Create(in)
	if err != nil {
		return &commonpb.Empty{}, err
	}
	return &commonpb.Empty{}, nil
}
func (s *RPCServer) DeleteHostIOC(ctx context.Context, in *clientpb.IOC) (*commonpb.Empty, error) {
	return nil, nil
}
