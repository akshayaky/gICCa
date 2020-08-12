package main

import (
	"fmt"
	"io/ioutil"
	"os"

	backlog "github.com/akshayaky/gICCa/backlog"
	Connection "github.com/akshayaky/gICCa/connection"
	login "github.com/akshayaky/gICCa/login"
	"github.com/akshayaky/gICCa/message"
	"github.com/inancgumus/screen"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	dat, _ := ioutil.ReadFile("session.txt")
	screen.Clear()
	screen.MoveTopLeft()

	if dat == nil {

		k := "Username"
		email := login.Val(k)
		fmt.Printf("PassWord : ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			fmt.Println("There was an error : ", err)
			return
		}
		tok := login.GetAuthToken()
		tok1 := tok[0 : len(tok)-1]

		session := login.GetSessionId(string(tok1), email, string(password))
		_ = session
		if session == "nil" {
			fmt.Println("Error occured")
			return
		}
		f, err := os.Create("session.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		l, err := f.WriteString(session)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		_ = l
	}
	session := string(dat)

	cid := backlog.EndpointConnection(session)
	cid = backlog.EndpointConnection(session)
	fmt.Println(cid)
	fmt.Println(len(cid))

	var options int
	fmt.Scanf("%d", &options)
	Connection.Connect(session, cid, "re")
	message.Say(session, cid)

}
