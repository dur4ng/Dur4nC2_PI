package models

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// *** Loot types ***
const (
	LOOT_FILE       = 0
	LOOT_CREDENTIAL = 1
)

// *** File types ***
const (
	NO_FILE = 0
	BINARY  = 1
	TEXT    = 2
)

// *** Credential types ***
const (
	LSASS = 0
	NTDS  = 1
	OTHER = 2
)

type Loot struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	HostID    uuid.UUID `gorm:"size:256"`
	CreatedAt time.Time

	Name           string
	Type           int
	CredentialType int
	FileType       int
	FileName       string
	FileData       []byte
}

func (l *Loot) BeforeCreate(tx *gorm.DB) (err error) {

	l.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	return nil
}

type LootUsecase interface {
	Create(loot *clientpb.Loot) (*clientpb.Loot, error)
	List() (*clientpb.Loots, error)
	Read(loot *clientpb.Loot) (*clientpb.Loot, error)
}

type LootRepository interface {
	Create(loot Loot, host Host) error
	List() ([]Loot, error)
	Read(id uuid.UUID) (Loot, error)
}
