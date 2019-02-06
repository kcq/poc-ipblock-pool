package main

import (
	"fmt"
	"os"

	"github.com/kcq/poc-ipblock-pool/internal/app/cli"
	"github.com/kcq/poc-ipblock-pool/pkg/pool"
)

func main() {
	fmt.Println("IP Block Allocator PoC (cli)...")

	pmanager := pool.New()
	app := cli.New(pmanager)
	app.Run(os.Args)
}
