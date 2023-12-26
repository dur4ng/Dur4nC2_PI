package postgres

import (
	cryptography "Dur4nC2/misc/crypto"
	"Dur4nC2/server/domain/models"
	"crypto/sha256"
	"encoding/base64"
	"gorm.io/gorm"
)

type postgresImplantRepository struct {
	Conn *gorm.DB
}

func NewPosgresImplantRepository(conn *gorm.DB) models.ImplantConfigRepository {
	return &postgresImplantRepository{conn}
}

func (p *postgresImplantRepository) Create() (*models.ImplantConfig, error) {
	keyPair, err := cryptography.RandomECCKeyPair()
	if err != nil {
		return &models.ImplantConfig{}, err
	}
	digest := sha256.Sum256(keyPair.Public[:])
	implantConfig := &models.ImplantConfig{
		ECCPublicKey:       keyPair.PublicBase64(),
		ECCPrivateKey:      keyPair.PrivateBase64(),
		ECCPublicKeyDigest: base64.RawStdEncoding.EncodeToString(digest[:]),
	}
	result := p.Conn.Create(implantConfig)

	if result.Error != nil {
		return &models.ImplantConfig{}, result.Error
	}
	//fmt.Println(implantConfig.ECCPublicKeyDigest)
	//fmt.Println(base64.RawStdEncoding.EncodeToString(digest[:]))
	//newImplantConfig, err := p.ReadByECCPlublicKeyDigest(&digest)
	return implantConfig, nil
}

func (p *postgresImplantRepository) ReadByECCPlublicKeyDigest(eccPublicKeyDigest *[32]byte) (models.ImplantConfig, error) {
	var ic models.ImplantConfig

	result := p.Conn.Where("ecc_public_key_digest = ?", base64.RawStdEncoding.EncodeToString(eccPublicKeyDigest[:])).First(&ic)
	if result.Error != nil {
		return models.ImplantConfig{}, gorm.ErrRecordNotFound
	}
	return ic, nil
}
