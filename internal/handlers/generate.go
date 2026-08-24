package handlers

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/MiroslavZaprazny/gossip-glomers/maelstrom"
)

type GenerateMessage struct {
	maelstrom.MessageBody
	Id string `json:"id"`
}

func RegisterGenerate(n *maelstrom.Node) {
	n.Handle(maelstrom.MsgGenerate, func(msg *maelstrom.Message) error {
		var generate GenerateMessage
		if err := json.Unmarshal(msg.Body, &generate); err != nil {
			return err
		}

		msgId := n.NextMsgId()

		if err := n.Reply(msg, &GenerateMessage{
			MessageBody: maelstrom.MessageBody{
				Type:  maelstrom.MsgGenerateOk,
				MsgId: msgId,
			},
			Id: fmt.Sprintf("%s-%d", n.NodeId(), msgId),
		}); err != nil {
			return err
		}

		return nil
	})
}
