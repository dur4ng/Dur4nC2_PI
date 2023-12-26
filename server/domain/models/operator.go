package models

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Operator struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Username  string
	CreatedAt time.Time
	Token     string
}

// GORM Hooks (triggers)
func (o *Operator) BeforeCreate(tx *gorm.DB) (err error) {

	o.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}

	return nil
}

type OperatorUsecase interface {
	List() (*clientpb.Operators, error)
	IsOnline(operatorID uuid.UUID) bool
}

type OperatorRepository interface {
	Create(operator *Operator) error
	List() ([]Operator, error)
	Read(id uuid.UUID) (Operator, error)
	Update(operator *Operator) error
	Delete(id uuid.UUID) error
}
