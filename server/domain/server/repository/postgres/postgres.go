package postgres

import (
	"Dur4nC2/server/domain/models"
	"gorm.io/gorm"
)

type postgresServerRepository struct {
	Conn *gorm.DB
}

func NewPostgresServerRepository(conn *gorm.DB) models.ServerRepository {
	return &postgresServerRepository{conn}
}

func (p *postgresServerRepository) SetKeyValue(key string, value string) error {
	err := p.Conn.Where(&models.Server{
		Key: key,
	}).First(&models.Server{}).Error
	if err == gorm.ErrRecordNotFound {
		err = p.Conn.Create(&models.Server{
			Key:   key,
			Value: value,
		}).Error
	} else {
		err = p.Conn.Where(&models.Server{
			Key: key,
		}).Updates(models.Server{
			Key:   key,
			Value: value,
		}).Error
	}
	return err
}

func (p *postgresServerRepository) GetKeyValue(key string) (string, error) {
	keyValue := &models.Server{}
	err := p.Conn.Where(&models.Server{
		Key: key,
	}).First(keyValue).Error
	return keyValue.Value, err
}

func (p *postgresServerRepository) DeleteKeyValue(key string) error {
	return nil
}
