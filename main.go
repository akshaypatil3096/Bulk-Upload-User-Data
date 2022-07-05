package main

import (
	"os"

	"github.com/akshaypatil3096/Bulk-Upload-User-Data/service"
	"github.com/urfave/cli"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "Golang App"
	cliApp.Version = "1.0.0"
	cliApp.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Starts the insert process for the provided number of records",
			Action: func(c *cli.Context) error {
				err := service.Start(c.Args().Get(0))
				return err
			},
		},
		{
			Name:  "resume",
			Usage: "resume the insert process",
			Action: func(c *cli.Context) error {
				err := service.ResumeInsertData()
				return err
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}

}
