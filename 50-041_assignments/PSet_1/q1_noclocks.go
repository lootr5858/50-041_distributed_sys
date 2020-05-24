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
func randomInt(id int, min int, max int) int {
	// generate random msg
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(max-min+1) + min

	return value
}

/* --- Client ---
   Periodically send msg

   Always receives msg from server, unless sending msg
*/
func Client(id int, CS chan int, SC chan int) {
	var receive int

	go ClientSend(id, CS)
	go ClientReceive(id, SC, receive)
}

func ClientSend(id int, channel chan int) {
	for {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(5000)))
		send := randomInt(id, 0, 100)
		channel <- send
		fmt.Println("Client ", id, " sending: ", send, "\n")
	}
}

func ClientReceive(id int, channel chan int, receive int) {
	for {
		select {
		case receive := <-channel:
			fmt.Println("Client ", id, " receives: ", receive, "\n")
		}
	}
}

/* --- Server ---
   Receives messages from 1 node
   Delay with random time units
   Broadcast to all other nodes

   Simultaneously receives message from multiple nodes
   Each msg broadcast delayed Randomly
*/
func Server(SC0 chan int, SC1 chan int, CS0 chan int, CS1 chan int) {
	for {
		select {
		case msg0 := <-CS0:
			fmt.Println("Server received from client_0: ", msg0, "\n")

			// Send to client_1 & _2
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			fmt.Println("Server sending: ", msg0, "\n")
			SC0 <- msg0

		case msg1 := <-CS1:
			fmt.Println("Server received from client_1: ", msg1, "\n")

			// send to client_0 & _2
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			fmt.Println("Server sending: ", msg1, "\n")
			SC0 <- msg1

		case <-time.After(time.Second):
			fmt.Println("No message received!\n")
		}
	}
}

/*  !!! --- EXECUTE --- !!! */
func main() {
	// create a channel of type integer
	var SC0 chan int = make(chan int)
	var SC1 chan int = make(chan int)
	var CS0 chan int = make(chan int)
	var CS1 chan int = make(chan int)

	// Create link b/w sender & receiver
	go Client(0, CS0, SC0)
	go Client(1, CS1, SC1)
	go Server(SC0, SC1, CS0, CS1)

	var input string
	fmt.Scanln(&input)
}
