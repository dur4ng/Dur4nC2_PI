package beacons

import (
	"Dur4nC2/server/console"
	"github.com/desertbit/grumble"
)

func BeaconsBackgroundCmd(ctx *grumble.Context, con *console.ServerConsoleClient) {
	con.ActiveTarget.Background()
}
