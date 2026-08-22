package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/MiroslavZaprazny/gossip-glomers/pkg"
)

type GenerateMessage struct {
	pkg.MessageBody
	Id string `json:"id"`
}

func RegisterGenerate(n *pkg.Node) {
	n.Handle(pkg.MsgGenerate, func(msg *pkg.Message) error {
		var generate GenerateMessage
		if err := json.Unmarshal(msg.Body, &generate); err != nil {
			return err
		}

		msgId := n.NextMsgId()

		if err := n.Reply(msg, &GenerateMessage{
			MessageBody: pkg.MessageBody{
				Type:  pkg.MsgGenerateOk,
				MsgId: msgId,
			},
			Id: fmt.Sprintf("%s-%d", n.NodeId(), msgId),
		}); err != nil {
			return err
		}

		return nil
	})
}
