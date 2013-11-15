package main

import (
	_ "log"

	"ircflu/app"

	"ircflu/msgsystem"
	_ "ircflu/msgsystem/catserver"
	_ "ircflu/msgsystem/irc"
	_ "ircflu/msgsystem/web"

	_ "ircflu/hooks"
	_ "ircflu/hooks/github"
	_ "ircflu/hooks/gitlab"
)

func main() {
	app.Run()
	msgsystem.StartSubSystems()

	ch := make(chan bool)
	<- ch
}
