package execCmd

import (
	"os/exec"
	"fmt"
	_ "log"
	"strings"
	"ircflu/auth"
	"ircflu/commands"
	"ircflu/msgsystem"
)

type ExecCommand struct {
	messagesIn chan msgsystem.Message
	messagesOut chan msgsystem.Message
}

func (h *ExecCommand) Name() string {
	return "exec"
}

func (h *ExecCommand) MessageInChan() chan msgsystem.Message {
	return h.messagesIn
}

func (h *ExecCommand) SetMessageInChan(channel chan msgsystem.Message) {
	h.messagesIn = channel
}

func (h *ExecCommand) MessageOutChan() chan msgsystem.Message {
	return h.messagesOut
}

func (h *ExecCommand) SetMessageOutChan(channel chan msgsystem.Message) {
	h.messagesOut = channel
}

func (h *ExecCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	cmd := m[0]
	params := strings.TrimSpace(strings.Join(m[1:], " "))

	switch cmd {
		case "!exec":
			if !auth.IsAuthed(msg.Source) || strings.Index(params, "rm ") >= 0 || strings.Index(params, "mv ") >= 0 {
				r := msgsystem.Message{
					To: channel,
					Msg: "Security breach. Talk to ircflu admin!",
				}
				h.messagesOut <- r
				return true
			}

			if len(params) > 0 {
				r := msgsystem.Message{
					To: channel,
					Msg: "Executing command!",
				}
				h.messagesOut <- r

				c := strings.Split(params, " ")
				e := exec.Command(c[0], c[1:]...)
				out, err := e.CombinedOutput()
				fmt.Println("Output:", string(out))
				fmt.Println("Error:", err)

				r = msgsystem.Message{
					To: []string{msg.Source},
				}
				if err != nil {
					r.Msg = "Command '" + params + "' failed: " + err.Error()
				} else {
					r.Msg = "Command '" + params + "' succeeded!"
				}

				h.messagesOut <- r
			} else {
				r := msgsystem.Message{
					To: channel,
					Msg: "Usage: !exec [command]",
				}
				h.messagesOut <- r
			}

			return true
	}

	return false
}

func init() {
	execCmd := ExecCommand{}
	commands.RegisterCommand(&execCmd)
}
