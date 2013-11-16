package joinCmd

import (
	"fmt"
	_ "log"
	"strings"
	"ircflu/commands"
	"ircflu/msgsystem"
)

type JoinCommand struct {
	messagesIn chan msgsystem.Message
	messagesOut chan msgsystem.Message
}

func (h *JoinCommand) Name() string {
	return "join"
}

func (h *JoinCommand) MessageInChan() chan msgsystem.Message {
	return h.messagesIn
}

func (h *JoinCommand) SetMessageInChan(channel chan msgsystem.Message) {
	h.messagesIn = channel
}

func (h *JoinCommand) MessageOutChan() chan msgsystem.Message {
	return h.messagesOut
}

func (h *JoinCommand) SetMessageOutChan(channel chan msgsystem.Message) {
	h.messagesOut = channel
}

func (h *JoinCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	cmd := m[0]
	params := strings.Join(m[1:], " ")

	switch cmd {
		case "!join":
			if len(params) > 0 {
				fmt.Println("Joining:", params)
//				h.client.Join(params)
			} else {
				r := msgsystem.Message{
					To: channel,
					Msg: "Usage: !join #chan  or  !join #chan key",
				}
				h.messagesOut <- r
			}
			return true

		case "!part":
			if len(params) > 0 {
				fmt.Println("Parting:", params)
//				h.client.Part(params)
			} else {
				r := msgsystem.Message{
					To: channel,
					Msg: "Usage: !part #chan",
				}
				h.messagesOut <- r
			}
			return true

		default:
			if !strings.HasPrefix(cmd, "!") {
				h.messagesOut <- msg
			}
	}
	return false
}

func init() {
	joinCmd := JoinCommand{}
	commands.RegisterCommand(&joinCmd)
}
