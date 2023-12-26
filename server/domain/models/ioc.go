package models

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

const (
	OPEN    = "open"
	PATCHED = "patched"
)

// IOC - Represents an indicator of compromise, generally a file we've
// uploaded to a remote system.
type IOC struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	HostID    uuid.UUID `gorm:"size:256"`
	CreatedAt time.Time `gorm:"->;<-:create;"`

	Path        string
	FileHash    string
	Name        string
	Description string
	Output      string
	State       string
}

func (i *IOC) BeforeCreate(tx *gorm.DB) (err error) {

	i.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}

	return nil
}

type IOCUsecase interface {
	Create(ioc_pb *clientpb.IOC) error
	List() ([]*clientpb.IOC, error)
	Read(id uuid.UUID) (*clientpb.IOC, error)
	Update(ioc_pb clientpb.IOC) error
}

type IOCRepository interface {
	Create(ioc IOC) error
	List() ([]IOC, error)
	Read(id uuid.UUID) (IOC, error)
	Update(ioc IOC) error
}
