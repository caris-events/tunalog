package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/caris-events/tunalog/config"
	"github.com/caris-events/tunalog/handler"
	_ "github.com/caris-events/tunalog/store"
	_ "github.com/caris-events/tunalog/translation"
)

func main() {
	fmt.Println(`
	████████╗██╗   ██╗███╗   ██╗ █████╗ ██╗      ██████╗  ██████╗
	╚══██╔══╝██║   ██║████╗  ██║██╔══██╗██║     ██╔═══██╗██╔════╝
	   ██║   ██║   ██║██╔██╗ ██║███████║██║     ██║   ██║██║  ███╗
	   ██║   ╚██████╔╝██║ ╚████║██║  ██║███████╗╚██████╔╝╚██████╔╝
	   ╚═╝    ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝ ╚═════╝  ╚═════╝ `)

	app := &cli.App{
		Name:    "tunalog",
		Version: config.Version,
		Usage:   "A simple blogging system written in Golang ✨",
		Action:  start,
		Commands: []*cli.Command{
			{
				Name:   "reset-user",
				Usage:  "reset the password or unsuspend a user, flag `-u` required",
				Action: resetUser,
			},
			{
				Name:   "reset-key",
				Usage:  "reset the auth key",
				Action: resetKey,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//=======================================================
// Start
//=======================================================

func start(c *cli.Context) error {
	return handler.Instance.Engine.Run(":8080")
}

//=======================================================
// Reset Key
//=======================================================

func resetKey(c *cli.Context) error {
	return nil
}

//=======================================================
// Reset User
//=======================================================

func resetUser(c *cli.Context) error {
	return nil
}
