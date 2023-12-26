package models

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"gorm.io/gorm"
	"time"

	"github.com/gofrs/uuid"
)

const (
	PENDING   = "pending"
	SENT      = "sent"
	COMPLETED = "completed"
	CANCELED  = "canceled"
)

type BeaconTask struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	EnvelopeID  int64
	BeaconID    uuid.UUID `gorm:"size:256"`
	CreatedAt   time.Time
	State       string
	SentAt      time.Time
	CompletedAt time.Time
	Description string
	Request     []byte // *sliverpb.Envelope
	Response    []byte // *sliverpb.Envelope
}

func (t *BeaconTask) BeforeCreate(tx *gorm.DB) (err error) {

	t.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}

	return nil
}

type BeaconTaskUsecase interface {
	List() (*clientpb.BeaconTasks, error)
	Read(task_pb *clientpb.BeaconTask) (*clientpb.BeaconTask, error)
	ListTasksStateAndBeaconID(id string, state string) (*clientpb.BeaconTasks, error)
	Cancel(task_pb *clientpb.BeaconTask) (*clientpb.BeaconTask, error)
}

type BeaconTaskRepository interface {
	Create(beacon Beacon, task BeaconTask) error
	List() ([]BeaconTask, error)
	ListPendingTasks(id uuid.UUID) ([]BeaconTask, error)
	ListTaskByEnvelopeID(envelopeID int64) (BeaconTask, error)
	ListTasksStateAndBeaconID(id string, state string) ([]BeaconTask, error)
	Read(id uuid.UUID) (BeaconTask, error)
	Cancel(taskId uuid.UUID) error
	CleanBeaconTasks(beaconID uuid.UUID) error
	Update()
	Delete()
}
