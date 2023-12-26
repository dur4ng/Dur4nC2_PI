package models

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Beacon struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	Deleted gorm.DeletedAt
	Name    string
	HostID  uuid.UUID `gorm:"size:256"`

	Username          string
	UID               string //user id
	GID               string //group id
	OS                string
	Arch              string
	Transport         string
	RemoteAddress     string
	PID               int32 //process id
	LastCheckin       time.Time
	ReconnectInterval int64
	ActiveC2          string
	Locale            string
	CreatedAt         time.Time

	Interval    int64
	Jitter      int64
	NextCheckin int64

	Tasks []BeaconTask
}

func (b *Beacon) BeforeCreate(tx *gorm.DB) (err error) {

	b.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}

	return nil
}

type BeaconUsecase interface {
	List() (*clientpb.Beacons, error)
	Read(beacon_pb *clientpb.Beacon) (*clientpb.Beacon, error)
	Delete(beacon_pb *clientpb.Beacon) (*commonpb.Empty, error)

	ListTasks(beacon_pb *clientpb.Beacon) (*clientpb.BeaconTasks, error)
}

type BeaconRepository interface {
	Create(beacon Beacon, host Host) (string, error)
	List() ([]Beacon, error)
	Read(id uuid.UUID) (Beacon, error)
	Update(beacon Beacon) error
	Delete(id uuid.UUID) error
	DeleteAll() error
	ListTasks(beacon Beacon) ([]BeaconTask, error)
}
