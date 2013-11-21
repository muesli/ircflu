package joinCmd

import (
	"fmt"
	"github.com/muesli/ircflu/commands"
	"github.com/muesli/ircflu/msgsystem"
	"github.com/muesli/ircflu/msgsystem/irc"
	_ "log"
	"strings"
)

type JoinCommand struct {
	messagesIn  chan msgsystem.Message
	messagesOut chan msgsystem.Message
}

func (h *JoinCommand) Name() string {
	return "join"
}

func (h *JoinCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	cmd := m[0]
	params := strings.Join(m[1:], " ")

	switch cmd {
	case "!join":
		if !msg.Authed {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Security breach. Talk to ircflu admin!",
			}
			h.messagesOut <- r
			return true
		}

		if len(params) > 0 {
			fmt.Println("Joining:", params)

			ircclient := (*msgsystem.SubSystem("irc")).(*irc.IrcSubSystem)
			if ircclient != nil {
				ircclient.Join(params)
			}
		} else {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Usage: !join #chan  or  !join #chan key",
			}
			h.messagesOut <- r
		}
		return true
	}
	return false
}

func (h *JoinCommand) Run(channelIn, channelOut chan msgsystem.Message) {
	h.messagesIn = channelIn
	h.messagesOut = channelOut
}

func init() {
	joinCmd := JoinCommand{}
	commands.RegisterCommand(&joinCmd)
}
