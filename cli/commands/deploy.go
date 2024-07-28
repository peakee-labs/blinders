package commands

import "github.com/urfave/cli/v2"

var DeployCommand = cli.Command{
	Name:  "deploy",
	Usage: "Deploy with Terraform",
}
