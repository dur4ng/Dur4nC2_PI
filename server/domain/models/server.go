package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

type Server struct {
	ID        uuid.UUID `gorm:"primaryKey`
	CreatedAt time.Time

	Key   string `gorm:"unique;"`
	Value string
}

// BeforeCreate - GORM hook
func (s *Server) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	s.CreatedAt = time.Now()
	return nil
}

type ServerRepository interface {
	SetKeyValue(key string, value string) error
	GetKeyValue(key string) (string, error)
	DeleteKeyValue(key string) error
}
