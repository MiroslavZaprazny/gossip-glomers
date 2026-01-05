package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Message struct {
	Src  string          `json:"src"`
	Dest string          `json:"dest"`
	Body json.RawMessage `json:"body"`
}

type MessageBody struct {
	Type      string `json:"type"`
	MsgId     int    `json:"msg_id,omitempty"` //unique only to the local node
	InReplyTo int    `json:"in_reply_to,omitempty"`
}

type InitMessageBody struct {
	NodeId  string   `json:"node_id"`
	NodeIds []string `json:"node_ids"`
}

type Replyable interface {
	SetInReplyTo(msgId int)
}

func (mb *MessageBody) SetInReplyTo(msgId int) {
	mb.InReplyTo = msgId
}

type Handler func(msg Message) error

type Node struct {
	id       string
	nodeIds  []string
	handlers map[string]Handler
	enc      *json.Encoder
}

func NewNode() *Node {
	return &Node{
		handlers: make(map[string]Handler),
		enc:      json.NewEncoder(os.Stdout),
	}
}

func (n *Node) Handle(msgType string, handler Handler) {
	n.handlers[msgType] = handler
}

func (n *Node) Listen() error {
	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		raw := r.Bytes()

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			return err
		}

		var body MessageBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if body.Type == "init" {
			var init InitMessageBody

			if err := json.Unmarshal(msg.Body, &init); err != nil {
				return err
			}

			n.id = init.NodeId
			n.nodeIds = init.NodeIds

			if err := n.Reply(msg, &MessageBody{Type: "init_ok"}); err != nil {
				return err
			}
		} else {
			handler, ok := n.handlers[body.Type]

			if !ok {
				fmt.Fprintf(os.Stderr, "Unhandled msg type: %s", body.Type)

				os.Exit(1)
			}

			if err := handler(msg); err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *Node) Reply(originalMsg Message, body Replyable) error {
	var msgBody MessageBody
	if err := json.Unmarshal(originalMsg.Body, &msgBody); err != nil {
		return err
	}

	body.SetInReplyTo(msgBody.MsgId)

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	if err := n.enc.Encode(Message{Src: originalMsg.Dest, Dest: originalMsg.Src, Body: bodyJson}); err != nil {
		return err
	}

	return nil
}
