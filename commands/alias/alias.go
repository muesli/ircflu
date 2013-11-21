// The alias command lets you setup shortcuts to lengthy commands.
package aliasCmd

import (
	"fmt"
	"github.com/muesli/ircflu/commands"
	"github.com/muesli/ircflu/msgsystem"
	"github.com/muesli/ircflu/msgsystem/irc/irctools"
	_ "log"
	"strings"
)

type AliasCommand struct {
	messagesIn  chan msgsystem.Message
	messagesOut chan msgsystem.Message

	aliases map[string]string
}

func (h *AliasCommand) Name() string {
	return "alias"
}

func (h *AliasCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	cmd := m[0]
	params := strings.Join(m[1:], " ")

	switch cmd {
	case "!alias":
		if !msg.Authed {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Security breach. Talk to ircflu admin!",
			}
			h.messagesOut <- r
			return true
		}

		a := strings.Split(params, "=")
		if len(a) == 2 {
			a[0] = strings.TrimSpace(a[0])
			a[1] = strings.TrimSpace(a[1])
			h.aliases[a[0]] = a[1]
			r := msgsystem.Message{
				To:  channel,
				Msg: "Added new alias '" + a[0] + "' for command '" + a[1] + "'",
			}
			h.messagesOut <- r
		} else {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Usage: !alias [new command] = [actual command]",
			}
			h.messagesOut <- r

			for k, v := range h.aliases {
				r.Msg = "Alias: " + irctools.Colored(k, "red") + " = " + irctools.Colored(v, "teal")
				h.messagesOut <- r
			}
		}

		return true

	default:
		v, ok := h.aliases[strings.TrimSpace(msg.Msg)[1:]]
		if ok {
			fmt.Println("Alias:", v, strings.TrimSpace(msg.Msg), ok)
			r := msgsystem.Message{
				To:     channel,
				Source: msg.Source,
				Authed: msg.Authed,
				Msg:    "!" + v,
			}
			h.messagesIn <- r

			return true
		}
	}

	return false
}

func (h *AliasCommand) Run(channelIn, channelOut chan msgsystem.Message) {
	h.messagesIn = channelIn
	h.messagesOut = channelOut
}

func init() {
	aliasCmd := AliasCommand{
		aliases: make(map[string]string),
	}
	commands.RegisterCommand(&aliasCmd)
}
