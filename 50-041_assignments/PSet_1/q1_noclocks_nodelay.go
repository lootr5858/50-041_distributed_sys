/* !!!!! ----- Q1. Client-Server Architecture ----- !!!!!
   1 server, multiple Client
   a client --> Server (random interval)
   When msg recevied, server --> all other clients
   Randomly delays each received msg before broadcasting
*/

/*  !!! --- PACKAGES --- !!!
    Main: create a executable portion of the script
*/
package main

/*  !!! --- Dependencies --- !!! */
import (
	"fmt"
	"math/rand"
	"time"
)

/*  !!! --- FUNCTIONS --- !!! */

// Random int generator

func RandomInt(id int, min int, max int) int {
	// generate random msg
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(max-min+1) + min
	return value
}

/* --- Client ---
   Periodically send msg

   Always receives msg from server, unless sending msg
*/

func Client(id int, channel chan int) {
	//var send int

	for {

		select {
		case receive := <-channel:
			fmt.Println("Cllient ", id, " received: ", receive, ".\n")

		default:
			fmt.Println("Nothing from Server.\n")
		}
	}

	for {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(5000)))

		send := RandomInt(id, 0, 100)
		fmt.Println("Client ", id, " sending: ", send, ".\n")
	}
}

/* --- Server ---
   Receives messages from 1 node
   Delay with random time units
   Broadcast to all other nodes

   Simultaneously receives message from multiple nodes
   Each msg broadcast delayed Randomly
*/

func Server(channel0 chan int, channel1 chan int) {
	for {
		select {
		case msg0 := <-channel0:
			fmt.Println("Server received from Client 0: ", msg0, ".\n")

			// random delay before sending
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
			channel1 <- msg0
			fmt.Println("Server broadcasting: ", msg0, ".\n")

		case msg1 := <-channel1:
			fmt.Println("Server received from Client 1: ", msg1, ".\n")

			// random delay before sending
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
			channel0 <- msg1
			fmt.Println("Server broadcasting: ", msg1, ".\n")

		case <-time.After(time.Second):
			fmt.Println("Server did not receive anthing!\n")
		}
	}
}

/*  !!! --- EXECUTE --- !!! */
func main() {
	// create a channel of type integer
	var channel0 chan int = make(chan int)
	var channel1 chan int = make(chan int)

	// Create link b/w sender & receiver
	go Client(0, channel0)
	go Client(1, channel1)
	go Server(channel0, channel1)

	var input string
	fmt.Scanln(&input)
}
