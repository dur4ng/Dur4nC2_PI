package beacons

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/server/console"
	"Dur4nC2/server/db"
	"Dur4nC2/server/domain/beacon/repository/postgres"
	"Dur4nC2/server/domain/beacon/usecase"
	"github.com/desertbit/grumble"
)

func BeaconsUseCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	beaconIDFlag := ctx.Args.String("id")
	if beaconIDFlag != "" {
		beaconRepo := postgres.NewPostgresBeaconRepository(db.Session())
		beaconUsecase := usecase.NewBeaconUsecase(beaconRepo)
		beaconpb := &clientpb.Beacon{
			ID: beaconIDFlag,
		}
		beacon, err := beaconUsecase.Read(beaconpb)
		if err != nil {
			return err
		}
		con.ActiveTarget.Set(beacon)
		con.PrintInfof("Active beacon %s (%s)\n", beacon.Name, beacon.ID)
	}
	return nil
}
