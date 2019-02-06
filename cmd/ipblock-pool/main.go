package main

import (
	"fmt"
	"os"

	"github.com/kcq/poc-ipblock-pool/pkg/pool"
	"github.com/kcq/poc-ipblock-pool/internal/app/cli"
)

func main() {
	fmt.Println("IP Block Allocator PoC (cli)...")
	
	pmanager := pool.New()
	app := cli.New(pmanager)
	app.Run(os.Args)
}
