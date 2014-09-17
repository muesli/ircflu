// ircflu's web-hook subsystem.
package hooks

import (
	"fmt"
	"github.com/hoisie/web"
	"github.com/pepl/ircflu/msgsystem"
	_ "strings"
	_ "time"
)

var (
	Hooks []*Hook
)

type Hook interface {
	Request(ctx *web.Context)

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
