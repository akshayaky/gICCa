package main

import (
	"fmt"
	"log"

	"github.com/akshayaky/gICCa/backlog"
	"github.com/akshayaky/gICCa/connection"
	"github.com/akshayaky/gICCa/login"
	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.FgColor = gocui.ColorGreen
	g.BgColor = gocui.ColorBlack

	session := login.Login()

	cid, name, _ := backlog.EndpointConnection(session, g)

	//control doesn't reach here
	fmt.Println("In the main function")
	fmt.Println(name)
	fmt.Println(cid)
	fmt.Println(len(cid))

	var options int
	fmt.Scanf("%d", &options)
	connection.Connect(session, cid[options], "re")

}
