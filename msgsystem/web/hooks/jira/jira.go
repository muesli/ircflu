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
	"io/ioutil"
)

var ()

type JiraHook struct {
	name     string
	path     string
	messages chan msgsystem.Message
}

type JiraWebhook struct {
	Id int `json:",string"`
	Issue JiraIssue
	User JiraUser
	Changelog *JiraChangelog
	Comment *JiraComment
	WebhookEvent string
}

type JiraIssue struct {
	Id int `json:",string"`
	Self string
	Key string
	Fields JiraIssueFields
}

type JiraIssueFields struct {
	Summary string
	Description string
	Labels []string
	Priority JiraPriority
}

type JiraPriority struct {
	Id int `json:",string"`
	Self string
	Name string
}

type JiraUser struct {
	Self string
	Name string
	EmailAddress string
	DisplayName string
	Active bool
}

type JiraChangelog struct {
	Id int `json:",string"`
	Items []JiraChangelogItem
}

type JiraChangelogItem struct {
	ToString string
	To string
	FromString string
	From string
	Fieldtype string
	Field string
}

type JiraComment struct {
	Self string
	Id int `json:",string"`
	Author JiraUser
	Body string
	UpdateAuthor JiraUser
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
	payload, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var data JiraWebhook
	err = json.Unmarshal(payload, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	reg, err := regexp.Compile("rest/api.*")
	if err != nil {
		fmt.Println(err)
	}
	url := reg.ReplaceAllLiteralString(data.Issue.Self, "browse/" + data.Issue.Key)
	action := ""
	event := data.WebhookEvent
	switch {
	case data.Comment != nil:
		action = "Commented"
	case event == "jira:issue_created":
		action = "Created issue"
	case event == "jira:issue_updated":
		action = "Updated issue"
	default:
		action = event
	}

	msg := msgsystem.Message{
		Msg: fmt.Sprintf("[%s] %s %s %s %s", irctools.Colored(data.Issue.Key, "lightblue"), data.Issue.Fields.Summary, irctools.Colored(data.User.DisplayName, "teal"), action, url),
	}
	hook.messages <- msg

	if data.Changelog != nil {
		for _, change := range data.Changelog.Items {
			field := change.Field
			var msg msgsystem.Message
			if change.ToString == "" {
				msg = msgsystem.Message{
					Msg: fmt.Sprintf("Deleted %s %s", irctools.Colored(field, "lightblue"), irctools.Colored(change.FromString, "teal")),
				}
			} else {
				switch field {
				case "Attachment":
					msg = msgsystem.Message{
						Msg: fmt.Sprintf("Added %s %s", irctools.Colored(strings.Title(field), "lightblue"), irctools.Colored(change.ToString, "teal")),
					}
				default:
					msg = msgsystem.Message{
						Msg: fmt.Sprintf("Changed %s from %s to %s", irctools.Colored(strings.Title(field), "lightblue"), irctools.Colored(change.FromString, "purple"), irctools.Colored(change.ToString, "teal")),
					}
				}

			}
			hook.messages <- msg
		}
	}
}
