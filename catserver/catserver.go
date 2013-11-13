// Imported from gocat by Richard Jones - https://github.com/RJ/gocat
// Listens on a TCP port, parses first line for addressees,
// puts CatMsg onto channel.
package catserver

import "net"
import "bufio"
import "log"
import "io"
import "strings"

type CatMsg struct {
    To  []string
    Msg string
}

func CatportServer(catmsgs chan CatMsg, catfamily string, catbind string) {
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

func ClientHandler(connection net.Conn, catmsgs chan CatMsg) {
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
            cm := CatMsg{to, line}
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
