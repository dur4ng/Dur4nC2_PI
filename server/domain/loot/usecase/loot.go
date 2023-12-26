package usecase

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"Dur4nC2/server/domain/models"
	"github.com/gofrs/uuid"
)

type lootUsecase struct {
	repository models.LootRepository
}

func NewLootUsecase(repository models.LootRepository) models.LootUsecase {
	return &lootUsecase{repository: repository}
}

func (u *lootUsecase) Create(loot_pb *clientpb.Loot) (*clientpb.Loot, error) {
	hostID, _ := uuid.FromString(loot_pb.HostID)
	host := models.Host{
		ID: hostID,
	}
	loot := models.Loot{
		HostID:         hostID,
		Name:           loot_pb.Name,
		Type:           int(loot_pb.Type),
		CredentialType: int(loot_pb.CredentialType),
		FileType:       int(loot_pb.FileType),
		FileName:       loot_pb.File.Name,
		FileData:       loot_pb.File.Data,
	}
	err := u.repository.Create(loot, host)
	if err != nil {
		return nil, err
	}
	return loot_pb, nil
}
func (u *lootUsecase) List() (*clientpb.Loots, error) {
	loots, err := u.repository.List()
	if err != nil {
		return nil, err
	}
	loots_pb := &clientpb.Loots{Loot: []*clientpb.Loot{}}
	for _, loot := range loots {
		loots_pb.Loot = append(loots_pb.Loot, &clientpb.Loot{
			LootID:         loot.ID.String(),
			HostID:         loot.HostID.String(),
			Name:           loot.Name,
			Type:           clientpb.LootType(loot.Type),
			CredentialType: clientpb.CredentialType(loot.CredentialType),
			FileType:       clientpb.FileType(loot.FileType),
			File:           &commonpb.File{Name: loot.FileName, Data: loot.FileData},
		})
	}
	return loots_pb, nil
}
func (u *lootUsecase) Read(loot_pb *clientpb.Loot) (*clientpb.Loot, error) {
	id, _ := uuid.FromString(loot_pb.LootID)
	loot, err := u.repository.Read(id)
	if err != nil {
		return nil, err
	}
	loot_pb = &clientpb.Loot{
		LootID:         loot.ID.String(),
		HostID:         loot.HostID.String(),
		Name:           loot.Name,
		Type:           clientpb.LootType(loot.Type),
		CredentialType: clientpb.CredentialType(loot.CredentialType),
		FileType:       clientpb.FileType(loot.FileType),
		File:           &commonpb.File{Name: loot.FileName, Data: loot.FileData},
	}
	return loot_pb, nil
}
