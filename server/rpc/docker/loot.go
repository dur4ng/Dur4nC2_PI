package docker

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"context"
)

func (s *RPCServer) AddLoot(ctx context.Context, in *clientpb.Loot) (*clientpb.Loot, error) {
	newLoot, err := s.lootUsecase.Create(in)
	if err != nil {
		return nil, err
	}
	return newLoot, nil
}
func (s *RPCServer) ListLoot(ctx context.Context, in *commonpb.Empty) (*clientpb.Loots, error) {
	loots, err := s.lootUsecase.List()
	if err != nil {
		return nil, err
	}
	return loots, nil
}
func (s *RPCServer) GetLootContent(ctx context.Context, in *clientpb.Loot) (*clientpb.Loot, error) {
	loot, err := s.lootUsecase.Read(in)
	if err != nil {
		return nil, err
	}
	return loot, nil
}
func (s *RPCServer) GetLoot(ctx context.Context, in *clientpb.Loot) (*clientpb.Loots, error) {
	return nil, nil
}
func (s *RPCServer) DeleteLoot(ctx context.Context, in *clientpb.Loot) (*commonpb.Empty, error) {
	return nil, nil
}
func (s *RPCServer) UpdateLoot(ctx context.Context, in *clientpb.Loot) (*clientpb.Loot, error) {
	return nil, nil
}
func (s *RPCServer) DeleteHostLoot(ctx context.Context, in *clientpb.Loot) (*commonpb.Empty, error) {
	return nil, nil
}
