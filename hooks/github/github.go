package github

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"io/ioutil"
	"ircflu/irctools"
	"strconv"
	"strings"
)

var (
	Messages chan string
)

func GitHubHook(ctx *web.Context) {
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
	user := ""
	commitData := data["commits"].([]interface{})
	commitCount := 0

	repoData := data["repository"].(map[string]interface{})
	repo := repoData["name"].(string)
	url := repoData["homepage"].(string) + "/compare/" + before[:8] + "..." + after[:8]

	commitToken := "commits"
	if commitCount == 1 {
		commitToken = "commit"
	}

	var ircmsgs []string
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
		msg := fmt.Sprintf("%s/%s %s %s: %s", irctools.Colored(repo, "lightblue"), irctools.Colored(ref, "purple"), irctools.Colored(commitId[:8], "grey"), irctools.Colored(user, "teal"), message)
		ircmsgs = append(ircmsgs, msg)
	}

	msg := fmt.Sprintf("[%s] %s pushed %s new %s to %s: %s", irctools.Colored(repo, "lightblue"), irctools.Colored(user, "teal"), irctools.Bold(strconv.Itoa(commitCount)), commitToken, irctools.Colored(ref, "purple"), url)
	Messages <- msg

	for _, m := range ircmsgs {
		Messages <- m
	}
}
