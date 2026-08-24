package handlers

import (
	"encoding/json"
	"sync"

	maelstrom "github.com/MiroslavZaprazny/gossip-glomers/maelstrom"
)

type NodeMessages struct {
	inner []int
	mu    sync.RWMutex
}

type ClusterTopology struct {
	mu       sync.RWMutex
	topology map[string][]string // an array of neighbours for each node
}

type BroadcastMessage struct {
	maelstrom.MessageBody
	Message int `json:"message"`
}

type ReadMessage struct {
	maelstrom.MessageBody
	Messages []int `json:"messages"`
}

type TopologyMessage struct {
	maelstrom.MessageBody
	Topology map[string][]string `json:"topology"`
}

func (tm *ClusterTopology) write(topology map[string][]string) {
	tm.mu.Lock()
	tm.topology = topology
	tm.mu.Unlock()
}

func (m *NodeMessages) appendMessage(message int) {
	m.mu.Lock()
	m.inner = append(m.inner, message)
	m.mu.Unlock()
}

func (m *NodeMessages) snapshot() []int {
	m.mu.RLock()
	msgCopy := make([]int, len(m.inner))
	copy(msgCopy, m.inner)
	m.mu.RUnlock()

	return msgCopy
}

func RegisterBroadcast(n *maelstrom.Node, messages *NodeMessages) {
	n.Handle(maelstrom.MsgBroadcast, func(msg *maelstrom.Message) error {
		var broadcast BroadcastMessage
		if err := json.Unmarshal(msg.Body, &broadcast); err != nil {
			return err
		}
		messages.appendMessage(broadcast.Message)

		if err := n.Reply(msg, &maelstrom.MessageBody{
			Type:  maelstrom.MsgBroadcastOk,
			MsgId: n.NextMsgId(),
		}); err != nil {
			return err
		}

		return nil
	})
}

func RegisterRead(n *maelstrom.Node, messages *NodeMessages) {
	n.Handle(maelstrom.MsgRead, func(msg *maelstrom.Message) error {

		if err := n.Reply(msg, &ReadMessage{
			MessageBody: maelstrom.MessageBody{
				Type:  maelstrom.MsgReadOk,
				MsgId: n.NextMsgId(),
			},
			Messages: messages.snapshot(),
		}); err != nil {
			return err
		}

		return nil
	})
}

func RegisterTopology(n *maelstrom.Node, clusterTopology *ClusterTopology) {
	n.Handle(maelstrom.MsgTopology, func(msg *maelstrom.Message) error {
		var topologyMessage TopologyMessage
		if err := json.Unmarshal(msg.Body, &topologyMessage); err != nil {
			return err
		}

		clusterTopology.write(topologyMessage.Topology)

		if err := n.Reply(msg, &maelstrom.MessageBody{
			Type:  maelstrom.MsgTopologyOk,
			MsgId: n.NextMsgId(),
		}); err != nil {
			return err
		}

		return nil
	})
}
