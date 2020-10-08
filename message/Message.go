package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	cookie "github.com/akshayaky/gICCa/cookie"
	"github.com/jroimartin/gocui"
)

//BufferMsg is an event received when a message is send by someone to the user
type BufferMsg struct {
	Cid      int    `json:"cid"`
	Bid      int    `json:"bid"`
	Type     string `json:"type"`
	Chan     string `json:"chan"`
	Eid      int    `json:"eid"`
	Msg      string `json:"msg"`
	From     string `json:"from"`
	StreamID string `json:"streamid"`
}

//Makeserver is a response when connecting to a network
type Makeserver struct {
	Cid          int    `json:"cid"`
	Type         string `json:"type"`
	Hostname     string `json:"hostname"`
	Port         int    `json:"port"`
	Ssl          bool   `json:"ssl"`
	Name         string `json:"name"`
	Nick         string `json:"nick"`
	Realname     string `json:"realname"`
	ServerPass   string `json:"server_pass"`
	NickservPass string `json:"nickserv_pass"`
	JoinCommands string `json:"join_commands"`
	//Ignores       string or bool `json:"ignores"`
	Away        string          `json:"away"`
	Status      string          `json:"status"`
	FailInfo    json.RawMessage `json:"fail_info"`
	Ircserver   string          `json:"ircserver"`
	IdentPrefix string          `json:"ident_prefix"`
	User        string          `json:"user"`
	Userhost    string          `json:"userhost"`
	Usermask    string          `json:"usermask"`
	NumBuffers  int             `json:"num_buffers"`
	Prefs       json.RawMessage `json:"prefs"`

	// /Disconnected string `json:"disconnected"`
}

//SendMessages sends messages to the specified nick or group
func SendMessages(session string, cid int, name string) func(string) {

	client := cookie.SetCookie(session, "say")
	Cid := strconv.Itoa(cid)
	var reciever string
	var hash string
	Send := func(msg string) {
		hash = ""
		if name[0] == '#' {
			hash = "%23"
			reciever = name[1:len(name)]
		} else {
			reciever = name
		}
		body := strings.NewReader(`msg=/msg%20` + hash + reciever + `%20` + msg + `&to=%2A&cid=` + Cid + `&session=` + session)
		req, err := http.NewRequest("POST", "https://www.irccloud.com/chat/say", body)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
	}
	return Send

}

/*
ViewMessages reads the stream of bytes
and displays the messages in the corresponding group or nick
*/
func ViewMessages(session string, reader *bufio.Reader, toName string, g *gocui.Gui, chan2 chan string) {

	var msg BufferMsg
	var last string
	last = ""

	var t time.Time
	for {

		line, _ := reader.ReadBytes('\n')

		if err := json.Unmarshal(line, &msg); err != nil {
			fmt.Println("I found the error here")
			panic(err)
		}

		if msg.Type == "buffer_msg" {
			chan2 <- msg.Msg
			g.Update(func(g *gocui.Gui) error {
				v, _ := g.View("mainView")
				if msg.Chan == toName {
					if msg.From != last {
					        fmt.Fprintln(v, fmt.Sprintf("\n\033[35;2m<%s>\033[0m", msg.From))
					        last = msg.From
					}
					t = time.Now()
					fmt.Fprintln(v, fmt.Sprintf("\033[39;2m(%s)\033[34;2m%s\033[0m", t.Format("15:04:05"), msg.Msg))
				return nil
			})

		}

	}
}
