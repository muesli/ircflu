// A Jira web-hook sending messages when new commits arrive.
package jira

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"github.com/muesli/ircflu/msgsystem"
	"github.com/muesli/ircflu/msgsystem/irc/irctools"
	"github.com/muesli/ircflu/msgsystem/web/hooks"
	//"strconv"
	"strings"
	"regexp"
)

var ()

type JiraHook struct {
	name     string
	path     string
	messages chan msgsystem.Message
}

func init() {
	hooks.RegisterWebHook(&JiraHook{name: "Jira", path: "/jira"})
}

func (hook *JiraHook) Name() string {
	return hook.name
}

func (hook *JiraHook) Path() string {
	return hook.path
}

func (hook *JiraHook) SetMessageChan(channel chan msgsystem.Message) {
	hook.messages = channel
}

func (hook *JiraHook) Request(ctx *web.Context) {
	decoder := json.NewDecoder(ctx.Request.Body)
	var payload interface{}
	err := decoder.Decode(&payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	data := payload.(map[string]interface{})
	issue := data["issue"].(map[string]interface{})
	issueFields := issue["fields"].(map[string]interface{})
	user := data["user"].(map[string]interface{})

	key := issue["key"].(string)
	summary := issueFields["summary"].(string)
	name := user["displayName"].(string)

	reg, err := regexp.Compile("rest/api.*")
	if err != nil {
		fmt.Println(err)
	}
	url := reg.ReplaceAllLiteralString(issue["self"].(string), "browse/" + key)
	action := ""
	event := data["webhookEvent"].(string)
	switch {
	case data["comment"] != nil:
		action = "Commented"
	case event == "jira:issue_created":
		action = "Created issue"
	case event == "jira:issue_updated":
		action = "Updated issue"
	default:
		action = event
	}

	msg := msgsystem.Message{
		Msg: fmt.Sprintf("[%s] %s %s %s %s", irctools.Colored(key, "lightblue"), summary, irctools.Colored(name, "teal"), action, url),
	}
	hook.messages <- msg

	if data["changelog"] != nil {
		changelog := data["changelog"].(map[string]interface{})
		changes := changelog["items"].([]interface{})
		for _, c := range changes {
			change := c.(map[string]interface{})
			field := change["field"].(string)
			var msg msgsystem.Message
			if change["toString"] == nil {
				deletedValue := change["fromString"].(string)
				msg = msgsystem.Message{
					Msg: fmt.Sprintf("Deleted %s %s", irctools.Colored(field, "lightblue"), irctools.Colored(deletedValue, "teal")),
				}
			} else {
				newStatus := change["toString"].(string)
				switch field {
				case "Attachment":
					action = "Added %s %s"
				default:
					action = "Changed %s to %s"
				}
				msg = msgsystem.Message{
					Msg: fmt.Sprintf(action, irctools.Colored(strings.Title(field), "lightblue"), irctools.Colored(newStatus, "teal")),
				}
			}
			hook.messages <- msg
		}
	}
}
