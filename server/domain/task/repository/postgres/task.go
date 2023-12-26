package postgres

import (
	"Dur4nC2/server/domain/models"
	"log"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PostgresBeaconTaskRepository struct {
	Conn *gorm.DB
}

func NewPostgresBeaconTaskRepository(conn *gorm.DB) models.BeaconTaskRepository {
	return &PostgresBeaconTaskRepository{conn}
}

func (m *PostgresBeaconTaskRepository) Create(beacon models.Beacon, task models.BeaconTask) error {
	result := m.Conn.Model(&beacon).Association("Tasks").Append([]models.BeaconTask{task})
	if result != nil {
		return result
	}
	return nil
}
func (m *PostgresBeaconTaskRepository) List() ([]models.BeaconTask, error) {
	var tasks []models.BeaconTask
	result := m.Conn.Find(&tasks)
	if result.Error != nil {
		log.Printf("Error: Could not list opetors...")
		return nil, result.Error
	}
	return tasks, nil
}
func (m *PostgresBeaconTaskRepository) ListPendingTasks(id uuid.UUID) ([]models.BeaconTask, error) {
	var tasks []models.BeaconTask
	result := m.Conn.Where("beacon_id = ? AND state = 'pending'", id.String()).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}
func (m *PostgresBeaconTaskRepository) ListTasksStateAndBeaconID(id string, state string) ([]models.BeaconTask, error) {
	var tasks []models.BeaconTask
	result := m.Conn.Where(`beacon_id = ? AND state = '`+state+`'`, id).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}
func (m *PostgresBeaconTaskRepository) Read(id uuid.UUID) (models.BeaconTask, error) {
	var task models.BeaconTask

	result := m.Conn.Where("id = ?", id.String()).First(&task)
	if result.Error != nil {
		return models.BeaconTask{}, result.Error
	}
	return task, nil
}
func (m *PostgresBeaconTaskRepository) ListTaskByEnvelopeID(envelopeID int64) (models.BeaconTask, error) {
	var task models.BeaconTask

	result := m.Conn.Where("envelope_id = ?", envelopeID).First(&task)
	if result.Error != nil {
		return models.BeaconTask{}, result.Error
	}
	return task, nil
}
func (m *PostgresBeaconTaskRepository) Cancel(taskId uuid.UUID) error {
	task, err := m.Read(taskId)
	if err != nil {
		return err
	}
	result := m.Conn.Model(&task).Update("state", models.CANCELED)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (m *PostgresBeaconTaskRepository) CleanBeaconTasks(beaconID uuid.UUID) error {
	var tasks []models.BeaconTask
	result := m.Conn.Where(`beacon_id = ?`, beaconID.String()).Find(&tasks)
	if result.Error != nil {
		return result.Error
	}
	m.Conn.Delete(tasks)

	return nil

}

func (m *PostgresBeaconTaskRepository) Update() {}
func (m *PostgresBeaconTaskRepository) Delete() {}
