package handlers

import (
	"encoding/json"

	maelstrom "github.com/MiroslavZaprazny/gossip-glomers/maelstrom"
)

type EchoMessage struct {
	maelstrom.MessageBody
	Echo string `json:"echo"`
}

func RegisterEcho(n *maelstrom.Node) {
	n.Handle(maelstrom.MsgEcho, func(msg *maelstrom.Message) error {
		var echo EchoMessage
		if err := json.Unmarshal(msg.Body, &echo); err != nil {
			return err
		}

		if err := n.Reply(msg, &EchoMessage{
			MessageBody: maelstrom.MessageBody{
				Type:  maelstrom.MsgEchoOk,
				MsgId: n.NextMsgId(),
			},
			Echo: echo.Echo,
		}); err != nil {
			return err
		}

		return nil
	})
}
