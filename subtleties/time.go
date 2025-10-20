package subtleties

import (
	"fmt"
	"time"
)

func DoneAfter() {
	ch := make(chan string, 1)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "done"
	}()

	select {
	case done := <-ch:
		fmt.Println(">" + done + "<")
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}

/*
The time.After function creates a channel that will be sent a message after x seconds.
When used in combination with a select statement it can be an easy way of setting a deadline for another routine.
*/
