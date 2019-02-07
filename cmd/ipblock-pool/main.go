package main

import (
	"fmt"
	"os"

	"github.com/kcq/poc-ipblock-pool/internal/app/cli"
	"github.com/kcq/poc-ipblock-pool/pkg/pool"
)

func main() {
	fmt.Println("IP Block Allocator PoC (cli)...")

	config := pool.Config{
		StartRange:    "169.254.51.0",
		EndRange:      "169.254.255.244",
		PoolBlockSize: 4,
		Store: &pool.StoreConfig{
			Address: "127.0.0.1:8500",
		},
	}

	if consulAddr, ok := os.LookupEnv("CONSUL_ADDR"); ok {
		config.Store.Address = consulAddr
		fmt.Println("Using Consul address from environment =", consulAddr)
	}

	pmanager := pool.New(&config, nil)
	app := cli.New(pmanager)
	app.Run(os.Args)
}
