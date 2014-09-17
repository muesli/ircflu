// The auth command lets clients authenticate to unlock secured commands.
package authCmd

import (
	_ "fmt"
	"github.com/pepl/ircflu/auth"
	"github.com/pepl/ircflu/commands"
	"github.com/pepl/ircflu/msgsystem"
	_ "log"
	"strings"
)

type AuthCommand struct {
	messagesIn  chan msgsystem.Message
	messagesOut chan msgsystem.Message
}

func (cmd *AuthCommand) Name() string {
	return "auth"
}

func (cmd *AuthCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	command := m[0]
	params := strings.Join(m[1:], " ")

	switch command {
	case "!auth":
		if len(params) > 0 && !strings.HasPrefix(channel[0], "#") {
			r := msgsystem.Message{
				To: []string{msg.Source},
			}

			if auth.Auth(msg.Source, params) {
				r.Msg = "Auth succeeded!"
			} else {
				r.Msg = "Auth failed!"
			}

			cmd.messagesOut <- r
		} else {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Usage in private query only: !auth password",
			}
			cmd.messagesOut <- r
		}
		return true
	}
	return false
}

func (cmd *AuthCommand) Run(channelIn, channelOut chan msgsystem.Message) {
	cmd.messagesIn = channelIn
	cmd.messagesOut = channelOut
}

func init() {
	authCmd := AuthCommand{}
	commands.RegisterCommand(&authCmd)
}
