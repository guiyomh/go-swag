package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/guiyomh/swagger/pkg/logger"
	"github.com/guiyomh/swagger/pkg/router"
	"github.com/guiyomh/swagger/pkg/swagger"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:                 "go-swag",
		Usage:                "Generate a openapi file",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "dir",
				Aliases:  []string{"d"},
				Required: true,
				Usage:    "Directories you want to parse",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output directory for all the generated files",
				Value:   "./docs/api",
			},
		},
		Action: generateOpenAPI,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func generateOpenAPI(ctx *cli.Context) error {
	log, err := logger.New(os.Stderr, logger.InfoLevel)
	if err != nil {
		return err
	}
	defer log.Close()

	log.
		WithField("output", ctx.String("output")).
		WithField("paths", strings.Join(ctx.StringSlice("dir"), ",")).
		Info("Starting...")

	routers := make([]*router.Router, 0)

	for _, dir := range ctx.StringSlice("dir") {
		rtes, err := router.FromComments(dir)
		if err != nil {
			return err
		}
		routers = append(routers, rtes...)

	}
	swag, err := swagger.New(
		"my title",
		"my description",
		"2.0.0",
		routers,
	)
	if err != nil {
		return err
	}
	buf, err := json.MarshalIndent(swag, "", "  ")
	if err != nil {
		return err
	}
	output := ctx.String("output")
	err = ensureExist(output)
	if err != nil {
		return err
	}
	log.Info("open api generated in:" + output + "/openapi.json")

	return ioutil.WriteFile(output+"/openapi.json", buf, os.ModePerm)
}

func ensureExist(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return os.MkdirAll(path, os.ModeDir|os.ModePerm)

}
