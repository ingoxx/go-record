package main

// Importing fmt and time
import (
	"fmt"
	"time"
)

// Main function
func main() {

	// Creating a channel
	// Using make keyword
	channel := make(chan string, 2)

	// time.Sleep(time.Second * 3)
	// channel <- "lxb"
	// Select statement
	select {

	// Using case statement to receive
	// or send operation on channel
	case output := <-channel:
		fmt.Println(output)

	// Calling After() method with its
	// parameter
	case <-time.After(5 * time.Second):

		// Printed after 5 seconds
		fmt.Println("Its timeout..")
	}
}
