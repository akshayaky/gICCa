package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/akshayaky/gICCa/backlog"
	"github.com/akshayaky/gICCa/connection"
	"github.com/akshayaky/gICCa/login"
	"github.com/inancgumus/screen"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	dat, _ := ioutil.ReadFile("session.txt")
	screen.Clear()
	screen.MoveTopLeft()
	var session string
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

		session = login.GetSessionId(string(tok1), email, string(password))
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
		main()
	}
	session = string(dat)

	cid, name, _ := backlog.EndpointConnection(session)

	//control doesn't reach here
	fmt.Println("In the main function")
	fmt.Println(name)
	fmt.Println(cid)
	fmt.Println(len(cid))

	var options int
	fmt.Scanf("%d", &options)
	connection.Connect(session, cid[options], "re")

}
