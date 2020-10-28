package backlog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/akshayaky/gICCa/concurrent"
	cookie "github.com/akshayaky/gICCa/cookie"
	"github.com/akshayaky/gICCa/message"
	"github.com/jroimartin/gocui"
)

//EndpointConnection connects to the stream
func EndpointConnection(session string, g *gocui.Gui) {

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

	channels, servers := GetAllData(session, g, reader)
	concurrent.Decider(session, g, reader, servers, channels)
}

//GetBacklog gets the backlog and assigns them to appropriate structs
func GetBacklog(streamid string, session string) io.ReadCloser {
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
	return resp.Body
}

/*
GetAllData returns the connnection id (CID)
and the corresponding name of the connection
*/
func GetAllData(session string, g *gocui.Gui, reader *bufio.Reader) ([10]message.ChannelInit, [10]message.Makeserver) {

	line, _ := reader.ReadBytes('\n')

	var msg message.BufferMsg

	if err := json.Unmarshal(line, &msg); err != nil {
		fmt.Println("I found the error here")
		panic(err)
	}

	var servers [10]message.Makeserver
	var channels [10]message.ChannelInit

	channels, servers = processBacklog(msg.StreamID, session)
	return channels, servers

}

func processBacklog(streamID string, session string) ([10]message.ChannelInit, [10]message.Makeserver) {

	var channels [10]message.ChannelInit
	var servers [10]message.Makeserver

	body := GetBacklog(streamID, session)
	bodyBytes := []message.Backlog{}
	err := json.NewDecoder(body).Decode(&bodyBytes)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	i := 0
	j := 0
	for _, event := range bodyBytes {

		if event.Type == "channel_init" {

			channels[i].Cid = event.Cid
			channels[i].Chan = event.Chan
			channels[i].Members = event.Members
			channels[i].Topics = event.Topics
			i++
		} else if event.Type == "makeserver" {
			servers[j].Cid = event.Cid
			servers[j].Name = event.Name
			j++
		}
	}
	fmt.Println(channels[0].Chan)
	return channels, servers
}
