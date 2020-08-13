package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	cookie "github.com/akshayaky/gICCa/cookie"
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
func SendMessages(session string, cid int, name string, chan1 chan string) {
	client := cookie.SetCookie(session, "say")
	reader := bufio.NewReader(os.Stdin)
	Cid := strconv.Itoa(cid)
	var msg string
	var reciever string
	hash := ""
	for {
		msg, _ = reader.ReadString('\n')

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
		chan1 <- msg
	}

}

/*
ViewMessages reads the stream of bytes
and displays the messages in the corresponding group or nick
*/
func ViewMessages(session string, reader *bufio.Reader, toName string, chan2 chan string) {

	// only viewing messages in a chat
	for {

		line, _ := reader.ReadBytes('\n')

		var msg BufferMsg

		if err := json.Unmarshal(line, &msg); err != nil {
			fmt.Println("I found the error here")
			panic(err)
		}

		if msg.Type == "buffer_msg" {
			if msg.From == toName {
				fmt.Printf(toName + " : ")
				fmt.Println(msg.Msg)

			} else if msg.Chan == toName {
				fmt.Printf(msg.From + " : ")
				fmt.Println(msg.Msg)
				chan2 <- msg.Msg
			}
		}

	}
}
