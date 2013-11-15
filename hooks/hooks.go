package hooks

import (
	"fmt"
	_ "strings"
	_ "time"
	"github.com/hoisie/web"
)

var (
	Messages = make(chan string)
	Hooks []*Hook
)

type Hook interface {
	Request(ctx *web.Context)

	MessageChan() chan string
	SetMessageChan(channel chan string)

	Name() string
	Path() string
}

func init() {
	fmt.Println("Initializing hooks subsystem...")
}

func RegisterWebHook(hook Hook) {
	fmt.Println("Registering web-hook:", hook.Name(), "on", hook.Path())

	hook.SetMessageChan(Messages)
	Hooks = append(Hooks, &hook)

	web.Post(hook.Path(), hook.Request)
}
