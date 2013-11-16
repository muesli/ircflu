package msgsystem

import (
	"fmt"
)

type MsgSubSystem interface {
	MessageInChan() chan Message
	SetMessageInChan(channel chan Message)

	MessageOutChan() chan Message
	SetMessageOutChan(channel chan Message)

	Name() string
	Run()
}

type Message struct {
	To  []string
	Msg string
	Source string
	Authed bool
}

var (
	CommandsIn = make(chan Message)
	MessagesOut = make(chan Message)

	subsystems map[string]*MsgSubSystem = make(map[string]*MsgSubSystem)
)

func init() {
	fmt.Println("Initializing messaging subsystem...")
}

func RegisterSubSystem(system MsgSubSystem) {
	fmt.Println("Registering msg-subsystem:", system.Name())

	system.SetMessageInChan(CommandsIn)
	system.SetMessageOutChan(MessagesOut)

	subsystems[system.Name()] = &system
}

func SubSystem(identifier string) *MsgSubSystem {
	system, ok := subsystems[identifier]
	if ok {
		return system
	}

	return nil
}

func StartSubSystems() {
	for _, system := range subsystems {
		(*system).Run()
	}
}
