package concurrent

import (
	"bufio"
	"fmt"

	"github.com/akshayaky/gICCa/message"
)

/*
Decider function uses concurrency to decided which
function to execute, viewMessages or Say1.
*/
func Decider(session string, reader *bufio.Reader, cid int) {

	chan1 := make(chan string)
	chan2 := make(chan string)

	var toName string

	fmt.Printf("\n\nEnter the Name : ")
	fmt.Scanln(&toName)
	fmt.Println(toName)

	//these two functions will run concurrently
	go message.SendMessages(session, cid, toName, chan1)
	go message.ViewMessages(session, reader, toName, chan2)

	for {
		select {
		case <-chan1:

		case <-chan2:

		}
	}
}
