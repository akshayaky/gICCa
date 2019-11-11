package main

import (
	"fmt"

	login "github.com/login"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	k := "Username"
	email := login.Val(k)
	fmt.Printf("PassWord : ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println("There was an error")
	}
	tok := login.GetAuthToken()
	tok1 := tok[0 : len(tok)-1]

	fmt.Println(tok)
	login.GetSessionId(string(tok1), email, string(password))

}
