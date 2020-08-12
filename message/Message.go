package message

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	cookie "github.com/akshayaky/gICCa/cookie"
)

func Say(session string, cid string) {
	client := cookie.SetCookie(session)
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Printf("name : ")
		name, _ := reader.ReadString('\n')
		hash := ""
		if name[0:len(name)-1] == "/end" {
			fmt.Println("end")
			break
		}
		for true {
			fmt.Printf(name[0:len(name)-1] + " : ")
			msg, _ := reader.ReadString('\n')
			if msg[0:len(msg)-1] == "/end" {
				fmt.Println("end")
				break

			}
			reciever := name
			if name[0] == '#' {
				hash = "%23"
				reciever = name[1:len(name)]
			}

			body := strings.NewReader(`msg=/msg%20` + hash + reciever[0:len(reciever)-1] + `%20` + msg + `&to=%2A&cid=` + cid + `&session=` + session)
			req, err := http.NewRequest("POST", "https://www.irccloud.com/chat/say", body)
			if err != nil {
				// handle err
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)
			if err != nil {
				// handle err
			}
			defer resp.Body.Close()
			//fmt.Println(resp)
			if 1 == 1 {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				bodyString := string(bodyBytes)
				_ = bodyString
			} else {
				fmt.Println("error")
			}
		}

	}
}
