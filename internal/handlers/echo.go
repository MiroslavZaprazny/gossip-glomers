package handlers

import (
	"encoding/json"

	"github.com/MiroslavZaprazny/gossip-glomers/pkg"
)

type EchoMessage struct {
	pkg.MessageBody
	Echo string `json:"echo"`
}

func RegisterEcho(n *pkg.Node) {
	n.Handle(pkg.MsgEcho, func (msg *pkg.Message) error{
		var echo EchoMessage
		if err := json.Unmarshal(msg.Body, &echo); err != nil {
			return err
		}

		if err := n.Reply(msg, &EchoMessage{
			MessageBody: pkg.MessageBody{
				Type: pkg.MsgEchoOk,
				MsgId: echo.MsgId,
			},
			Echo: echo.Echo,
		}); err != nil {
			//todo do something?
		}

		return nil
	})
}

