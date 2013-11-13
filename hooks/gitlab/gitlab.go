package gitlab

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

func GitLabHook(ctx *web.Context ) {
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
	ircmsg := fmt.Sprintf("[%s] %s pushed %s new %s to %s: %s", irctools.Colored(repo, "lightblue"), irctools.Colored(user, "teal"), irctools.Bold(strconv.Itoa(commitCount)), commitToken, irctools.Colored(ref, "purple"), url)
	//	irccat.PostToIrc(irccataddr, ircmsg)
	Messages <- ircmsg

	for _, c := range commitData {
		commit := c.(map[string]interface{})
		commitId := commit["id"].(string)
		if commitId == before {
			continue
		}

		message := commit["message"].(string)
		ircmsg = fmt.Sprintf("%s/%s %s %s: %s", irctools.Colored(repo, "lightblue"), irctools.Colored(ref, "purple"), irctools.Colored(commitId[:8], "grey"), irctools.Colored(user, "teal"), message)
		//		irccat.PostToIrc(irccataddr, ircmsg)
		Messages <- ircmsg
	}
}
