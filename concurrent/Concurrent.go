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
func Decider(session string, g *gocui.Gui, reader *bufio.Reader, servers [10]message.Makeserver) {

	chan1 := make(chan string)
	chan2 := make(chan string)

	var ToName string
	ToName = "#amfoss-test"
	ui.ChangeTitle(&ToName, "mainView", g)
	Send := message.SendMessages(session, servers[3].Cid, g)

	//these two functions will run concurrently
	go ui.Control(g, chan2, &ToName)
	go message.ViewMessages(session, reader, &ToName, g, chan1)
	var msg string
	for {
		select {
		case <-chan1:

		case msg = <-chan2:
			Send(msg, &ToName)
		}
	}
}
