package concurrent

import (
	"bufio"
	"fmt"
	"io/ioutil"

	connection "github.com/akshayaky/gICCa/connection"
	"github.com/akshayaky/gICCa/message"
	"github.com/akshayaky/gICCa/ui"
	"github.com/jroimartin/gocui"
)

/*
Decider function uses concurrency to decided which
function to execute, viewMessages or Say1.
*/
func Decider(session string, g *gocui.Gui, reader *bufio.Reader, servers [10]message.Makeserver, channels [10]message.ChannelInit) {

	chan1 := make(chan string)
	chan2 := make(chan string)

	var ToName string
	dat, _ := ioutil.ReadFile("defaultChan.txt")
	ToName = string(dat)
	ui.ChangeTitle(&ToName, "mainView", g)
	connection.Connect(session, servers[3].Cid, "re")
	Send := message.SendMessages(session, servers[3].Cid, g)

	g.Update(func(g *gocui.Gui) error {
		v, _ := g.View("channels")
		for _, channel := range channels {
			fmt.Fprintln(v, fmt.Sprintf("%s", channel.Chan))
		}
		return nil
	})
	g.Update(func(g *gocui.Gui) error {
		v, _ := g.View("members")
		for _, member := range channels[0].Members {
			fmt.Fprintln(v, fmt.Sprintf("%s", member.Nick))
		}
		return nil
	})

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
