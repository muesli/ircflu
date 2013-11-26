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

func (sys *WebSubSystem) Name() string {
	return "web"
}

func (sys *WebSubSystem) Run(channelIn, channelOut chan msgsystem.Message) {
	go web.Run(sys.addr)
}

func (sys *WebSubSystem) Handle(cm msgsystem.Message) bool {
	return false
}

func init() {
	w := WebSubSystem{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&w.addr, "webaddr", "0.0.0.0:12346", "net.Listen spec, to listen for json-api calls"},
	})

	msgsystem.RegisterSubSystem(&w)
}
