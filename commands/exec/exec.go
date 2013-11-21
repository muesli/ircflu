// The exec command lets you remotely execute arbitrary commands.
package execCmd

import (
	"fmt"
	"github.com/muesli/ircflu/commands"
	"github.com/muesli/ircflu/msgsystem"
	"github.com/muesli/ircflu/msgsystem/irc/irctools"
	_ "log"
	"os/exec"
	"strings"
)

type ExecCommand struct {
	messagesIn  chan msgsystem.Message
	messagesOut chan msgsystem.Message
}

func (cmd *ExecCommand) Name() string {
	return "exec"
}

func (cmd *ExecCommand) Parse(msg msgsystem.Message) bool {
	channel := msg.To
	m := strings.Split(msg.Msg, " ")
	command := m[0]
	params := strings.TrimSpace(strings.Join(m[1:], " "))

	switch command {
	case "!exec":
		if !msg.Authed || strings.Index(params, "rm ") >= 0 || strings.Index(params, "mv ") >= 0 {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Security breach. Talk to ircflu admin!",
			}
			cmd.messagesOut <- r
			return true
		}

		if len(params) > 0 {
			r := msgsystem.Message{
				To:  channel,
				Msg: irctools.Colored("Executing command!", "red"),
			}
			cmd.messagesOut <- r

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

			cmd.messagesOut <- r
		} else {
			r := msgsystem.Message{
				To:  channel,
				Msg: "Usage: !exec [command]",
			}
			cmd.messagesOut <- r
		}

		return true
	}

	return false
}

func (cmd *ExecCommand) Run(channelIn, channelOut chan msgsystem.Message) {
	cmd.messagesIn = channelIn
	cmd.messagesOut = channelOut
}

func init() {
	execCmd := ExecCommand{}
	commands.RegisterCommand(&execCmd)
}
