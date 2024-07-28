package commands

import (
	"fmt"
	"log"
	"os"

	"blinders/packages/auth"
	authutils "blinders/packages/auth/utils"
	"blinders/packages/utils"

	firebaseAuth "firebase.google.com/go/auth"
	"github.com/urfave/cli/v2"
)

var client *firebaseAuth.Client

var AuthCommand = cli.Command{
	Name:        "auth",
	Subcommands: []*cli.Command{&loadAuthCommand, &genWSCatCommand},
	Before: func(ctx *cli.Context) error {
		env := ctx.String("env")
		adminJSON, err := utils.GetFile(fmt.Sprintf("firebase.admin.%v.json", env))
		if err != nil {
			return err
		}

		a, err := auth.NewFirebaseManager(adminJSON, nil)
		if err != nil {
			return err
		}
		client = a.Client

		return nil
	},
}

var loadAuthCommand = cli.Command{
	Name:  "jwt",
	Usage: "load user jwt by using uid",
	Args:  true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "uid",
		},
	},
	Action: func(ctx *cli.Context) error {
		webAPIKey := os.Getenv("WEB_API_KEY")
		if webAPIKey == "" {
			log.Fatal("WEB_API_KEY is required from environment")
		}
		env := ctx.String("env")
		var uid string
		if uid = ctx.String("uid"); uid == "" {
			uid = os.Getenv("USER_UID")
		}
		if uid == "" {
			log.Fatal("USER_UID is required from environment")
		}

		cacheFile := fmt.Sprintf("auth.%v.json", env)
		idToken, authToken, err := authutils.LoadFirebaseAuthForUserWithCache(
			client,
			uid,
			webAPIKey,
			cacheFile,
		)
		if err != nil {
			return fmt.Errorf("failed to load firebase auth: %v", err)
		}

		fmt.Printf("JWT of %v: %v\n", authToken.Firebase.Identities["email"], idToken)

		return nil
	},
}

var genWSCatCommand = cli.Command{
	Name:  "wscat",
	Usage: "load user jwt and generate wscat command for easily connect to websocket api by using uid",
	Args:  true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "uid",
		},
		&cli.StringFlag{
			Name:     "endpoint",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		webAPIKey := os.Getenv("WEB_API_KEY")
		if webAPIKey == "" {
			log.Fatal("WEB_API_KEY is required from environment")
		}
		env := ctx.String("env")
		var uid string
		if uid = ctx.String("uid"); uid == "" {
			uid = os.Getenv("USER_UID")
		}
		if uid == "" {
			log.Fatal("USER_UID is required from environment")
		}

		endpoint := ctx.String("endpoint")

		cacheFile := fmt.Sprintf("auth.%v.json", env)
		idToken, _, err := authutils.LoadFirebaseAuthForUserWithCache(
			client,
			uid,
			webAPIKey,
			cacheFile,
		)
		if err != nil {
			return fmt.Errorf("failed to load firebase auth: %v", err)
		}

		fmt.Printf("wscat -c \"wss://%v?token=%v\"\n", endpoint, idToken)

		return nil
	},
}
