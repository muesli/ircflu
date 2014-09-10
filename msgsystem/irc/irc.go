// ircflu's IRC client subsystem.
package irc

import (
	_ "fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/muesli/ircflu/app"
	"github.com/muesli/ircflu/auth"
	"github.com/muesli/ircflu/msgsystem"
	"log"
	"strings"
	"time"
)

type IrcSubSystem struct {
	// channel signaling irc connection status
	ConnectedState chan bool

	// setup IRC client:
	client *irc.Conn

	irchost	 string
	ircnick	 string
	ircpassword string
	ircssl	  bool
	ircchannel  string

	channels	[]string
}

func (sys *IrcSubSystem) Name() string {
	return "irc"
}

func (sys *IrcSubSystem) Rejoin() {
	for _, channel := range sys.channels {
		sys.client.Join(channel)
	}
}

func (sys *IrcSubSystem) Join(channel string) {
	channel = strings.TrimSpace(channel)
	sys.client.Join(channel)

	sys.channels = append(sys.channels, channel)
}

func (sys *IrcSubSystem) Part(channel string) {
	channel = strings.TrimSpace(channel)
	sys.client.Part(channel)

	for k, v := range sys.channels {
		if v == channel {
			sys.channels = append(sys.channels[:k], sys.channels[k+1:]...)
			return
		}
	}
}

func (sys *IrcSubSystem) Run(channelIn, channelOut chan msgsystem.Message) {
	if len(sys.irchost) == 0 {
		return
	}

	// channel signaling irc connection status
	sys.ConnectedState = make(chan bool)

	// setup IRC client:
	sys.client = irc.SimpleClient(sys.ircnick, "ircflu", "ircflu")
	sys.client.Config().SSL = sys.ircssl

	sys.client.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) { sys.ConnectedState <- true })

	sys.client.HandleFunc("disconnected",
		func(conn *irc.Conn, line *irc.Line) { sys.ConnectedState <- false })

	sys.client.HandleFunc("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		channel := line.Args[0]
		text := ""
		if len(line.Args) > 1 {
			text = line.Args[1]
		}
		if channel == sys.ircnick {
			log.Println("PM from " + line.Src)
			channel = line.Src // replies go via PM too.
		} else {
			log.Println("Message in channel " + line.Args[0] + " from " + line.Src)
		}

		msg := msgsystem.Message{
			To:	 []string{channel},
			Msg:	text,
			Source: line.Src,
			Authed: auth.IsAuthed(line.Src),
		}
		channelIn <- msg
	})

	// loop on IRC dis/connected events
	go func() {
		for {
			log.Println("Connecting to IRC:", sys.irchost)
			sys.client.Config().Server = sys.irchost
			sys.client.Config().Pass = sys.ircpassword
			err := sys.client.Connect()
			if err != nil {
				log.Println("Failed to connect to IRC:", sys.irchost)
				log.Println(err)
				continue
			}
			for {
				status := <-sys.ConnectedState
				if status {
					log.Println("Connected to IRC:", sys.irchost)

					if len(sys.channels) == 0 {
						// join default channel
						sys.Join(sys.ircchannel)
					} else {
						// we must have been disconnected, rejoin channels
						sys.Rejoin()
					}
				} else {
					log.Println("Disconnected from IRC:", sys.irchost)
					break
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func (sys *IrcSubSystem) Handle(cm msgsystem.Message) bool {
	if len(cm.To) == 0 {
		sys.client.Privmsg(sys.ircchannel, cm.Msg)
		return true
	} else {
		for _, recv := range cm.To {
			if recv == "*" {
				// special: send to all joined channels
				for _, to := range sys.channels {
					sys.client.Privmsg(to, cm.Msg)
				}
			} else {
				// needs stripping hostname when sending to user!host
				if strings.Index(recv, "!") > 0 {
					recv = recv[0:strings.Index(recv, "!")]
				}

				sys.client.Privmsg(recv, cm.Msg)
			}
		}

		return true
	}

	return false
}

func init() {
	irc := IrcSubSystem{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&irc.irchost, "irchost", "", "Hostname of IRC server, eg: irc.example.org:6667"},
		app.CliFlag{&irc.ircnick, "ircnick", "ircflu", "Nickname to use for IRC"},
		app.CliFlag{&irc.ircpassword, "ircpassword", "", "Password to use to connect to IRC server"},
		app.CliFlag{&irc.ircchannel, "ircchannel", "#ircflutest", "Which channel to join"},
		app.CliFlag{&irc.ircssl, "ircssl", false, "Use SSL for IRC connection"},
	})

	msgsystem.RegisterSubSystem(&irc)
}
