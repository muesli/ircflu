package sendCmd

import (
	_ "fmt"
	_ "log"
	"strings"
	"ircflu/commands"
	"ircflu/msgsystem"
)

type SendCommand struct {
	messagesIn chan msgsystem.Message
	messagesOut chan msgsystem.Message
}

func (h *SendCommand) Name() string {
	return "join"
}

func (h *SendCommand) MessageInChan() chan msgsystem.Message {
	return h.messagesIn
}

func (h *SendCommand) SetMessageInChan(channel chan msgsystem.Message) {
	h.messagesIn = channel
}

func (h *SendCommand) MessageOutChan() chan msgsystem.Message {
	return h.messagesOut
}

func (h *SendCommand) SetMessageOutChan(channel chan msgsystem.Message) {
	h.messagesOut = channel
}

func (h *SendCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	cmd := m[0]
	params := strings.Join(m[1:], " ")

	switch cmd {
		case "!send":
			if len(params) > 0 {
				r := msgsystem.Message{
					Msg: params,
				}
				h.messagesOut <- r
			} else {
				r := msgsystem.Message{
					To: channel,
					Msg: "Usage: !send text",
				}
				h.messagesOut <- r
			}
			return true
	}
	return false
}

func init() {
	sendCmd := SendCommand{}
	commands.RegisterCommand(&sendCmd)
}
