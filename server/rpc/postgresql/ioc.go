package postgresql

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"context"
)

func (s *RPCServer) AddHostIOC(ctx context.Context, in *clientpb.IOC) (*commonpb.Empty, error) {
	//err := s.iocUsecase.Create(&clientpb.Host{HostUUID: in.HostID}, in) // need add two parameters in rpc, is possible? No
	err := s.iocUsecase.Create(in)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (s *RPCServer) DeleteHostIOC(ctx context.Context, in *clientpb.IOC) (*commonpb.Empty, error) {
	return nil, nil
}
