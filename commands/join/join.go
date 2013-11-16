package joinCmd

import (
	"fmt"
	_ "log"
	"strings"
	"ircflu/auth"
	"ircflu/commands"
	"ircflu/msgsystem"
	"ircflu/msgsystem/irc"
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
			if !auth.IsAuthed(msg.Source) {
				r := msgsystem.Message{
					To: channel,
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
					To: channel,
					Msg: "Usage: !join #chan  or  !join #chan key",
				}
				h.messagesOut <- r
			}
			return true
	}
	return false
}

func init() {
	joinCmd := JoinCommand{}
	commands.RegisterCommand(&joinCmd)
}
