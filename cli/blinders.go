package main

import (
	"fmt"
	"log"
	"os"

	"blinders/cli/commands"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Blinders",
		Usage: "CLI tools for backend development",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "default",
				Usage: "Define environment for the CLI",
			},
		},
		Commands: []*cli.Command{&commands.AuthCommand},
		Before: func(ctx *cli.Context) error {
			env := ctx.String("env")
			fmt.Println("CLI is running on environment:", env)

			envFile := ".env"
			if env != "default" {
				envFile = fmt.Sprintf(".env.%s", env)
			}

			if godotenv.Load(envFile) != nil {
				log.Fatal("Can not load env:", envFile)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
