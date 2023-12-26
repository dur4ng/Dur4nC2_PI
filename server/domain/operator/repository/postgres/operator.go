package postgres

import (
	"Dur4nC2/server/domain/models"
	"log"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PostgresOperatorRepository struct {
	Conn *gorm.DB
}

func NewPostgresOperatorRepository(conn *gorm.DB) models.OperatorRepository {
	return &PostgresOperatorRepository{conn}
}

func (m *PostgresOperatorRepository) Create(operator *models.Operator) error {
	result := m.Conn.Create(operator)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (m *PostgresOperatorRepository) List() ([]models.Operator, error) {
	var operators []models.Operator
	result := m.Conn.Find(&operators)
	if result.Error != nil {
		log.Printf("Error: Could not list operators...")
		return nil, result.Error
	}
	return operators, nil
}
func (m *PostgresOperatorRepository) Read(id uuid.UUID) (models.Operator, error) {
	var operator models.Operator
	result := m.Conn.First(&operator, id)
	if result.Error != nil {
		return models.Operator{}, result.Error
	}
	return operator, nil
}
func (m *PostgresOperatorRepository) Update(operator *models.Operator) error {
	err := m.Conn.Model(&operator).Where("id = ?", uuid.UUID.String(operator.ID)).Updates(models.Operator{Username: operator.Username, Token: operator.Token}) // aunque le pase un objeto con datos al original no peta porque usa el id
	if err.Error != nil {
		return err.Error
	}
	return nil
}

// TODO
func (m *PostgresOperatorRepository) Delete(id uuid.UUID) error {
	return nil
}
