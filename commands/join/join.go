// The join command makes ircflu join another channel.
package joinCmd

import (
	"fmt"
	"github.com/pepl/ircflu/commands"
	"github.com/pepl/ircflu/msgsystem"
	"github.com/pepl/ircflu/msgsystem/irc"
	_ "log"
	"strings"
)

type JoinCommand struct {
	messagesIn  chan msgsystem.Message
	messagesOut chan msgsystem.Message
}

func (cmd *JoinCommand) Name() string {
	return "join"
}

func (cmd *JoinCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	command := m[0]
	params := strings.Join(m[1:], " ")

	switch command {
	case "!join":
		if !msg.Authed {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Security breach. Talk to ircflu admin!",
			}
			cmd.messagesOut <- r
			return true
		}

		if len(params) > 0 {
			fmt.Println("Joining:", params)

			ircclient := (*msgsystem.GetSubSystem("irc")).(*irc.IrcSubSystem)
			if ircclient != nil {
				ircclient.Join(params)
			}
		} else {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Usage: !join #chan  or  !join #chan key",
			}
			cmd.messagesOut <- r
		}
		return true
	}
	return false
}

func (cmd *JoinCommand) Run(channelIn, channelOut chan msgsystem.Message) {
	cmd.messagesIn = channelIn
	cmd.messagesOut = channelOut
}

func init() {
	joinCmd := JoinCommand{}
	commands.RegisterCommand(&joinCmd)
}
