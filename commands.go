package main

import (
	"log"
	"os"

	"fmt"

	"github.com/codegangsta/cli"
	"github.com/kyokomi/appConfig"
	"github.com/kyokomi/scan"
	"github.com/skratchdot/open-golang/open"
)


var ac *appConfig.AppConfig

var Commands = []cli.Command{
	commandAdd,
	commandList,
	commandDelete,
}

var commandAdd = cli.Command{
	Name:  "add",
	Usage: "",
	Description: `
`,
	Action: doAdd,
}

var commandList = cli.Command{
	Name:  "list",
	Usage: "",
	Description: `
`,
	Action: doList,
}

var commandDelete = cli.Command{
	Name:  "delete",
	Usage: "",
	Description: `
`,
	Action: doDelete,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doAdd(c *cli.Context) {
}

func doList(c *cli.Context) {
	ac = appConfig.NewAppConfig(c.App.Name)

	t, err  := readAccessToken()
	if err != nil {
		log.Fatal(err)
	}

	d := NewDropBox(t)
	if err := d.SetupCache(ac.ConfigDirPath); err != nil {
		log.Fatal(err)
	}

	l, err := d.ReadImageList()
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range l {
		fmt.Printf("![%s](%s)\n", s.Name, s.URL)
	}
}

func doDelete(c *cli.Context) {
}

func readAccessToken() (string, error) {
	s := scan.CliScan{
		Scans: []scan.Scan{
			{Name: "token",
				Value: "",
				Usage: "please your dropbox accessToken",
			},
		},
	}

	data, err := ac.ReadAppConfig()
	if err != nil {

		// TODO: OAuth jump
		open.Run("https://localhost:8443/access")

		// Scan accessToken
		t := s.Scan("token")

		// config write
		if err := ac.WriteAppConfig([]byte(t)); err != nil {
			return "", err
		}
	}

	return string(data), nil
}
