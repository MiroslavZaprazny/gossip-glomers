package handlers

import (
	"encoding/json"

	"github.com/MiroslavZaprazny/gossip-glomers/pkg"
)

type EchoReplyBody struct {
	BodyType string `json:"type"`
	MsgId int `json:"msg_id"`
	InReplyTo int `json:"in_reply_to"`
	Echo string `json:"echo"`
}

type EchoReplyMessage struct {
	BodyType string `json:"type"`
	MsgId int `json:"msg_id"`
	Echo string `json:"echo"`
}

func (b *EchoReplyBody) SetInReplyTo(msgId int) {
	b.InReplyTo = msgId
}

func RegisterEcho(n *pkg.Node) {
	n.Handle("echo", func (msg *pkg.Message) error{
		var msgBody EchoReplyMessage
		if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
			return err
		}

		n.Reply(msg, &EchoReplyBody{
			BodyType: "echo_ok",
			MsgId: msgBody.MsgId,
			Echo: msgBody.Echo,
		})

		return nil
	})
}

