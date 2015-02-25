package main

import (
	_ "log"

	"github.com/pepl/ircflu/app"

	"github.com/pepl/ircflu/commands"
	_ "github.com/pepl/ircflu/commands/alias"
	_ "github.com/pepl/ircflu/commands/auth"
	_ "github.com/pepl/ircflu/commands/exec"
	_ "github.com/pepl/ircflu/commands/join"
	_ "github.com/pepl/ircflu/commands/part"
	_ "github.com/pepl/ircflu/commands/send"

	"github.com/pepl/ircflu/msgsystem"
	_ "github.com/pepl/ircflu/msgsystem/catserver"
	_ "github.com/pepl/ircflu/msgsystem/irc"
	//	_ "github.com/pepl/ircflu/msgsystem/jabber"
	_ "github.com/pepl/ircflu/msgsystem/web"
	_ "github.com/pepl/ircflu/msgsystem/web/hooks"
	_ "github.com/pepl/ircflu/msgsystem/web/hooks/github"
	_ "github.com/pepl/ircflu/msgsystem/web/hooks/gitlab"
	_ "github.com/pepl/ircflu/msgsystem/web/hooks/jira"
	_ "github.com/pepl/ircflu/msgsystem/web/hooks/fisheye"
)

func main() {
	// Parse command-line args for all registered sub modules
	app.Run()

	// Initialize commands and messaging sub-systems
	commands.StartCommands()
	msgsystem.StartSubSystems()

	// Keep app alive
	ch := make(chan bool)
	<-ch
}
