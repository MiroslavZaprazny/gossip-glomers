package handlers

import (
	"encoding/json"

	"github.com/MiroslavZaprazny/gossip-glomers/pkg"
)

type InitMessageBody struct {
	NodeId  string   `json:"node_id"`
	NodeIds []string `json:"node_ids"`
}

func RegisterInit(n *pkg.Node) {
	n.Handle(pkg.MsgInit, func (msg *pkg.Message) error {
		var init InitMessageBody

		if err := json.Unmarshal(msg.Body, &init); err != nil {
			return err
		}

		n.Init(init.NodeId, init.NodeIds)

		if err := n.Reply(msg, &pkg.MessageBody{Type: pkg.MsgInitOk}); err != nil {
			return err
		}

		return nil
	})
}

