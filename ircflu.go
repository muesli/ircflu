package main

import (
	_ "log"

	"ircflu/app"

	_ "ircflu/commands"
	_ "ircflu/commands/alias"
	_ "ircflu/commands/exec"
	_ "ircflu/commands/join"
	_ "ircflu/commands/send"

	"ircflu/msgsystem"
	_ "ircflu/msgsystem/catserver"
	_ "ircflu/msgsystem/irc"
	_ "ircflu/msgsystem/web"
	_ "ircflu/msgsystem/web/hooks"
	_ "ircflu/msgsystem/web/hooks/github"
	_ "ircflu/msgsystem/web/hooks/gitlab"
)

func main() {
	app.Run()
	msgsystem.StartSubSystems()

	ch := make(chan bool)
	<- ch
}
