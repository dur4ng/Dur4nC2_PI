package postgresql

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/rpcpb"
	_beaconUsecase "Dur4nC2/server/domain/beacon/usecase"
	_hostUsecase "Dur4nC2/server/domain/host/usecase"
	_iocUsecase "Dur4nC2/server/domain/ioc/usecase"
	_lootUsecase "Dur4nC2/server/domain/loot/usecase"
	"Dur4nC2/server/domain/models"
	_operatorUsecase "Dur4nC2/server/domain/operator/usecase"
	_taskUsecase "Dur4nC2/server/domain/task/usecase"
	"context"
)

type RPCServer struct {
	rpcpb.UnimplementedTeamServerRPCServer
	operatorUsecase models.OperatorUsecase
	beaconUsecase   models.BeaconUsecase
	taskUsecase     models.BeaconTaskUsecase
	hostUsecase     models.HostUsecase
	iocUsecase      models.IOCUsecase
	lootUsecase     models.LootUsecase
}

func NewServer() *RPCServer {
	return &RPCServer{
		operatorUsecase: _operatorUsecase.NewOperatorUsecase(),
		beaconUsecase:   _beaconUsecase.NewBeaconUsecase(),
		taskUsecase:     _taskUsecase.NewTaskUsecase(),
		hostUsecase:     _hostUsecase.NewHostUsecase(),
		iocUsecase:      _iocUsecase.NewIOCUsecase(),
		lootUsecase:     _lootUsecase.NewLootUsecase(),
	}
}

// *** Listeners ***
func (s *RPCServer) StartHTTPSListener(ctx context.Context, in *clientpb.HTTPListenerReq) (*clientpb.HTTPListener, error) {
	return nil, nil
}
func (s *RPCServer) StartHTTPListener(ctx context.Context, in *clientpb.HTTPListenerReq) (*clientpb.HTTPListener, error) {
	return nil, nil
}

// *** Stager Listener ***
func (s *RPCServer) StartTCPStagerListener(ctx context.Context, in *clientpb.StagerListenerReq) (*clientpb.StagerListener, error) {
	return nil, nil
}
func (s *RPCServer) StartHTTPStagerListener(ctx context.Context, in *clientpb.StagerListenerReq) (*clientpb.StagerListener, error) {
	return nil, nil
}
