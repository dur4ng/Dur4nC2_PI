package docker

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"context"
)

func (s *RPCServer) ListHosts(ctx context.Context, in *commonpb.Empty) (*clientpb.Hosts, error) {
	hosts, err := s.hostUsecase.List()
	if err != nil {
		return &clientpb.Hosts{}, err
	}
	return hosts, nil
}
func (s *RPCServer) GetHost(ctx context.Context, in *clientpb.Host) (*clientpb.Host, error) {
	host, err := s.hostUsecase.Get(in)
	if err != nil {
		return &clientpb.Host{}, err
	}
	return host, nil
}
func (s *RPCServer) DeleteHost(ctx context.Context, in *clientpb.Host) (*commonpb.Empty, error) {
	message, err := s.hostUsecase.Delete(in)
	if err != nil {
		return message, err
	}
	return message, nil
}
