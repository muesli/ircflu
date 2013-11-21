package msgsystem

import (
	"fmt"
)

type MsgSubSystem interface {
	Name() string
	Run(channelIn, channelOut chan Message)
}

type Message struct {
	To     []string
	Msg    string
	Source string
	Authed bool
}

var (
	CommandsIn  = make(chan Message)
	MessagesOut = make(chan Message)

	subsystems map[string]*MsgSubSystem = make(map[string]*MsgSubSystem)
)

func init() {
	fmt.Println("Initializing messaging subsystem...")
}

func RegisterSubSystem(system MsgSubSystem) {
	fmt.Println("Registering msg-subsystem:", system.Name())

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
		(*system).Run(CommandsIn, MessagesOut)
	}
}
