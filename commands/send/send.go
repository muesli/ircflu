package sendCmd

import (
	_ "fmt"
	"github.com/muesli/ircflu/commands"
	"github.com/muesli/ircflu/msgsystem"
	_ "log"
	"strings"
)

type SendCommand struct {
	messagesIn  chan msgsystem.Message
	messagesOut chan msgsystem.Message
}

func (h *SendCommand) Name() string {
	return "send"
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
				To:  channel,
				Msg: params,
			}
			h.messagesOut <- r
		} else {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Usage: !send text",
			}
			h.messagesOut <- r
		}
		return true
	}
	return false
}

func (h *SendCommand) Run(channelIn, channelOut chan msgsystem.Message) {
	h.messagesIn = channelIn
	h.messagesOut = channelOut
}

func init() {
	sendCmd := SendCommand{}
	commands.RegisterCommand(&sendCmd)
}
