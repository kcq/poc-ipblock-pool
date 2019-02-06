package main

import (
	"fmt"

	"github.com/kcq/poc-ipblock-pool/internal/app/server"
	"github.com/kcq/poc-ipblock-pool/pkg/pool"
)

func main() {
	fmt.Println("IP Block Allocator PoC...")

	pmanager := pool.New()
	app := server.New(pmanager)
	app.Run()
}
