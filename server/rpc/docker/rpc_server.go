package docker

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/rpcpb"
	_beaconRepository "Dur4nC2/server/domain/beacon/repository/postgres"
	_beaconUsecase "Dur4nC2/server/domain/beacon/usecase"
	_hostRepository "Dur4nC2/server/domain/host/repository/postgres"
	_hostUsecase "Dur4nC2/server/domain/host/usecase"
	_iocRepo "Dur4nC2/server/domain/ioc/repository/postgres"
	_iocUsecase "Dur4nC2/server/domain/ioc/usecase"
	_lootRepository "Dur4nC2/server/domain/loot/repository/postgres"
	_lootUsecase "Dur4nC2/server/domain/loot/usecase"
	"Dur4nC2/server/domain/models"
	_operatorRepository "Dur4nC2/server/domain/operator/repository/docker"
	_operatorUsecase "Dur4nC2/server/domain/operator/usecase/docker"
	_taskUsecase "Dur4nC2/server/domain/task/usecase"
	"context"
	"gorm.io/gorm"
)

type RPCServer struct {
	rpcpb.UnimplementedTeamServerRPCServer

	operatorUsecase    models.OperatorUsecase
	operatorDestructor func()

	beaconUsecase models.BeaconUsecase
	taskUsecase   models.BeaconTaskUsecase
	hostUsecase   models.HostUsecase
	iocUsecase    models.IOCUsecase
	lootUsecase   models.LootUsecase
}

func NewServer(conn *gorm.DB) *RPCServer {
	return &RPCServer{
		operatorUsecase: _operatorUsecase.NewOperatorUsecase(_operatorRepository.NewDockerOpertatorRepository(conn)),
		beaconUsecase:   _beaconUsecase.NewBeaconUsecase(_beaconRepository.NewPostgresBeaconRepository(conn)),
		taskUsecase:     _taskUsecase.NewTaskUsecase(),
		hostUsecase:     _hostUsecase.NewHostUsecase(_hostRepository.NewPostgresHostRepository(conn)),
		iocUsecase:      _iocUsecase.NewIOCUsecase(_iocRepo.NewPostgresIOCRepository(conn)),
		lootUsecase:     _lootUsecase.NewLootUsecase(_lootRepository.NewPostgresLootRepository(conn)),
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
