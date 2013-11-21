// The send command makes ircflu talk on an IRC channel.
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

func (cmd *SendCommand) Name() string {
	return "send"
}

func (cmd *SendCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	command := m[0]
	params := strings.Join(m[1:], " ")

	switch command {
	case "!send":
		if len(params) > 0 {
			r := msgsystem.Message{
				To:  channel,
				Msg: params,
			}
			cmd.messagesOut <- r
		} else {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Usage: !send text",
			}
			cmd.messagesOut <- r
		}
		return true
	}
	return false
}

func (cmd *SendCommand) Run(channelIn, channelOut chan msgsystem.Message) {
	cmd.messagesIn = channelIn
	cmd.messagesOut = channelOut
}

func init() {
	sendCmd := SendCommand{}
	commands.RegisterCommand(&sendCmd)
}
