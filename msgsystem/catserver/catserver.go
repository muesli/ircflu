// Based on gocat's catserver by Richard Jones - https://github.com/RJ/gocat
// Listens on a TCP port, parses first line for addressees,
// Puts Message onto the out channel.
package catserver

import (
	"bufio"
	"github.com/muesli/ircflu/app"
	"github.com/muesli/ircflu/msgsystem"
	"io"
	"log"
	"net"
	"strings"
)

type CatSubSystem struct {
	name        string
	messagesIn  chan msgsystem.Message
	messagesOut chan msgsystem.Message

	catbind string
	catfam  string
}

func (h *CatSubSystem) Name() string {
	return h.name
}

func (h *CatSubSystem) MessageInChan() chan msgsystem.Message {
	return h.messagesIn
}

func (h *CatSubSystem) SetMessageInChan(channel chan msgsystem.Message) {
	h.messagesIn = channel
}

func (h *CatSubSystem) MessageOutChan() chan msgsystem.Message {
	return h.messagesOut
}

func (h *CatSubSystem) SetMessageOutChan(channel chan msgsystem.Message) {
	h.messagesOut = channel
}

func CatportServer(catmsgs chan msgsystem.Message, catfamily string, catbind string) {
	netListen, error := net.Listen(catfamily, catbind)
	if error != nil {
		log.Fatal(error)
	} else {
		defer netListen.Close()
		for {
			log.Println("Waiting for clients")
			connection, error := netListen.Accept()
			if error == io.EOF {
				break
			}
			if error != nil {
				log.Println("Client error: ", error)
			} else {
				go ClientHandler(connection, catmsgs)
			}
		}
	}
}

func ClientHandler(connection net.Conn, catmsgs chan msgsystem.Message) {
	log.Println(connection.RemoteAddr().String() + " connected.")
	reader := bufio.NewReader(connection)
	seenFirst := false
	to := []string{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(connection.RemoteAddr().String() + " error reading line")
			log.Println(err)
			break
		} else {
			line = strings.TrimRight(line, "\r\n")
			log.Println(connection.RemoteAddr().String() + " -> " + line)
			// ensure we have captured the to-address
			if !seenFirst {
				seenFirst = true
				// replace line (potentially)
				newto, line2 := ParseFirstLine(line)
				line = line2
				to = newto
			}
			cm := msgsystem.Message{
				To:  to,
				Msg: "!send " + line,
			}
			catmsgs <- cm
		}
	}
	log.Println(connection.RemoteAddr().String() + " closed.")
}

func ParseFirstLine(str string) ([]string, string) {
	strparts := strings.SplitN(str, " ", 2)
	if len(strparts) == 1 {
		return []string{}, str
	}
	firstword := strparts[0]
	rest := strparts[1]

	// special spam mode, all joined channels:
	if firstword == "#*" {
		return []string{firstword}, rest
	}
	// maybe nothing specified, we end up using the default channel from config
	if strings.Index(firstword, "#") != 0 &&
		strings.Index(firstword, "@") != 0 {
		return []string{}, str
	}
	parts := strings.Split(firstword, ",")
	for i, p := range parts {
		// User nicks start with @, which needs stripping:
		if strings.Index(p, "@") == 0 {
			parts[i] = strings.TrimLeft(p, "@")
		}
		// otherwise considered to be a channel name, leave as-is
	}
	return parts, rest
}

func (h *CatSubSystem) Run() {
	// Listen on catport:
	go CatportServer(h.messagesIn, h.catfam, h.catbind)
}

func init() {
	cat := CatSubSystem{name: "cat"}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&cat.catbind, "catbind", ":12345", "net.Listen spec, to listen for IRCCat msgs"},
		app.CliFlag{&cat.catfam, "catfamily", "tcp4", "net.Listen address family for IRCCat msgs"},
	})

	msgsystem.RegisterSubSystem(&cat)
}
