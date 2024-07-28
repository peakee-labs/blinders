package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var BuildCommand = cli.Command{
	Name:  "build",
	Usage: "Build lambdas/services",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "all",
			Value: false,
			Usage: "Lookup all lambdas/services, ignore --name",
		},
		&cli.BoolFlag{
			Name:  "force",
			Value: false,
			Usage: "Force rebuild even if lambdas/services have no change",
		},
		&cli.SliceFlag[*cli.StringSliceFlag, []string, string]{
			Target: &cli.StringSliceFlag{
				Name:  "name",
				Usage: "specify which lambda/service to build",
			},
		},
	},
	Action: BuildAction,
}

const ServicesPathPrefix = "services/%s"

func BuildAction(ctx *cli.Context) error {
	names := ctx.StringSlice("name")

	for _, name := range names {
		servicePath := fmt.Sprintf(ServicesPathPrefix, name)
		_, err := Execute("sh", "scripts/check_go_mod_changes.sh", servicePath)
		if err != nil {
			return err
		}
	}

	return nil
}
