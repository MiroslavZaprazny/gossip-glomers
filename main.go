package main

import (
	"log"

	handlers "github.com/MiroslavZaprazny/gossip-glomers/internal/handlers"
	maelstrom "github.com/MiroslavZaprazny/gossip-glomers/maelstrom"
)

func main() {
	n := maelstrom.NewNode()

	//TODO: this should probably be registered in the maelstrom lib
	maelstrom.RegisterInit(n)

	handlers.RegisterEcho(n)
	handlers.RegisterGenerate(n)

	nodeMsgs := &handlers.NodeMessages{}
	topology := &handlers.ClusterTopology{}

	handlers.RegisterBroadcast(n, nodeMsgs)
	handlers.RegisterRead(n, nodeMsgs)
	handlers.RegisterTopology(n, topology)

	if err := n.Listen(); err != nil {
		log.Fatal(err)
	}
}
