package docker

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"context"
)

func (s *RPCServer) ListOperators(ctx context.Context, in *commonpb.Empty) (*clientpb.Operators, error) {
	operators, err := s.operatorUsecase.List()
	if err != nil {
		return &clientpb.Operators{}, err
	}
	return operators, err
}
