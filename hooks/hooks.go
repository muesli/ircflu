package hooks

import (
	"fmt"
	_ "strings"
	_ "time"
	"github.com/hoisie/web"
	"ircflu/msgsystem"
)

var (
	Hooks []*Hook
)

type Hook interface {
	Request(ctx *web.Context)

	MessageChan() chan msgsystem.Message
	SetMessageChan(channel chan msgsystem.Message)

	Name() string
	Path() string
}

func init() {
	fmt.Println("Initializing hooks subsystem...")
}

func RegisterWebHook(hook Hook) {
	fmt.Println("Registering web-hook:", hook.Name(), "on", hook.Path())

	hook.SetMessageChan(msgsystem.MessagesOut)
	Hooks = append(Hooks, &hook)

	web.Post(hook.Path(), hook.Request)
}
