package postgres

import (
	"Dur4nC2/server/domain/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PostgresHostRepository struct {
	Conn *gorm.DB
}

func NewPostgresHostRepository(conn *gorm.DB) models.HostRepository {
	return &PostgresHostRepository{conn}
}

func (m *PostgresHostRepository) Create(host *models.Host) error {
	result := m.Conn.Create(host)
	if result != nil {
		return result.Error
	}
	return nil
}
func (m *PostgresHostRepository) List() ([]models.Host, error) {
	var hosts []models.Host
	result := m.Conn.Find(&hosts)
	if result.Error != nil {
		return nil, result.Error
	}
	return hosts, nil
}
func (m *PostgresHostRepository) Read(id uuid.UUID) (models.Host, error) {
	var h models.Host
	result := m.Conn.Where("id = ?", uuid.UUID.String(id)).First(&h)
	if result.Error != nil {
		return models.Host{}, result.Error
	}
	return h, nil
}
func (m *PostgresHostRepository) ReadByHostname(hostname string) (models.Host, error) {
	var h models.Host
	result := m.Conn.Where("hostname = ?", hostname).First(&h)
	if result.Error != nil {
		return models.Host{}, result.Error
	}
	return h, nil
}
func (m *PostgresHostRepository) Delete(id uuid.UUID) error {
	result := m.Conn.Delete(&models.Host{ID: id})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *PostgresHostRepository) AddBeacon(beacon models.Beacon, host models.Host) (string, error) {
	result := m.Conn.Model(&host).Association("Beacons").Append([]models.Beacon{beacon})
	//result := m.Conn.Model(&host).Debug().Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Table("beacons").Association("Beacons").Append([]models.Beacon{beacon})
	if result != nil {
		return "", result
	}
	//m.Conn.Save(&beacon)
	return beacon.ID.String(), nil
}

func (m *PostgresHostRepository) Update(id uuid.UUID) error { return nil }

func (m *PostgresHostRepository) ListHostIOC(host models.Host) ([]models.IOC, error) {
	var iocs []models.IOC
	err := m.Conn.Model(&host).Association("IOCs").Find(&iocs)
	if err != nil {
		return nil, err
	}
	return iocs, nil
}
func (m *PostgresHostRepository) ListHostLoot(host *models.Host) ([]models.Loot, error) {
	var loots []models.Loot
	err := m.Conn.Model(&host).Association("Loots").Find(&loots)
	if err != nil {
		return nil, err
	}
	return loots, nil
}
