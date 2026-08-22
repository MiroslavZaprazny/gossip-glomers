package main

import (
	"log"

	handlers "github.com/MiroslavZaprazny/gossip-glomers/internal/handlers"
	pkg "github.com/MiroslavZaprazny/gossip-glomers/pkg"
)

func main() {
	n := pkg.NewNode()

	pkg.RegisterInit(n)
	handlers.RegisterEcho(n)
	handlers.RegisterGenerate(n)

	if err := n.Listen(); err != nil {
		log.Fatal(err)
	}
}
