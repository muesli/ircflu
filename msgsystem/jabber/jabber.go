// ircflu's Jabber subsystem.
package jabber

import (
	"fmt"
	"github.com/mattn/go-xmpp"
	"github.com/muesli/ircflu/app"
	"github.com/muesli/ircflu/auth"
	"github.com/muesli/ircflu/msgsystem"
	"log"
)

type JabberSubSystem struct {
	client *xmpp.Client

	server string
	username string
	password string
	notls bool
}

func (sys *JabberSubSystem) Name() string {
	return "jabber"
}

func (sys *JabberSubSystem) Run(channelIn, channelOut chan msgsystem.Message) {
	if len(sys.server) == 0 {
		return
	}

	var talk *xmpp.Client
	var err error
	if sys.notls {
		talk, err = xmpp.NewClientNoTLS(sys.server, sys.username, sys.password, false)
	} else {
		talk, err = xmpp.NewClient(sys.server, sys.username, sys.password, false)
	}
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			chat, err := talk.Recv()
			if err != nil {
				log.Fatal(err)
			}
			switch v := chat.(type) {
				case xmpp.Chat:
					fmt.Println(v.Remote, v.Text)

					msg := msgsystem.Message{
						To:     []string{sys.username},
						Msg:    v.Text,
						Source: v.Remote,
						Authed: auth.IsAuthed(v.Remote),
					}
					channelIn <- msg
				case xmpp.Presence:
//					fmt.Println(v.From, v.Show)
			}
		}
	}()

	go func() {
		for {
			cm := <-channelOut
			fmt.Println("Sending:", cm.To, cm.Msg)
		}
	}()
}

func init() {
	jabber := JabberSubSystem{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&jabber.server, "jabberhost", "localhost:443", "Hostname of Jabber server, eg: talk.google.com:443"},
		app.CliFlag{&jabber.username, "jabberuser", "ircflu", "Username to authenticate with Jabber server"},
		app.CliFlag{&jabber.password, "jabberpassword", "", "Password to use to connect to Jabber server"},
		app.CliFlag{&jabber.notls, "jabbernotls", false, "If you don't want to connect with TLS"},
	})

	msgsystem.RegisterSubSystem(&jabber)
}
