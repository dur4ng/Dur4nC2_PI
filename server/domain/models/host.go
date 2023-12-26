package models

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"gorm.io/gorm"
	"time"

	"github.com/gofrs/uuid"
)

// Host - Represents a host machine
type Host struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	Deleted gorm.DeletedAt
	//HostID  uuid.UUID //`gorm:"type:uuid;"` has one gorm association
	//BeaconID  uuid.UUID `gorm:"size:256"`
	CreatedAt time.Time `gorm:"->;<-:create;"`

	Hostname  string
	OSVersion string // Verbore OS version
	Locale    string // Detected language code

	Beacons []Beacon
	IOCs    []IOC
	Loots   []Loot
}

func (h *Host) BeforeCreate(tx *gorm.DB) (err error) {

	h.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	return nil
}

type HostUsecase interface {
	List() (*clientpb.Hosts, error)
	Get(host_pb *clientpb.Host) (*clientpb.Host, error)
	Delete(host_pb *clientpb.Host) (*commonpb.Empty, error)

	//ListBeacons(host Host) ([]Host, error)
	ListHostIOC(host_pb clientpb.Host) ([]*clientpb.IOC, error)
	ListHostLoot(loot_pb *clientpb.Loot) (*clientpb.Loots, error)
}

type HostRepository interface {
	Create(host *Host) error
	List() ([]Host, error)
	Read(id uuid.UUID) (Host, error)
	ReadByHostname(hostname string) (Host, error)
	Update(id uuid.UUID) error
	Delete(id uuid.UUID) error

	//ListBeacons(host Host) ([]Host, error)
	ListHostIOC(host Host) ([]IOC, error)
	ListHostLoot(host *Host) ([]Loot, error)
}
