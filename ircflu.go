package main

import (
	"flag"
	"fmt"
	"log"
	irc "github.com/fluffle/goirc/client"
	"github.com/hoisie/web"
	"ircflu/catserver"
	"ircflu/hooks/github"
	"ircflu/hooks/gitlab"
	"strings"
	"time"
)

var (
	// to read parser output:
	messages = make(chan string)

	irchost     string
	ircnick     string
	ircpassword string
	ircssl      bool
	ircchannel  string

	catbind string
	catfam  string
)

func init() {
	flag.StringVar(&irchost, "irchost", "localhost:6667",
		"Hostname of IRC server, eg: irc.example.org:6667")
	flag.StringVar(&ircnick, "ircnick", "ircflu",
		"Nickname to use for IRC")
	flag.StringVar(&ircpassword, "ircpassword", "",
		"Password to use to connect to IRC server")
	flag.BoolVar(&ircssl, "ircssl", false,
		"Use SSL for IRC connection")
	flag.StringVar(&ircchannel, "ircchannel", "#ircflutest",
		"Which channel to join")
	flag.StringVar(&catbind, "catbind", ":12345",
		"net.Listen spec, to listen for IRCCat msgs")
	flag.StringVar(&catfam, "catfamily", "tcp4",
		"net.Listen address family for IRCCat msgs")

	flag.Parse()

	github.Messages = messages
	gitlab.Messages = messages
}

func CatMsgSender(ch chan catserver.CatMsg, client *irc.Conn) {
	for {
		cm := <-ch
		if len(cm.To) == 0 {
			client.Privmsg(ircchannel, cm.Msg)
		} else {
			for _, to := range cm.To {
				client.Privmsg(to, cm.Msg)
			}
		}
	}
}

func setupClient(c *irc.Conn, chConnected chan bool) {
	c.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		chConnected <- true
	})
	c.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		chConnected <- false
	})
	c.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		if len(line.Args) < 2 || !strings.HasPrefix(line.Args[1], "!") {
			return
		}
		to := line.Args[0]
		sender := to

		if to == c.Me.Nick {
			// TODO: check if sender is in main chan, else return
			log.Println("Got via PM from " + line.Src)
			sender = line.Src // replies go via PM too.
		} else {
			log.Println("Got via channel " + line.Args[0] + " from " + line.Src)
		}

		msg := strings.Split(line.Args[1], " ")
		cmd := msg[0]
		params := strings.Join(msg[1:], " ")
		fmt.Println("Command: " + cmd + " with Params: " + params)

		switch cmd {
			case "!join":
				if len(line.Args) == 2 {
					c.Join(params)
				} else {
					c.Privmsg(sender, "Usage: !join #chan  or  !join #chan key")
				}
			case "!part":
				if len(line.Args) == 2 {
					c.Part(params)
				} else {
					c.Privmsg(sender, "Usage: !part #chan")
				}
			default:
				c.Privmsg(sender, "Invalid command: " + cmd)
				return
		}
	})
}

func decodeParam(param string) string {
	return strings.Replace(param, "+", " ", -1)
}


func main() {
	// msgs from tcp catport to this channel
	catmsgs := make(chan catserver.CatMsg)
	// channel signaling irc connection status
	chConnected := make(chan bool)

	// setup IRC client:
	c := irc.SimpleClient(ircnick)
	c.SSL = ircssl

	// Listen on catport:
	go catserver.CatportServer(catmsgs, catfam, catbind)
	go CatMsgSender(catmsgs, c)

	// loop on IRC dis/connected events
	setupClient(c, chConnected)
	go func() {
		for {
			log.Println("Connecting to IRC...")
			err := c.Connect(irchost, ircpassword)
			if err != nil {
				log.Println("Failed to connect to IRC")
				log.Println(err)
				continue
			}
			for {
				status := <-chConnected
				if status {
					log.Println("Connected to IRC")
					c.Join(ircchannel)
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
			m := <-messages
			c.Privmsg(ircchannel, m)
		}
	}()

	web.Post("/github", github.GitHubHook)
	web.Post("/gitlab", gitlab.GitLabHook)
	web.Run("0.0.0.0:12346")
}
