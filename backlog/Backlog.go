package backlog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/akshayaky/gICCa/concurrent"
	cookie "github.com/akshayaky/gICCa/cookie"
	"github.com/akshayaky/gICCa/message"
	"github.com/jroimartin/gocui"
)

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

	bodyBytes := []message.Makeserver{}
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
func EndpointConnection(session string, g *gocui.Gui) ([10]int, [10]string, *bufio.Reader) {

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

	name, cid = GetNameAndCid(session, g, reader)

	return cid, name, reader
}

/*
GetNameAndCid returns the connnection id (CID)
and the corresponding name of the connection
*/
func GetNameAndCid(session string, g *gocui.Gui, reader *bufio.Reader) ([10]string, [10]int) {

	line, _ := reader.ReadBytes('\n')

	var msg message.BufferMsg

	if err := json.Unmarshal(line, &msg); err != nil {
		fmt.Println("I found the error here")
		panic(err)
	}

	name, cid := GetBacklog(msg.StreamID, session)

	//fmt.Println(cid, name)

	//var option int
	//fmt.Printf("Connection Index : ")
	//fmt.Scanln(&option)

	concurrent.Decider(session, g, reader, cid[3])

	return name, cid

}
