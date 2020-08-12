package connection

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	cookie "github.com/akshayaky/gICCa/cookie"
)

func Reconnect(session string, cid string, ReorDis string) {
	client := cookie.SetCookie(session)
	body := strings.NewReader(`cid=` + cid + `&session=` + session)
	req, err := http.NewRequest("POST", "https://www.irccloud.com/chat/"+ReorDis+"connect", body)
	//fmt.Println(body)
	if err != nil {
		fmt.Println("error sending the request :  ", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	_ = bodyBytes
	//fmt.Println(string(bodyBytes))

}
