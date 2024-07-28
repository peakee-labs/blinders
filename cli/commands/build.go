package commands

import (
	"fmt"

	"blinders/packages/utils"

	"github.com/fatih/color"
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

type ChangeResult struct {
	Change string
	Path   string
}

func BuildAction(ctx *cli.Context) error {
	var names []string

	if ctx.Bool("all") {
		color.Yellow("\nUse --all, finding all lambdas, ignoring --name")
		result, err := Execute("sh", "scripts/lookup_lambdas.sh")
		if err != nil {
			return err
		}

		lambdaMap, err := utils.ParseJSON[map[string]string]([]byte(result))
		if err != nil {
			color.Red("Can not parse ChangeResult: %v", err)
		}

		for key := range *lambdaMap {
			names = append(names, key)
		}
	} else {
		names = ctx.StringSlice("name")
	}

	for _, name := range names {
		color.Magenta("\n[%v] checking and building", name)

		servicePath := fmt.Sprintf(ServicesPathPrefix, name)

		result, err := Execute("sh", "scripts/check_go_mod_changes.sh", servicePath)
		if err != nil {
			return err
		}

		changeResult, err := utils.ParseJSON[ChangeResult]([]byte(result))
		if err != nil {
			color.Red("Can not parse ChangeResult: %v", err)
		}

		switch changeResult.Change {
		case "unchange":
			if ctx.Bool("force") {
				color.Magenta("[%v] unchange, force building", name)
			} else {
				color.Yellow("[%v] unchange, ignore this build", name)
			}
		case "change":
			color.Magenta("[%v] Change found, building", name)
		default:
			return fmt.Errorf("[%v] Unknown change result: %v", name, changeResult.Change)
		}
	}

	return nil
}
