package irc

import (
	irc "github.com/fluffle/goirc/client"
	"fmt"
	"log"
	_ "strings"
	"time"
	"ircflu/app"
	"ircflu/msgsystem"
)

type IrcSubSystem struct {
	name string
	messagesIn chan msgsystem.Message
	messagesOut chan msgsystem.Message

	// channel signaling irc connection status
	chConnected chan bool

	// setup IRC client:
	client *irc.Conn

	irchost     string
	ircnick     string
	ircpassword string
	ircssl      bool
	ircchannel  string
}

func (h *IrcSubSystem) Name() string {
	return h.name
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

func (h *IrcSubSystem) Run() {
	// channel signaling irc connection status
	h.chConnected = make(chan bool)

	// setup IRC client:
	h.client = irc.SimpleClient(h.ircnick, "ircflu", "ircflu")
	h.client.SSL = h.ircssl

	h.client.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		h.chConnected <- true
	})
	h.client.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		h.chConnected <- false
	})
	h.client.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		channel := line.Args[0]
		if channel == h.client.Me.Nick {
			// TODO: check if source is in main chan, else return
			log.Println("Got via PM from " + line.Src)
			channel = line.Src // replies go via PM too.
		} else {
			log.Println("Got via channel " + line.Args[0] + " from " + line.Src)
		}

		msg := msgsystem.Message{
			To: []string{channel},
			Msg: line.Args[1],
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
				status := <-h.chConnected
				if status {
					log.Println("Connected to IRC")
					h.client.Join(h.ircchannel)
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
				for _, to := range cm.To {
					h.client.Privmsg(to, cm.Msg)
				}
			}
		}
	}()
}

func init() {
	irc := IrcSubSystem{name: "irc"}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&irc.irchost, "irchost", "localhost:6667", "Hostname of IRC server, eg: irc.example.org:6667"},
		app.CliFlag{&irc.ircnick, "ircnick", "ircflu", "Nickname to use for IRC"},
		app.CliFlag{&irc.ircpassword, "ircpassword", "", "Password to use to connect to IRC server"},
		app.CliFlag{&irc.ircchannel, "ircchannel", "#ircflutest", "Which channel to join"},
	//	app.CliFlag{&irc.ircssl, "ircssl", false, "Use SSL for IRC connection"},
	})

	msgsystem.RegisterSubSystem(&irc)
}
