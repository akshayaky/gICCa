package backlog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func GetBacklog(streamid string, session string) string {
	jar, _ := cookiejar.New(nil)

	var cookies []*http.Cookie

	firstCookie := &http.Cookie{
		Name:   "session",
		Value:  session,
		Path:   "/",
		Domain: ".irccloud.com",
	}

	cookies = append(cookies, firstCookie)

	// URL for cookies to remember. i.e reply when encounter this URL
	cookieURL, _ := url.Parse("http://www.irccloud.com/chat/backlog/" + streamid)

	jar.SetCookies(cookieURL, cookies)

	//setup our client based on the cookies data
	client := &http.Client{
		Jar: jar,
	}

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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	scanner := bufio.NewScanner(strings.NewReader(bodyString))
	k := 1
	for scanner.Scan() {
		backlog, _ := json.Marshal(scanner.Text()[0 : len(scanner.Text())-1])
		strback := string(backlog)
		strback = strings.Replace(strback, "\\", "", -1)

		k = k + 1
		if strings.Contains(strback, "cid") {

			//fmt.Println(strback)
			ind := strings.Index(strback, "cid")
			cid := strback[ind+5 : ind+12]
			fmt.Println(cid)
			//_ = cid
			return cid
		}

		//fmt.Println(strback)
	}

	//	lines := string(line)
	//	print(lines)
	//	print("\n")

	//}
	return "nil"

}

func EndpointConnection(session string) string {

	jar, _ := cookiejar.New(nil)

	var cookies []*http.Cookie

	firstCookie := &http.Cookie{
		Name:   "session",
		Value:  session,
		Path:   "/",
		Domain: ".irccloud.com",
	}

	cookies = append(cookies, firstCookie)

	cookieURL, _ := url.Parse("https://www.irccloud.com/chat/stream")

	jar.SetCookies(cookieURL, cookies)

	client := &http.Client{
		Jar: jar,
	}

	urlData := url.Values{}
	urlData.Set("session", session)

	req, err := http.NewRequest("POST", "https://www.irccloud.com/chat/stream", strings.NewReader(urlData.Encode()))
	if err != nil {
		fmt.Println("error sending the request :  ", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(resp.Body)

	for {
		line, _ := reader.ReadBytes('\n')
		lines := string(line)
		ind := strings.Index(lines, "streamid")
		streamid := lines[ind+11 : ind+43]
		//fmt.Println(streamid)

		cid := GetBacklog(streamid, session)
		return cid
		break

	}
	defer resp.Body.Close()
	return ""
}
