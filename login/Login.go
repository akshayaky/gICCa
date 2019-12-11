package login

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/manifoldco/promptui"
)

func GetAuthToken() string { //function tp get the Auth token

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

func GetSessionId(tok1 string, email string, pass string) string { // function to get the session id from the Auth token
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

func Val(lab string) string { //function to validate the given input
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
