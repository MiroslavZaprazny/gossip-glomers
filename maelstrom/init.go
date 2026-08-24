package maelstorm

import "encoding/json"

type InitMessageBody struct {
	NodeId  string   `json:"node_id"`
	NodeIds []string `json:"node_ids"`
}

func RegisterInit(n *Node) {
	n.Handle(MsgInit, func(msg *Message) error {
		var init InitMessageBody

		if err := json.Unmarshal(msg.Body, &init); err != nil {
			return err
		}

		n.Init(init.NodeId, init.NodeIds)

		if err := n.Reply(msg, &MessageBody{Type: MsgInitOk}); err != nil {
			return err
		}

		return nil
	})
}
