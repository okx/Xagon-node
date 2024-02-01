package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/log"
	"github.com/0xPolygonHermez/zkevm-node/tools/signer/config"
	"github.com/0xPolygonHermez/zkevm-node/tools/signer/service"
	"github.com/urfave/cli/v2"
)

const (
	appName      = "x1-signer" //nolint:gosec
	appUsage     = "x1 signer tool"
	timeout      = 5 * time.Second
	httpGetPath  = "/priapi/v1/assetonchain/ecology/querySignDataByOrderNo"
	httpPostPath = "/priapi/v1/assetonchain/ecology/ecologyOperate"
)

var (
	configFileFlag = cli.StringFlag{
		Name:        config.FlagCfg,
		Aliases:     []string{"c"},
		Usage:       "Configuration `FILE`",
		DefaultText: "./config/signer.config.toml",
		Required:    true,
	}
)

// main is the entry point for the tool
func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage

	app.Commands = []*cli.Command{
		{
			Name:    "http",
			Aliases: []string{},
			Usage:   "Generate stream file from scratch",
			Action:  HttpService,
			Flags: []cli.Flag{
				&configFileFlag,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("Error: %v", err)
		os.Exit(1)
	}
}

// HttpService is the entry point for the http service
func HttpService(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx)
	if err != nil {
		log.Errorf("Error: %v", err)
		os.Exit(1)
	}

	log.Init(c.Log)
	srv := service.NewServer(c, cliCtx.Context)
	http.HandleFunc(httpGetPath, srv.GetSignDataByOrderNo)
	http.HandleFunc(httpPostPath, srv.PostSignDataByOrderNo)

	log.Infof("Listen port:%v", c.Port)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", c.Port),
		ReadHeaderTimeout: timeout,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Errorf("Error:%v", err)
	}

	return nil
}
