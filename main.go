package main

import (
	"fmt"

	"github.com/akshayaky/gICCa/backlog"
	"github.com/akshayaky/gICCa/connection"
	"github.com/akshayaky/gICCa/login"
)

func main() {

	session := login.Login()

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
