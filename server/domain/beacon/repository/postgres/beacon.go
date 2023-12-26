package postgres

import (
	"Dur4nC2/server/db"
	"Dur4nC2/server/domain/models"
	taskRepository "Dur4nC2/server/domain/task/repository/postgres"
	"log"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PostgresBeaconRepository struct {
	Conn *gorm.DB
}

func NewPostgresBeaconRepository(conn *gorm.DB) models.BeaconRepository {
	return &PostgresBeaconRepository{conn}
}

func (m *PostgresBeaconRepository) Create(beacon models.Beacon, host models.Host) (string, error) {
	result := m.Conn.Model(&host).Association("Beacons").Append([]models.Beacon{beacon})
	//result := m.Conn.Model(&host).Debug().Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Table("beacons").Association("Beacons").Append([]models.Beacon{beacon})
	if result != nil {
		return "", result
	}
	//m.Conn.Save(&beacon)
	return beacon.ID.String(), nil
}
func (m *PostgresBeaconRepository) List() ([]models.Beacon, error) {
	var beacons []models.Beacon
	//result := m.Conn.Find(&beacons)
	result := m.Conn.Order("created_at").Find(&beacons)
	if result.Error != nil {
		log.Printf("Error: Could not list opetors...")
		return nil, result.Error
	}
	return beacons, nil
}
func (m *PostgresBeaconRepository) Read(id uuid.UUID) (models.Beacon, error) {
	var b models.Beacon
	if id == uuid.Nil {
		return models.Beacon{}, gorm.ErrRecordNotFound
	}
	result := m.Conn.Where("id = ?", uuid.UUID.String(id)).First(&b)
	if result.Error != nil {
		return models.Beacon{}, result.Error
	}
	return b, nil
}
func (m *PostgresBeaconRepository) Update(beacon models.Beacon) error { return nil }
func (m *PostgresBeaconRepository) Delete(id uuid.UUID) error {
	_, err := m.Read(id)
	if err != nil {
		return err
	}
	result := m.Conn.Delete(&models.Beacon{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (m *PostgresBeaconRepository) DeleteAll() error {
	beacons, err := m.List()
	if err != nil {
		return err
	}
	taskRepo := taskRepository.NewPostgresBeaconTaskRepository(db.Session())
	for _, beacon := range beacons {
		taskRepo.CleanBeaconTasks(beacon.ID)
	}
	m.Conn.Delete(beacons)
	return nil
}
func (m *PostgresBeaconRepository) ListTasks(beacon models.Beacon) ([]models.BeaconTask, error) {
	var tasks []models.BeaconTask
	result := m.Conn.Where("state = ?", "pending").Model(&beacon).Association("Tasks").Find(&tasks)
	if result != nil {
		return nil, result
	}
	return tasks, nil
}
