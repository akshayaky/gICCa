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

	GetAllData(session, g, reader)
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
		// fmt.Println(event)

		if event.Type == "channel_init" {

			channels[i].Cid = event.Cid
			channels[i].Chan = event.Chan
			channels[i].Members = event.Members
			channels[i].Topics = event.Topics
			// fmt.Println(channels[i].Chan + strconv.Itoa(channels[i].Cid))
			i++
		} else if event.Type == "makeserver" {
			// fmt.Println(event.Cid, event.Name)
			servers[j].Cid = event.Cid
			servers[j].Name = event.Name
			j++
		}
	}
	// os.Exit(6)
	fmt.Println(channels[0].Chan)
	return channels, servers
}

/*
GetAllData returns the connnection id (CID)
and the corresponding name of the connection
*/
func GetAllData(session string, g *gocui.Gui, reader *bufio.Reader) {

	line, _ := reader.ReadBytes('\n')

	var msg message.BufferMsg

	if err := json.Unmarshal(line, &msg); err != nil {
		fmt.Println("I found the error here")
		panic(err)
	}

	// fmt.Printf("lknasdvfonasdfjinvfjnfvlglkjre")
	// name, cid, channels := GetBacklog(msg.StreamID, session)

	var servers [10]message.Makeserver
	var channels [10]message.ChannelInit

	channels, servers = processBacklog(msg.StreamID, session)

	g.Update(func(g *gocui.Gui) error {
		v, _ := g.View("channels")
		for _, channel := range channels {
			fmt.Fprintln(v, fmt.Sprintf("%s", channel.Chan)) //+strconv.Itoa(channel.Cid)))
		}
		return nil
	})

	concurrent.Decider(session, g, reader, servers)

}
