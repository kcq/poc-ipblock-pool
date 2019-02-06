package cli

import (
	"bytes"
	"encoding/json"
	"fmt"

	ucli "github.com/urfave/cli"

	"github.com/kcq/poc-ipblock-pool/pkg/pool"
)

const (
	flagBlock = "block"
	flagKey   = "key"
)

type App struct {
	pm  *pool.Manager
	cli *ucli.App
}

func New(pmanager *pool.Manager) *App {
	app := &App{
		pm:  pmanager,
		cli: ucli.NewApp(),
	}

	app.init()
	return app
}

func (a *App) init() {
	a.cli.Name = "ipblock-pool"
	a.cli.Usage = "IP Block allocator PoC"

	blockKeyFlag := ucli.StringFlag{
		Name:  flagKey,
		Value: "",
		Usage: "Block key",
	}

	blockIPFlag := ucli.StringFlag{
		Name:  flagBlock,
		Value: "",
		Usage: "Starting IP address of the IP block",
	}

	a.cli.Commands = []ucli.Command{
		{
			Name:    "lookup",
			Aliases: []string{"g"},
			Usage:   "lookup IP block allocation by IP block or key",
			Flags: []ucli.Flag{
				blockKeyFlag,
				blockIPFlag,
			},
			Action: func(ctx *ucli.Context) error {
				key := ctx.String(flagKey)
				block := ctx.String(flagBlock)

				blockInfo := a.pm.Lookup(block, key)

				if blockInfo == nil {
					fmt.Println("Block not found")
				} else {
					printBlockInfo(blockInfo)
				}

				return nil
			},
		},
		{
			Name:    "allocate",
			Aliases: []string{"a"},
			Usage:   "allocate a new IP block",
			Flags: []ucli.Flag{
				blockKeyFlag,
			},
			Action: func(ctx *ucli.Context) error {
				key := ctx.String(flagKey)

				blockInfo := a.pm.Allocate(key, false)
				printBlockInfo(blockInfo)

				return nil
			},
		},
		{
			Name:    "free",
			Aliases: []string{"d"},
			Usage:   "free an IP block",
			Flags: []ucli.Flag{
				blockKeyFlag,
				blockIPFlag,
			},
			Action: func(ctx *ucli.Context) error {
				key := ctx.String(flagKey)
				block := ctx.String(flagBlock)

				err := a.pm.Free(block, key)

				switch err {
				case pool.ErrBlockNotFound:
					fmt.Println("Block not found!")
				case nil:
					fmt.Println("Done!")
				default:
					fmt.Println(err)
				}
				return nil
			},
		},
	}
}

func (a *App) Run(args []string) {
	a.cli.Run(args)
}

func printBlockInfo(value interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(value); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(buf.String())
}
