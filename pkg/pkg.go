package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Message struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
	Body Body   `json:"body"`
}

type Body struct {
	Type      string `json:"type"`
	MsgId     int    `json:"msg_id,omitempty"` //unique only to the local node
	InReplyTo int    `json:"in_reply_to,omitempty"`
}

type InitRequest struct {
	NodeId  string   `json:"node_id"`
	NodeIds []string `json:"node_ids"`
}

func MainLoop() {
	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		raw := r.Bytes()
		var msg Message

		if err := json.Unmarshal(raw, &msg); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to unmarshal message: %s, err: %s", string(raw), err)

			os.Exit(1)
		}

		switch msg.Body.Type {
		case "init":
			var msg Message

			if err := json.Unmarshal(raw, &msg); err != nil {
				fmt.Fprintf(os.Stderr, "Unable to unmarshal init request: %s, err: %s", string(raw), err)

				os.Exit(1)
			}

			response := Message{
				Src:  msg.Dest,
				Dest: msg.Src,
				Body: Body{
					Type:      "init_ok",
					InReplyTo: msg.Body.MsgId,
				},
			}

			if err := Reply(response); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write %s, to stdout err: %s", response.Body.Type, err)

				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unimplemented msg type: %s", msg.Body.Type)

			os.Exit(1)
		}
	}
}

func Reply(msg any) error {
	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(msg); err != nil {
		return err
	}

	return nil
}
