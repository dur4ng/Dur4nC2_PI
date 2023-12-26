package docker

import (
	"Dur4nC2/server/domain/models"
	"log"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type dockerOperatorRepository struct {
	Conn *gorm.DB
}

func NewDockerOpertatorRepository(conn *gorm.DB) models.OperatorRepository {
	return &dockerOperatorRepository{conn}
}

func (m *dockerOperatorRepository) Create(operator *models.Operator) error {
	result := m.Conn.Create(operator)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (m *dockerOperatorRepository) List() ([]models.Operator, error) {
	var operators []models.Operator
	result := m.Conn.Find(&operators)
	if result.Error != nil {
		log.Printf("Repository Error: Could not list operators...")
		return nil, result.Error
	}
	return operators, nil
}
func (m *dockerOperatorRepository) Read(id uuid.UUID) (models.Operator, error) {
	var operator models.Operator
	result := m.Conn.First(&operator, id)
	if result.Error != nil {
		return models.Operator{}, result.Error
	}
	return operator, nil
}
func (m *dockerOperatorRepository) Update(operator *models.Operator) error {
	err := m.Conn.Model(&operator).Where("id = ?", uuid.UUID.String(operator.ID)).Updates(models.Operator{Username: operator.Username, Token: operator.Token}) // aunque le pase un objeto con datos al original no peta porque usa el id
	if err.Error != nil {
		return err.Error
	}
	return nil
}

// TODO
func (m *dockerOperatorRepository) Delete(id uuid.UUID) error {
	return nil
}
