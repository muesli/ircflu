package irc

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/muesli/ircflu/app"
	"github.com/muesli/ircflu/auth"
	"github.com/muesli/ircflu/msgsystem"
	"log"
	"strings"
	"time"
)

type IrcSubSystem struct {
	messagesIn  chan msgsystem.Message
	messagesOut chan msgsystem.Message

	// channel signaling irc connection status
	ConnectedState chan bool

	// setup IRC client:
	client *irc.Conn

	irchost     string
	ircnick     string
	ircpassword string
	ircssl      bool
	ircchannel  string

	channels    []string
}

func (h *IrcSubSystem) Name() string {
	return "irc"
}

func (h *IrcSubSystem) MessageInChan() chan msgsystem.Message {
	return h.messagesIn
}

func (h *IrcSubSystem) SetMessageInChan(channel chan msgsystem.Message) {
	h.messagesIn = channel
}

func (h *IrcSubSystem) MessageOutChan() chan msgsystem.Message {
	return h.messagesOut
}

func (h *IrcSubSystem) SetMessageOutChan(channel chan msgsystem.Message) {
	h.messagesOut = channel
}

func (h *IrcSubSystem) Rejoin() {
	for _, channel := range h.channels {
		h.client.Join(channel)
	}
}

func (h *IrcSubSystem) Join(channel string) {
	channel = strings.TrimSpace(channel)
	h.client.Join(channel)

	h.channels = append(h.channels, channel)
}

func (h *IrcSubSystem) Part(channel string) {
	channel = strings.TrimSpace(channel)
	h.client.Part(channel)

	for k, v := range h.channels {
		if v == channel {
			h.channels = append(h.channels[:k], h.channels[k+1:]...)
			return
		}
	}
}

func (h *IrcSubSystem) Run() {
	// channel signaling irc connection status
	h.ConnectedState = make(chan bool)

	// setup IRC client:
	h.client = irc.SimpleClient(h.ircnick, "ircflu", "ircflu")
	h.client.SSL = h.ircssl

	h.client.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		h.ConnectedState <- true
	})
	h.client.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		h.ConnectedState <- false
	})
	h.client.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		channel := line.Args[0]
		if channel == h.client.Me.Nick {
			log.Println("PM from " + line.Src)
			channel = line.Src // replies go via PM too.
		} else {
			log.Println("Message in channel " + line.Args[0] + " from " + line.Src)
		}

		msg := msgsystem.Message{
			To:     []string{channel},
			Msg:    line.Args[1],
			Source: line.Src,
			Authed: auth.IsAuthed(line.Src),
		}
		h.messagesIn <- msg
	})

	// loop on IRC dis/connected events
	go func() {
		for {
			log.Println("Connecting to IRC...")
			err := h.client.Connect(h.irchost, h.ircpassword)
			if err != nil {
				log.Println("Failed to connect to IRC")
				log.Println(err)
				continue
			}
			for {
				status := <-h.ConnectedState
				if status {
					log.Println("Connected to IRC")

					if len(h.channels) == 0 {
						// join default channel
						h.Join(h.ircchannel)
					} else {
						// we must have been disconnected, rejoin channels
						h.Rejoin()
					}
				} else {
					log.Println("Disconnected from IRC")
					break
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			cm := <-h.messagesOut
			fmt.Println("Sending:", cm.To, cm.Msg)
			if len(cm.To) == 0 {
				h.client.Privmsg(h.ircchannel, cm.Msg)
			} else {
				for _, recv := range cm.To {
					if recv == "#*" {
						// special: send to all joined channels
						for _, to := range h.channels {
							h.client.Privmsg(to, cm.Msg)
						}
					} else {
						// needs stripping hostname when sending to user!host
						if strings.Index(recv, "!") > 0 {
							recv = recv[0:strings.Index(recv, "!")]
						}

						h.client.Privmsg(recv, cm.Msg)
					}
				}
			}
		}
	}()
}

func init() {
	irc := IrcSubSystem{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&irc.irchost, "irchost", "localhost:6667", "Hostname of IRC server, eg: irc.example.org:6667"},
		app.CliFlag{&irc.ircnick, "ircnick", "ircflu", "Nickname to use for IRC"},
		app.CliFlag{&irc.ircpassword, "ircpassword", "", "Password to use to connect to IRC server"},
		app.CliFlag{&irc.ircchannel, "ircchannel", "#ircflutest", "Which channel to join"},
		app.CliFlag{&irc.ircssl, "ircssl", false, "Use SSL for IRC connection"},
	})

	msgsystem.RegisterSubSystem(&irc)
}
