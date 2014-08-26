package main

import (
	"log"

	"fmt"

	"github.com/codegangsta/cli"
	"github.com/kyokomi/appConfig"
	"github.com/kyokomi/scan"
	"github.com/skratchdot/open-golang/open"
)

const accessTokenUrl = "https://kyokomi-oauth2.herokuapp.com/access"

var ac *appConfig.AppConfig

var Commands = []cli.Command{
	commandAdd,
	commandList,
	commandDeleteConfig,
}

var commandAdd = cli.Command{
	Name:  "add",
	Usage: "",
	Description: `
`,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "path", Value: "", Usage: "", EnvVar: ""},
	},
	Action: doAdd,
}

var commandList = cli.Command{
	Name:  "list",
	Usage: "",
	Description: `
`,
	Action: doList,
}

var commandDeleteConfig = cli.Command{
	Name:  "delete-config",
	Usage: "",
	Description: `
`,
	Action: doDeleteConfig,
}

func doAdd(c *cli.Context) {
	ac = appConfig.NewAppConfig(c.App.Name)
	t, err := readAccessToken()
	if err != nil {
		log.Fatal(err)
	}

	d := NewDropBox(t)
	if err := d.SetupCache(ac.ConfigDirPath); err != nil {
		log.Fatal(err)
	}

	filePath := c.String("path")
	image, err := d.AddImage(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("![%s](%s)\n", image.Name, image.URL)
}

func doList(c *cli.Context) {
	ac = appConfig.NewAppConfig(c.App.Name)
	t, err := readAccessToken()
	if err != nil {
		log.Fatal(err)
	}

	d := NewDropBox(t)
	if err := d.SetupCache(ac.ConfigDirPath); err != nil {
		log.Fatal(err)
	}
	defer d.level.Close()

	l, err := d.ReadImageList()
	if err != nil {
		log.Fatal(err)
	}

	if len(l) == 0 {
		fmt.Println("not files")
	}

	for _, s := range l {
		fmt.Printf("![%s](%s)\n", s.Name, s.URL)
	}
}

func doDeleteConfig(c *cli.Context) {
	ac = appConfig.NewAppConfig(c.App.Name)

	if err := resetAccessToken(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("delete config successful!")
	fmt.Println("path: ", ac.ConfigDirPath)
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
	accessToken := string(data)
	if err != nil || accessToken == "" {

		// OAuth jump
		open.Run(accessTokenUrl)

		// Scan accessToken
		accessToken = s.Scan("token")

		// config write
		if err := ac.WriteAppConfig([]byte(accessToken)); err != nil {
			return "", err
		}
	}

	return accessToken, nil
}

func resetAccessToken() error {
	return ac.RemoveAppConfig()
}
