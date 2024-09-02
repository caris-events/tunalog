package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thanhpk/randstr"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/caris-events/tunalog/handler"
	"github.com/caris-events/tunalog/store"
	"github.com/caris-events/tunalog/system"
)

func main() {
	fmt.Println(`
████████╗██╗   ██╗███╗   ██╗ █████╗ ██╗      ██████╗  ██████╗
╚══██╔══╝██║   ██║████╗  ██║██╔══██╗██║     ██╔═══██╗██╔════╝
   ██║   ██║   ██║██╔██╗ ██║███████║██║     ██║   ██║██║  ███╗
   ██║   ╚██████╔╝██║ ╚████║██║  ██║███████╗╚██████╔╝╚██████╔╝
   ╚═╝    ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝ ╚═════╝  ╚═════╝`)

	app := &cli.App{
		Name:    "tunalog",
		Version: system.Version,
		Usage:   "A simple blogging system written in Golang ✨",
		Action:  start,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Usage:   "port to listen on",
				Aliases: []string{"p"},
				Value:   ":8080",
			},
			&cli.StringFlag{
				Name:  "tls-key",
				Usage: "path to TLS key file",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "tls-crt",
				Usage: "path to TLS certificate file",
				Value: "",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "reset-password",
				Usage:  "reset the password of a user, email address is required ",
				Action: resetUser,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) error {
	if c.String("tls-crt") != "" && c.String("tls-key") != "" {
		fmt.Printf("👋 Visit https://localhost%s to use Tunalog\n", c.String("port"))
		return handler.Router.RunTLS(c.String("port"), c.String("tls-crt"), c.String("tls-key"))
	} else {
		fmt.Printf("👋 Visit http://localhost%s to use Tunalog\n", c.String("port"))
		return handler.Router.Run(c.String("port"))
	}
}

func resetUser(c *cli.Context) error {
	u, err := store.GetUserByEmail(c.Args().First())
	if err != nil {
		return err
	}
	pwd := randstr.String(16, randstr.Base64Chars)
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := store.UpdateUserPassword(u.ID, string(b)); err != nil {
		return err
	}
	log.Printf(`Password for user %s has been reset to: "%s"`, u.Email, pwd)
	return nil
}
