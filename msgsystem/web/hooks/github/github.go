// A GitHub web-hook sending messages when new commits arrive.
package github

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"github.com/pepl/ircflu/msgsystem"
	"github.com/pepl/ircflu/msgsystem/irc/irctools"
	"github.com/pepl/ircflu/msgsystem/web/hooks"
	"strconv"
	"strings"
)

var ()

type GitHubHook struct {
	name     string
	path     string
	messages chan msgsystem.Message
}

func init() {
	hooks.RegisterWebHook(&GitHubHook{name: "GitHub", path: "/github"})
}

func (hook *GitHubHook) Name() string {
	return hook.name
}

func (hook *GitHubHook) Path() string {
	return hook.path
}

func (hook *GitHubHook) SetMessageChan(channel chan msgsystem.Message) {
	hook.messages = channel
}

func (hook *GitHubHook) Request(ctx *web.Context) {
	payloadString, ok := ctx.Params["payload"]
	if !ok {
		fmt.Println("Couldn't find GitHub payload!")
		return
	}

	b := []byte(payloadString)

	var payload interface{}
	err := json.Unmarshal(b, &payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	data := payload.(map[string]interface{})

	before := data["before"].(string)
	after := data["after"].(string)
	ref := data["ref"].(string)
	ref = ref[strings.LastIndex(ref, "/")+1:]
	user := ""
	commitData := data["commits"].([]interface{})
	commitCount := 0

	repoData := data["repository"].(map[string]interface{})
	repo := repoData["name"].(string)
	url := repoData["url"].(string) + "/compare/" + before[:8] + "..." + after[:8]

	var ircmsgs []msgsystem.Message
	for _, c := range commitData {
		commit := c.(map[string]interface{})
		commitId := commit["id"].(string)
		if commitId == before {
			continue
		}

		if len(user) == 0 {
			author := commit["author"].(map[string]interface{})
			user = author["name"].(string)
		}

		commitCount++
		message := commit["message"].(string)

		msg := msgsystem.Message{
			Msg: fmt.Sprintf("%s/%s %s %s: %s", irctools.Colored(repo, "lightblue"), irctools.Colored(ref, "purple"), irctools.Colored(commitId[:8], "grey"), irctools.Colored(user, "teal"), message),
		}
		ircmsgs = append(ircmsgs, msg)
	}

	commitToken := "commits"
	if commitCount == 1 {
		commitToken = "commit"
	}
	msg := msgsystem.Message{
		Msg: fmt.Sprintf("[%s] %s pushed %s new %s to %s: %s", irctools.Colored(repo, "lightblue"), irctools.Colored(user, "teal"), irctools.Bold(strconv.Itoa(commitCount)), commitToken, irctools.Colored(ref, "purple"), url),
	}
	hook.messages <- msg

	for _, m := range ircmsgs {
		hook.messages <- m
	}
}
