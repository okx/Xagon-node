package main

import (
	"os"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/urfave/cli/v2"
)

const (
	flagInput = "input"
)

func main() {
	app := cli.NewApp()
	app.Name = "x1-node-scripts"
	app.Commands = []*cli.Command{
		{
			Name:   "updatedeps",
			Usage:  "Updates external dependencies like images, test vectors or proto files",
			Action: updateDeps,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "compilesc",
			Usage:  "Compiles smart contracts required for testing",
			Action: compileSC,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     flagInput,
					Aliases:  []string{"in"},
					Usage:    "Target path where the source solidity files are located. It can be a file or a directory, in which case the command will traverse all the descendants from it.",
					Required: true,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
