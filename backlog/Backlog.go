package backlog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	cookie "github.com/akshayaky/gICCa/cookie"
	"github.com/akshayaky/gICCa/message"
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

//GetBacklog gets the backlog and assigns them to appropriate structs
func GetBacklog(streamid string, session string) ([10]string, [10]int) {
	client := cookie.SetCookie(session, "backlog/"+streamid)

	urlData := url.Values{}
	urlData.Set("session", session)

	req, err := http.NewRequest("GET", "http://www.irccloud.com/chat/backlog/"+streamid, strings.NewReader(urlData.Encode()))
	if err != nil {
		fmt.Println("error : ", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	bodyBytes := []Makeserver{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&bodyBytes)
	if err != nil {
		log.Fatal(err)
	}
	var name [10]string
	var cid [10]int

	i := 0

	for _, event := range bodyBytes {
		if event.Type == "makeserver" {

			cid[i] = event.Cid
			name[i] = event.Name
			i++
		}
		if i == 4 {
			break
		}
	}

	return name, cid

}

//EndpointConnection connects to the stream
func EndpointConnection(session string) ([10]int, [10]string, *bufio.Reader) {

	client := cookie.SetCookie(session, "stream")

	urlData := url.Values{}
	urlData.Set("session", session)

	req, err := http.NewRequest("POST", "https://www.irccloud.com/chat/stream",
		strings.NewReader(urlData.Encode()))
	if err != nil {
		fmt.Println("error sending the request :  ", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error is here")
		fmt.Println(err)
	}
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	var cid [10]int
	var name [10]string

	name, cid = GetNameAndCid(session, reader)

	return cid, name, reader
}

/*
GetNameAndCid returns the connnection id (CID)
and the corresponding name of the connection
*/
func GetNameAndCid(session string, reader *bufio.Reader) ([10]string, [10]int) {

	line, _ := reader.ReadBytes('\n')

	var msg BufferMsg

	if err := json.Unmarshal(line, &msg); err != nil {
		fmt.Println("I found the error here")
		panic(err)
	}

	name, cid := GetBacklog(msg.StreamID, session)

	fmt.Println(cid, name)

	var option int
	fmt.Printf("Connection Index : ")
	fmt.Scanln(&option)

	Decider(session, reader, cid[option])

	return name, cid

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

/*
Decider function uses concurrency to decided which
function to execute, viewMessages or Say1.
*/
func Decider(session string, reader *bufio.Reader, cid int) {

	chan1 := make(chan string)
	chan2 := make(chan string)

	var toName string

	fmt.Printf("\n\nEnter the Name : ")
	fmt.Scanln(&toName)
	fmt.Println(toName)

	//these two functions will run concurrently
	go message.SendMessages(session, cid, toName, chan1)
	go ViewMessages(session, reader, toName, chan2)

	for {
		select {
		case <-chan1:

		case <-chan2:

		}
	}
}
