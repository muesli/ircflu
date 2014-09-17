// A GitLab web-hook sending messages when new commits arrive.
package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"github.com/pepl/ircflu/msgsystem"
	"github.com/pepl/ircflu/msgsystem/irc/irctools"
	"github.com/pepl/ircflu/msgsystem/web/hooks"
	"io/ioutil"
	"strconv"
	"strings"
)

var ()

type GitLabHook struct {
	name     string
	path     string
	messages chan msgsystem.Message
}

func init() {
	hooks.RegisterWebHook(&GitLabHook{name: "GitLab", path: "/gitlab"})
}

func (hook *GitLabHook) Name() string {
	return hook.name
}

func (hook *GitLabHook) Path() string {
	return hook.path
}

func (hook *GitLabHook) SetMessageChan(channel chan msgsystem.Message) {
	hook.messages = channel
}

func (hook *GitLabHook) Request(ctx *web.Context) {
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Params:", string(b))

	var payload interface{}
	err = json.Unmarshal(b, &payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	data := payload.(map[string]interface{})

	before := data["before"].(string)
	after := data["after"].(string)
	ref := data["ref"].(string)
	ref = ref[strings.LastIndex(ref, "/")+1:]
	user := data["user_name"].(string)
	commitData := data["commits"].([]interface{})
	commitCount := int(data["total_commits_count"].(float64))

	repoData := data["repository"].(map[string]interface{})
	repo := repoData["name"].(string)
	url := repoData["homepage"].(string) + "/compare/" + before[:8] + "..." + after[:8]

	commitToken := "commits"
	if commitCount == 1 {
		commitToken = "commit"
	}
	ircmsg := msgsystem.Message{
		Msg: fmt.Sprintf("[%s] %s pushed %s new %s to %s: %s", irctools.Colored(repo, "lightblue"), irctools.Colored(user, "teal"), irctools.Bold(strconv.Itoa(commitCount)), commitToken, irctools.Colored(ref, "purple"), url),
	}
	hook.messages <- ircmsg

	for _, c := range commitData {
		commit := c.(map[string]interface{})
		commitId := commit["id"].(string)
		if commitId == before {
			continue
		}

		message := commit["message"].(string)
		ircmsg := msgsystem.Message{
			Msg: fmt.Sprintf("%s/%s %s %s: %s", irctools.Colored(repo, "lightblue"), irctools.Colored(ref, "purple"), irctools.Colored(commitId[:8], "grey"), irctools.Colored(user, "teal"), message),
		}
		hook.messages <- ircmsg
	}
}
