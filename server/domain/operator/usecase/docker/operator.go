package docker

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/server/domain/models"
	"github.com/gofrs/uuid"
)

type operatorUsecase struct {
	repository models.OperatorRepository
}

func NewOperatorUsecase(repository models.OperatorRepository) models.OperatorUsecase {
	return &operatorUsecase{repository: repository}
}

func (u *operatorUsecase) List() (*clientpb.Operators, error) {
	operators, err := u.repository.List()
	if err != nil {
		return &clientpb.Operators{}, err
	}

	operators_pb := &clientpb.Operators{Operators: []*clientpb.Operator{}}
	for _, o := range operators {
		operators_pb.Operators = append(operators_pb.Operators, &clientpb.Operator{Name: o.Username, Online: u.IsOnline(o.ID)})
	}

	return operators_pb, nil
}

func (u *operatorUsecase) IsOnline(operatorID uuid.UUID) bool {
	return false
}
