package concurrent

import (
	"bufio"

	"github.com/akshayaky/gICCa/message"
	"github.com/akshayaky/gICCa/ui"
	"github.com/jroimartin/gocui"
)

/*
Decider function uses concurrency to decided which
function to execute, viewMessages or Say1.
*/
func Decider(session string, g *gocui.Gui, reader *bufio.Reader, cid int) {

	chan1 := make(chan string)
	chan2 := make(chan string)

	var toName string

	// fmt.Printf("\n\nEnter the Name : ")
	// fmt.Scanln(&toName)
	// fmt.Println(toName)
	toName = "chan or nick"
	ui.ChangeTitle(toName, "mainView", g)
	Send := message.SendMessages(session, cid, toName)
	//these two functions will run concurrently
	// go message.SendMessages(session, cid, toName, chan1)
	go message.ViewMessages(session, reader, toName, g, chan1)
	go ui.Control(g, chan2)

	var msg string
	for {
		select {
		case <-chan1:

		case msg = <-chan2:
			Send(msg)
		}
	}
}
