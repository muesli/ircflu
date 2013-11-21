package commands

import (
	"fmt"
	"github.com/muesli/ircflu/msgsystem"
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
)

func init() {
	fmt.Println("Initializing command parsers...")

	go func() {
		for {
			msg := <-msgsystem.CommandsIn
			fmt.Println("Commands:", msg.To, msg.Msg)

			go func() {
				for _, c := range commands {
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
