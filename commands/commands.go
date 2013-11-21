package commands

import (
	"fmt"
	"github.com/muesli/ircflu/app"
	"github.com/muesli/ircflu/msgsystem"
	"strings"
)

type Command interface {
	MessageInChan() chan msgsystem.Message
	SetMessageInChan(channel chan msgsystem.Message)
	MessageOutChan() chan msgsystem.Message
	SetMessageOutChan(channel chan msgsystem.Message)

	Name() string

	Parse(msg msgsystem.Message) bool
}

var (
	commands []*Command
	enabledCommands string
)

func IsCommandEnabled(name string) bool {
	foundAsEnabled := false
	cmds := strings.Split(enabledCommands, ",")
	for _, cmdName := range cmds {
		cmdName = strings.TrimSpace(cmdName)
		if cmdName == name {
			foundAsEnabled = true
		}
	}

	return foundAsEnabled
}

func init() {
	fmt.Println("Initializing command parsers...")

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&enabledCommands, "commands", "alias,auth,join,part,send", "Comma-separated list of commands (alias,auth,exec,join,part,send) you want to enable"},
	})

	go func() {
		for {
			msg := <-msgsystem.CommandsIn
			fmt.Println("Commands:", msg.To, msg.Msg)

			go func() {
				for _, c := range commands {
					if !IsCommandEnabled((*c).Name()) {
						continue
					}
					fmt.Println("Handing out to:", (*c).Name(), (*c).Parse(msg))
				}
			}()
		}
	}()
}

func RegisterCommand(command Command) {
	fmt.Println("Registering command:", command.Name())
	command.SetMessageInChan(msgsystem.CommandsIn)
	command.SetMessageOutChan(msgsystem.MessagesOut)
	commands = append(commands, &command)
}
