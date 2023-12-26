package generate

import (
	"Dur4nC2/server/console"
	"Dur4nC2/server/domain/models"
	serverGenerate "Dur4nC2/server/generate"
	"errors"
	"fmt"
	"github.com/desertbit/grumble"
	"time"
)

var (
	minBeaconInterval         = 5 * time.Second
	ErrBeaconIntervalTooShort = fmt.Errorf("beacon interval must be %v or greater", minBeaconInterval)
)

func GenerateBeaconCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	err := GenerateBeacon(ctx, con)
	return err
}

func GenerateBeacon(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	config := parseCompileFlags(ctx, con)
	if config == nil {
		return errors.New("non valid beacon config provided")
	}
	config.IsBeacon = true
	err := parseBeaconFlags(ctx, con, config)
	if err != nil {
		return errors.New("non valid flags provided")
	}
	_, err = serverGenerate.CreateImplant(*config)
	if err != nil {
		return errors.New("server could not create the implant")
	}
	con.PrintSuccessf("Beacon generated: /tmp/implant.*")
	return nil
}

func parseBeaconFlags(ctx *grumble.Context, con *console.ServerConsoleClient, config *models.ImplantConfig) error {
	interval := time.Duration(ctx.Flags.Int64("days")) * time.Hour * 24
	interval += time.Duration(ctx.Flags.Int64("hours")) * time.Hour
	interval += time.Duration(ctx.Flags.Int64("minutes")) * time.Minute

	if (ctx.Flags["seconds"].IsDefault && interval.Seconds() == 0) || (!ctx.Flags["seconds"].IsDefault) {
		interval += time.Duration(ctx.Flags.Int64("seconds")) * time.Second
	}

	if interval < minBeaconInterval {
		return ErrBeaconIntervalTooShort
	}
	config.BeaconInterval = int64(interval)
	config.BeaconJitter = int64(time.Duration(ctx.Flags.Int64("jitter")) * time.Second)
	return nil
}
