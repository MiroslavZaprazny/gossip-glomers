package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Data types implmeneting the maelstrom protocol
// see: https://github.com/jepsen-io/maelstrom/blob/main/doc/protocol.md

type Message struct {
	Src  string          `json:"src"` // A string identifying the node this message came from
	Dest string          `json:"dest"` // A string identifying the node this message is to
	Body json.RawMessage `json:"body"` // An object: the payload of the message
}

type MessageBody struct {
	Type      MsgType `json:"type"` // A string identifying the type of message this is
	MsgId     int    `json:"msg_id,omitempty"` // A unique integer identifier(is only unique to the local node)
	InReplyTo int    `json:"in_reply_to,omitempty"` // For req/response, the msg_id of the request
}

type MsgType string

const (
	MsgInit MsgType = "init"
	MsgInitOk MsgType = "init_ok"
	MsgEcho MsgType = "echo"
	MsgEchoOk MsgType = "echo_ok"
	MsgError MsgType = "error"
)

type Replyable interface {
	SetInReplyTo(msgId int)
}

func (mb *MessageBody) SetInReplyTo(msgId int) {
	mb.InReplyTo = msgId
}

type Handler func(msg *Message) error

type Node struct {
	id       string
	nodeIds  []string
	handlers map[MsgType]Handler
	enc      *json.Encoder
}

func NewNode() *Node {
	return &Node{
		handlers: make(map[MsgType]Handler),
		enc:      json.NewEncoder(os.Stdout),
	}
}

func (n *Node) Init(id string, nodeIds []string) {
	n.id = id
	n.nodeIds = nodeIds
}

func (n *Node) Handle(msgType MsgType, handler Handler) {
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

		handler, ok := n.handlers[body.Type]

		if !ok {
			//Todo: maybe we should reply with code 10 - not-supported
			fmt.Fprintf(os.Stderr, "Unhandled msg type: %s", body.Type)
			continue
		}

		if err := handler(&msg); err != nil {
			return err
		}
	}

	if r.Err() != nil {
		//Todo: maybe we should reply with some error code?
		fmt.Fprintf(os.Stderr, "Scanner failed with error: %s", r.Err().Error())
	}

	return nil
}

func (n *Node) Reply(to *Message, body Replyable) error {
	var msgBody MessageBody
	if err := json.Unmarshal(to.Body, &msgBody); err != nil {
		return err
	}

	body.SetInReplyTo(msgBody.MsgId)

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	if err := n.enc.Encode(Message{Src: to.Dest, Dest: to.Src, Body: bodyJson}); err != nil {
		return err
	}

	return nil
}
