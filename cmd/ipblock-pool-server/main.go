package main

import (
	"fmt"
	"os"

	"github.com/kcq/poc-ipblock-pool/internal/app/server"
	"github.com/kcq/poc-ipblock-pool/pkg/pool"
)

func main() {
	fmt.Println("IP Block Allocator PoC...")

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
	app := server.New(pmanager)
	app.Run()
}
