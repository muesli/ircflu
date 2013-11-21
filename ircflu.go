package main

import (
	_ "log"

	"github.com/muesli/ircflu/app"

	"github.com/muesli/ircflu/commands"
	_ "github.com/muesli/ircflu/commands/alias"
	_ "github.com/muesli/ircflu/commands/auth"
	_ "github.com/muesli/ircflu/commands/exec"
	_ "github.com/muesli/ircflu/commands/join"
	_ "github.com/muesli/ircflu/commands/part"
	_ "github.com/muesli/ircflu/commands/send"

	"github.com/muesli/ircflu/msgsystem"
	_ "github.com/muesli/ircflu/msgsystem/catserver"
	_ "github.com/muesli/ircflu/msgsystem/irc"
	_ "github.com/muesli/ircflu/msgsystem/web"
	_ "github.com/muesli/ircflu/msgsystem/web/hooks"
	_ "github.com/muesli/ircflu/msgsystem/web/hooks/github"
	_ "github.com/muesli/ircflu/msgsystem/web/hooks/gitlab"
)

func main() {
	app.Run()
	commands.StartCommands()
	msgsystem.StartSubSystems()

	ch := make(chan bool)
	<-ch
}
