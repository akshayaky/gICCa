package login

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/inancgumus/screen"
	"github.com/manifoldco/promptui"
	"golang.org/x/crypto/ssh/terminal"
)

//getAuthToken returns Auth token
func getAuthToken() string {

	req, err := http.NewRequest("POST", "https://www.irccloud.com/chat/auth-formtoken", nil)
	if err != nil {
		fmt.Println("\nError sending request(1)")
	}
	req.Header.Set("Content-Length", "0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("\nError sending request(2)")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		ind := strings.Index(bodyString, "token")
		token := bodyString[ind+8 : ind+52]
		//fmt.Println(token)
		return token
	} else {
		fmt.Println("\nerror")
		return "nil"
	}

}

//getSessionID returns session id from the Auth token
func getSessionID(tok1 string, email string, pass string) string {
	body := strings.NewReader(`email=` + email + `&password=` + pass + `&token=` + tok1)
	req, err := http.NewRequest("POST", "https://www.irccloud.com/chat/login", body)
	if err != nil {
		fmt.Println("S errer")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Auth-Formtoken", tok1)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("\nError authenticating")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		ind := strings.Index(bodyString, "session")
		session := bodyString[ind+10 : ind+44]

		return session

	} else {
		fmt.Println("\nError authenticating")
		return "nil"
	}

}

//val  validates the given input
func val(lab string) string {
	validate := func(input string) error {

		if input == "" {
			return errors.New("empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    lab,
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}

//Login returns session ID. Checks session.txt for the session key,
//if not found gets auth key and session ID and then stores and returns
func Login() string {
	dat, _ := ioutil.ReadFile("session.txt")
	screen.Clear()
	screen.MoveTopLeft()
	var session string
	if dat == nil {

		k := "Username"
		email := val(k)
		fmt.Printf("PassWord : ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			fmt.Println("There was an error : ", err)
			os.Exit(0)
		}
		tok := getAuthToken()
		tok1 := tok[0 : len(tok)-1]

		session = getSessionID(string(tok1), email, string(password))
		if session == "nil" {
			fmt.Println("Error occured")
			os.Exit(1)
		}
		f, err := os.Create("session.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		_, err = f.WriteString(session)
		if err != nil {
			fmt.Println(err)
			f.Close()
			os.Exit(3)
		}

		return session
	}
	return string(dat)
}
