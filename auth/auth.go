package auth

import (
	"fmt"
	"github.com/muesli/ircflu/app"
)

var (
	auths map[string]bool = make(map[string]bool)

	password string
)

func init() {
	fmt.Println("Initializing auth subsystem...")

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&password, "authpassword", "", "Password required to authenticate"},
	})
}

func Auth(source string, passwordIn string) bool {
	if len(password) == 0 {
		fmt.Println("Auth system disabled!", source)
		return false
	}
	if password == passwordIn {
		fmt.Println("Registering authed user:", source)
		auths[source] = true
		return true
	} else {
		fmt.Println("Auth'ing user failed:", source)
		auths[source] = false
		return false
	}
}

func IsAuthed(source string) bool {
	v, ok := auths[source]
	if ok {
		return v
	}

	return false
}
