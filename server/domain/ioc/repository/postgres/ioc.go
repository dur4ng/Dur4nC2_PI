package postgres

import (
	"Dur4nC2/server/domain/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PostgresIOCRepository struct {
	Conn *gorm.DB
}

func NewPostgresIOCRepository(conn *gorm.DB) models.IOCRepository {
	return &PostgresIOCRepository{conn}
}

func (m *PostgresIOCRepository) Create(ioc models.IOC) error {
	result := m.Conn.Model(&models.Host{
		ID: ioc.HostID,
	}).Association("IOCs").Append([]models.IOC{ioc})
	if result != nil {
		return result
	}
	return nil
}

func (m *PostgresIOCRepository) List() ([]models.IOC, error) {
	var iocs []models.IOC
	result := m.Conn.Find(&iocs)
	if result.Error != nil {
		return nil, result.Error
	}
	return iocs, nil
}

func (m *PostgresIOCRepository) Read(id uuid.UUID) (models.IOC, error) {
	var i models.IOC
	result := m.Conn.Where("id = ?", uuid.UUID.String(id)).First(&i)
	if result.Error != nil {
		return models.IOC{}, result.Error
	}
	return i, nil
}

func (m *PostgresIOCRepository) Update(ioc models.IOC) error {
	err := m.Conn.Model(&ioc).Where("id = ?", uuid.UUID.String(ioc.ID)).Updates(models.IOC{Path: ioc.Path, FileHash: ioc.FileHash, Name: ioc.Name, Description: ioc.Description, Output: ioc.Output, State: ioc.State}) // aunque le pase un objeto con datos al original no peta porque usa el id
	if err.Error != nil {
		return err.Error
	}
	return nil
}
