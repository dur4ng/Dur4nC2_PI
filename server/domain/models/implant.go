package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

type OutputFormat int32

const (
	OutputFormat_SHARED_LIB OutputFormat = 0
	OutputFormat_SHELLCODE  OutputFormat = 1
	OutputFormat_EXECUTABLE OutputFormat = 2
	OutputFormat_SERVICE    OutputFormat = 3
)

// ImplantConfig - An implant build configuration
type ImplantConfig struct {
	ID uuid.UUID `gorm:"primaryKey"`
	//Beacon    []Beacon
	CreatedAt time.Time

	Name string

	ConnectionMethod string

	Domain     string
	URL        string
	PathPrefix string

	IsBeacon       bool
	BeaconInterval int64
	BeaconJitter   int64

	ECCPublicKey          string
	ECCPublicKeyDigest    string
	ECCPrivateKey         string
	ECCPublicKeySignature string
	ECCServerPublicKey    string

	OS          string
	Format      OutputFormat
	IsSharedLib bool
	IsService   bool
	IsShellcode bool

	ImplantPackagePath string
}

// BeforeCreate - GORM hook
func (ic *ImplantConfig) BeforeCreate(tx *gorm.DB) (err error) {
	ic.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	ic.CreatedAt = time.Now()
	return nil
}

type ImplantConfigUsecase interface {
	Create(implantConfig ImplantConfig) error
	ReadByECCPlublicKeyDigest(eccPublicKeyDigest string) (ImplantConfig, error)
}
type ImplantConfigRepository interface {
	Create() (*ImplantConfig, error)
	ReadByECCPlublicKeyDigest(eccPublicKeyDigest *[32]byte) (ImplantConfig, error)
}
