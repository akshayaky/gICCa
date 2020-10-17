package main

import (
	"log"

	"github.com/akshayaky/gICCa/backlog"
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

	backlog.EndpointConnection(session, g)

	// var options int
	// fmt.Scanf("%d", &options)
	// connection.Connect(session, cid[options], "re")

}
