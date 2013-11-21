// ircflu's integrated web-server to handle web-hooks.
package irc

import (
	"github.com/hoisie/web"
	"github.com/muesli/ircflu/app"
	"github.com/muesli/ircflu/msgsystem"
)

type WebSubSystem struct {
	addr string
}

func (h *WebSubSystem) Name() string {
	return "web"
}

func (h *WebSubSystem) Run(channelIn, channelOut chan msgsystem.Message) {
	go web.Run(h.addr)
}

func init() {
	w := WebSubSystem{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&w.addr, "webaddr", "0.0.0.0:12346", "net.Listen spec, to listen for json-api calls"},
	})

	msgsystem.RegisterSubSystem(&w)
}
