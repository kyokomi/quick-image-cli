package main

import (
	"log"

	"fmt"

	"github.com/codegangsta/cli"
	"github.com/kyokomi/appConfig"
	"github.com/kyokomi/scan"
	"github.com/skratchdot/open-golang/open"
)

const accessTokenURL = "https://kyokomi-oauth2.herokuapp.com/access"

type DropBoxAppConfig struct {
	appConfig.AppConfig
}

func (d *DropBoxAppConfig) readAccessToken() (string, error) {
	data, err := d.ReadAppConfig()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readDropBoxAppConfig(appName string) (*DropBoxAppConfig, error) {
	s := scan.CliScan{
		Scans: []scan.Scan{
			{Name: "token",
				Value: "",
				Usage: "please your dropbox accessToken",
				Env:   "",
			},
		},
	}

	ac := appConfig.NewDefaultAppConfig(appName)
	data, err := ac.ReadAppConfig()
	accessToken := string(data)
	if err != nil || accessToken == "" {

		// OAuth jump
		open.Run(accessTokenURL)

		// Scan accessToken
		accessToken = s.Scan("token")

		// config write
		if err := ac.WriteAppConfig([]byte(accessToken)); err != nil {
			return nil, err
		}
	}

	return &DropBoxAppConfig{*ac}, nil
}

func newDropBox(appName string) (*DropBox, error) {
	ac, err := readDropBoxAppConfig(appName)
	if err != nil {
		return nil, err
	}
	t, err := ac.readAccessToken()
	if err != nil {
		return nil, err
	}

	return NewDropBox(t), nil
}

var Commands = []cli.Command{
	commandAdd,
	commandList,
	commandDeleteConfig,
	commandCreateFolder,
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

var commandCreateFolder = cli.Command{
	Name:      "create-folder",
	ShortName: "C",
	Usage:     "",
	Description: `
`,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "path", Value: "", Usage: "", EnvVar: ""},
	},
	Action: doCreateFolder,
}

var commandList = cli.Command{
	Name:  "list",
	Usage: "",
	Description: `
`,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "dir", Usage: "", EnvVar: ""},
	},
	Action: doList,
}

var commandDeleteConfig = cli.Command{
	Name:      "delete-config",
	ShortName: "D",
	Usage:     "",
	Description: `
`,
	Action: doDeleteConfig,
}

func doAdd(c *cli.Context) {
	d, err := newDropBox(c.App.Name)
	if err != nil {
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

	d, err := newDropBox(c.App.Name)
	if err != nil {
		log.Fatal(err)
	}

	l, err := d.ReadImageList(c.Bool("dir"))
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
	ac := appConfig.NewDefaultAppConfig(c.App.Name)
	if err := ac.RemoveAppConfig(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("delete config successful!")
	fmt.Println("path: ", ac.ConfigDirPath)
}

func doCreateFolder(c *cli.Context) {
	d, err := newDropBox(c.App.Name)
	if err != nil {
		log.Fatal(err)
	}

	res, err := d.CreateFolder(c.String("path"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res))
}
