// A FishEye web-hook sending messages when new commits arrive.
package fisheye

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"github.com/pepl/ircflu/msgsystem"
	"github.com/pepl/ircflu/msgsystem/irc/irctools"
	"github.com/pepl/ircflu/msgsystem/web/hooks"
	"io/ioutil"
)

var ()

type FishEyeHook struct {
	name     string
	path     string
	messages chan msgsystem.Message
}

func init() {
	hooks.RegisterWebHook(&FishEyeHook{name: "FishEye", path: "/fisheye"})
}

func (hook *FishEyeHook) Name() string {
	return hook.name
}

func (hook *FishEyeHook) Path() string {
	return hook.path
}

func (hook *FishEyeHook) SetMessageChan(channel chan msgsystem.Message) {
	hook.messages = channel
}

func (hook *FishEyeHook) Request(ctx *web.Context) {
	payload, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("ReadAll Error:", err)
		return
	}

	//fmt.Println("Payload:", string(payload))
	var payloadif interface{}
	err = json.Unmarshal(payload, &payloadif)
	if err != nil {
		fmt.Println("JSON Unmarshal Error:", err)
		return
	}

	data := payloadif.(map[string]interface{})
	//fmt.Println("%#v", data)

	repository := data["repository"].(map[string]interface{})
	repoName := repository["name"].(string)

	changeset := data["changeset"].(map[string]interface{})
	branches := changeset["branches"].([]interface{})
	// only one cs and thus only one branch at a time
	first_branch := branches[0].(string)

	msg := msgsystem.Message{
		Msg: fmt.Sprintf("[%s/%s] %s %s %s", irctools.Colored(repoName, "purple"), irctools.Colored(first_branch, "lightblue"), changeset["csid"].(string), irctools.Colored(changeset["author"].(string), "teal"), changeset["comment"].(string)),
	}
	hook.messages <- msg
}
