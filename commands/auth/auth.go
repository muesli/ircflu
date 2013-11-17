package authCmd

import (
	_ "fmt"
	_ "log"
	"strings"
	"github.com/muesli/ircflu/auth"
	"github.com/muesli/ircflu/commands"
	"github.com/muesli/ircflu/msgsystem"
)

type AuthCommand struct {
	messagesIn chan msgsystem.Message
	messagesOut chan msgsystem.Message
}

func (h *AuthCommand) Name() string {
	return "auth"
}

func (h *AuthCommand) MessageInChan() chan msgsystem.Message {
	return h.messagesIn
}

func (h *AuthCommand) SetMessageInChan(channel chan msgsystem.Message) {
	h.messagesIn = channel
}

func (h *AuthCommand) MessageOutChan() chan msgsystem.Message {
	return h.messagesOut
}

func (h *AuthCommand) SetMessageOutChan(channel chan msgsystem.Message) {
	h.messagesOut = channel
}

func (h *AuthCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	cmd := m[0]
	params := strings.Join(m[1:], " ")

	switch cmd {
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

				h.messagesOut <- r
			} else {
				r := msgsystem.Message{
					To: channel,
					Msg: "Usage in private query only: !auth password",
				}
				h.messagesOut <- r
			}
			return true
	}
	return false
}

func init() {
	authCmd := AuthCommand{}
	commands.RegisterCommand(&authCmd)
}
