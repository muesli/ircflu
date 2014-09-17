// ircflu's central commands controller and dispatcher.
package commands

import (
	"fmt"
	"github.com/pepl/ircflu/app"
	"github.com/pepl/ircflu/msgsystem"
	"strings"
)

type Command interface {
	Name() string
	Run(channelIn, channelOut chan msgsystem.Message)
	Parse(msg msgsystem.Message) bool
}

var (
	commands map[string]*Command = make(map[string]*Command)
	enabledCommands map[string]*Command = make(map[string]*Command)

	activateCommands string
)

func init() {
	fmt.Println("Initializing command parsers...")

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&activateCommands, "commands", "alias,auth,join,part,send", "Comma-separated list of commands (alias,auth,exec,join,part,send) you want to enable"},
	})

	go func() {
		for {
			msg := <-msgsystem.CommandsIn
			if len(strings.TrimSpace(msg.Msg)) == 0 {
				fmt.Println("Received empty message. Discarding!")
				continue
			}

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

// Commands need to call this method to register themselves
func RegisterCommand(command Command) {
//	fmt.Println("Registering command:", command.Name())

	commands[command.Name()] = &command
}

// Returns command with this name
func GetCommand(identifier string) *Command {
	command, ok := commands[identifier]
	if ok {
		return command
	}

	return nil
}

// Returns true when command with this name is enabled
func IsCommandEnabled(name string) bool {
	_, ok := enabledCommands[name]
	return ok
}

// Starts all enabled commands
func StartCommands() {
	cmds := strings.Split(activateCommands, ",")
	for _, cmdName := range cmds {
		cmdName = strings.TrimSpace(cmdName)
		command := GetCommand(cmdName)
		if command != nil {
			enabledCommands[(*command).Name()] = command
		} else {
			fmt.Println("Command not found:", cmdName)
		}
	}

	for _, command := range commands {
		if IsCommandEnabled((*command).Name()) {
			fmt.Println("Starting command:", (*command).Name())
			(*command).Run(msgsystem.CommandsIn, msgsystem.MessagesOut)
		}
	}
}
