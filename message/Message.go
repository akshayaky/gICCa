package message

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	cookie "github.com/akshayaky/gICCa/cookie"
)

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
