package usecase

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/server/domain/models"
	"fmt"
	"github.com/gofrs/uuid"
)

type iocUsecase struct {
	repository models.IOCRepository
}

func NewIOCUsecase(repository models.IOCRepository) models.IOCUsecase {
	return &iocUsecase{repository: repository}
}

func (u *iocUsecase) Create(ioc_pb *clientpb.IOC) error {
	hostID, _ := uuid.FromString(ioc_pb.HostID)
	ioc := models.IOC{
		HostID:   hostID,
		Path:     ioc_pb.Path,
		FileHash: ioc_pb.FileHash,
	}
	err := u.repository.Create(ioc)
	if err != nil {
		fmt.Println("Usecase IOC ERROR...")
		return err
	}
	return nil
}
func (u *iocUsecase) List() ([]*clientpb.IOC, error) {
	iocs, err := u.repository.List()
	if err != nil {
		return nil, err
	}
	var iocs_pb []*clientpb.IOC
	for _, ioc := range iocs {
		iocs_pb = append(iocs_pb, &clientpb.IOC{
			ID:       ioc.ID.String(),
			FileHash: ioc.FileHash,
			Path:     ioc.Path,
		})
	}
	return iocs_pb, nil
}
func (u *iocUsecase) Read(id uuid.UUID) (*clientpb.IOC, error) {
	ioc, err := u.repository.Read(id)
	if err != nil {
		return nil, err
	}
	ioc_pb := &clientpb.IOC{
		ID:       ioc.ID.String(),
		FileHash: ioc.FileHash,
		Path:     ioc.Path,
	}
	return ioc_pb, nil
}
func (u *iocUsecase) Update(ioc_pb clientpb.IOC) error {
	return nil
}
