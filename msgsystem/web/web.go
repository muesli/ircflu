package irc

import (
	"github.com/hoisie/web"
	"github.com/muesli/ircflu/app"
	"github.com/muesli/ircflu/msgsystem"
)

type WebSubSystem struct {
	messagesIn chan msgsystem.Message
	messagesOut chan msgsystem.Message

	addr string
}

func (h *WebSubSystem) Name() string {
	return "web"
}

func (h *WebSubSystem) MessageInChan() chan msgsystem.Message {
	return h.messagesIn
}

func (h *WebSubSystem) SetMessageInChan(channel chan msgsystem.Message) {
	h.messagesIn = channel
}

func (h *WebSubSystem) MessageOutChan() chan msgsystem.Message {
	return h.messagesOut
}

func (h *WebSubSystem) SetMessageOutChan(channel chan msgsystem.Message) {
	h.messagesOut = channel
}

func (h *WebSubSystem) Run() {
	go web.Run(h.addr)
}

func init() {
	w := WebSubSystem{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&w.addr, "webaddr", "0.0.0.0:12346", "net.Listen spec, to listen for json-api calls"},
	})

	msgsystem.RegisterSubSystem(&w)
}
