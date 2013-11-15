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
}

var (
	MessagesIn = make(chan Message)
	MessagesOut = make(chan Message)

	subsystems []*MsgSubSystem
)

func init() {
	fmt.Println("Initializing messaging subsystem...")
}

func RegisterSubSystem(system MsgSubSystem) {
	fmt.Println("Registering msg-subsystem:", system.Name())

	system.SetMessageInChan(MessagesIn)
	system.SetMessageOutChan(MessagesOut)

	subsystems = append(subsystems, &system)
}

func StartSubSystems() {
	for _, system := range subsystems {
		(*system).Run()
	}
}
