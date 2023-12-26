package postgres

import (
	"Dur4nC2/server/domain/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PostgresLootRepository struct {
	Conn *gorm.DB
}

func NewPostgresLootRepository(conn *gorm.DB) models.LootRepository {
	return &PostgresLootRepository{conn}
}

func (m *PostgresLootRepository) Create(loot models.Loot, host models.Host) error {
	result := m.Conn.Model(&host).Association("Loots").Append([]models.Loot{loot})
	if result != nil {
		return result
	}
	return nil
}
func (m *PostgresLootRepository) List() ([]models.Loot, error) {
	var loots []models.Loot
	result := m.Conn.Find(&loots)
	if result.Error != nil {
		return nil, result.Error
	}
	return loots, nil
}
func (m *PostgresLootRepository) Read(id uuid.UUID) (models.Loot, error) {
	var l models.Loot
	result := m.Conn.Where("id = ?", uuid.UUID.String(id)).First(&l)
	if result.Error != nil {
		return models.Loot{}, result.Error
	}
	return l, nil
}
