package connection

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	cookie "github.com/akshayaky/gICCa/cookie"
)

//Connect disconnects or reconnects from a connection whose ID is given
func Connect(session string, cid int, ReorDis string) {
	client := cookie.SetCookie(session, ReorDis+"connect")
	body := strings.NewReader(`cid=` + strconv.Itoa(cid) + `&session=` + session)
	req, err := http.NewRequest("POST", "https://www.irccloud.com/chat/"+ReorDis+"connect", body)
	if err != nil {
		fmt.Println("error sending the request :  ", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	_ = bodyBytes

}
