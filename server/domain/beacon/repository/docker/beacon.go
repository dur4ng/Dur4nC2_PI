package postgres

import (
	"Dur4nC2/server/domain/models"
	"log"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type mysqlBeaconRepository struct {
	Conn *gorm.DB
}

func NewBeaconRepository(conn *gorm.DB) models.BeaconRepository {
	return &mysqlBeaconRepository{conn}
}

func (m *mysqlBeaconRepository) Create(beacon models.Beacon, host models.Host) error {
	result := m.Conn.Model(host).Association("Beacons").Append([]models.Beacon{beacon})
	if result != nil {
		return result
	}
	return nil
}
func (m *mysqlBeaconRepository) List() ([]models.Beacon, error) {
	var beacons []models.Beacon
	result := m.Conn.Find(&beacons)
	if result.Error != nil {
		log.Printf("Error: Could not list opetors...")
		return nil, result.Error
	}
	return beacons, nil
}
func (m *mysqlBeaconRepository) Read(id uuid.UUID) (models.Beacon, error) {
	var b models.Beacon

	m.Conn.Where("id = ?", uuid.UUID.String(id)).First(&b)
	return b, nil
}
func (m *mysqlBeaconRepository) Update(beacon models.Beacon) error { return nil }
func (m *mysqlBeaconRepository) Delete(id uuid.UUID) error         { return nil }

func (m *mysqlBeaconRepository) ListTasks(beacon models.Beacon) ([]models.BeaconTask, error) {
	var tasks []models.BeaconTask
	result := m.Conn.Model(&beacon).Association("Tasks").Find(&tasks)
	if result != nil {
		return nil, result
	}
	return tasks, nil
}
