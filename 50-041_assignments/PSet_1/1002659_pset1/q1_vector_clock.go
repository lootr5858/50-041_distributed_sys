/* !!!!! ----- Q1. Client-Server Architecture ----- !!!!!
   1 server, multiple Client
   a client --> Server (random interval)
   When msg recevied, server --> all other clients
   Randomly delays each received msg before broadcasting

   Lamport's logic clock to determine total order
   Print order
	- to know order of message read
*/
/*  !!! --- PACKAGES --- !!!
    Main: create a executable portion of the script
*/
package main

/*  !!! --- Dependencies --- !!! */
import (
	"fmt"
	"math/rand"
	"sync"
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
func Client(id int, CS chan []int, SC chan []int) {
	send := []int{0, 0, 0, 0}
	var clk int = 0
	var wg sync.WaitGroup

	/*
	   2 simultaneous threads running at the same time
	     - Periodically SENDING message to server
	     - Constantly LISTENING for message from server
	*/
	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func(i int) {
			defer wg.Done()

			switch {
			case i == 0:
				// Periodically send message
				for {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(5000)))

					clk++
					send[id+1] = clk
					send[len(send)-1] = randomInt(id, 0, 100)
					CS <- send
					fmt.Println("Client ", id, " vector clock: ", clk)
					fmt.Println("Client ", id, " sending: ", send, "\n")
					clk++
				}

			case i == 1:
				// Constantly listening for message
				for {
					select {
					case receive := <-SC:
						switch {
						case clk >= receive[id+1]:
							clk++

						case clk < receive[id+1]:
							fmt.Println("ERROR! Message clock came before client ", id, "\n")
							wg.Done()
							break
						}

						fmt.Println("Client ", id, " vector clock: ", clk)
						fmt.Println("Client ", id, " receives: ", receive, "\n")

					default:
						continue
					}
				}
			}
		}(i)
	}
	wg.Wait()
}

/* --- Server ---
   Receives messages from 1 node
   Delay with random time units
   Broadcast to all other nodes

   Simultaneously receives message from multiple nodes
   Each msg broadcast delayed Randomly
*/
func Server(CS0 chan []int, CS1 chan []int, SC0 chan []int, SC1 chan []int) {
	clock := []int{0, 0, 0}

	for {
		send0 := []int{0, 0, 0, 0}
		send1 := []int{0, 0, 0, 0}

		select {
		case msg0 := <-CS0:

			if clock[0] >= msg0[0] {
				clock[0]++
			}

			if clock[1] < msg0[1] {
				clock[1] = msg0[1]
			}

			if clock[2] < msg0[2] {
				clock[2] = msg0[2]
			}

			fmt.Println("Server vector clock: ", clock)
			fmt.Println("Server received from client_0: ", msg0, "\n")

			// Send to client_1 & _2
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			clock[0]++
			send0[0] = clock[0]
			send0[1] = clock[1]
			send0[2] = clock[2]
			send0[len(send0)-1] = msg0[len(msg0)-1]
			SC1 <- send0
			fmt.Println("Server vector clock: ", clock)
			fmt.Println("Server sending: ", send0, "\n")
			clock[0]++

		case msg1 := <-CS1:

			if clock[0] >= msg1[0] {
				clock[0]++
			}

			if clock[1] < msg1[1] {
				clock[1] = msg1[1]
			}

			if clock[2] < msg1[2] {
				clock[2] = msg1[2]
			}

			fmt.Println("Server vector clock: ", clock)
			fmt.Println("Server received from client_1: ", msg1, "\n")

			// send to client_0 & _2
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			clock[0]++
			send1[0] = clock[0]
			send1[1] = clock[1]
			send1[2] = clock[2]
			send1[len(send1)-1] = msg1[len(msg1)-1]
			SC0 <- send1
			fmt.Println("Server vector clock: ", clock)
			fmt.Println("Server sending: ", send1, "\n")
			clock[0]++

		case <-time.After(time.Second):
			clock[0]++
			fmt.Println("Server vector clock: ", clock)
			fmt.Println("No message received!\n")
		}
	}
}

/*  !!! --- EXECUTE --- !!! */
func main() {
	// create a channel of type integer
	var CS0 chan []int = make(chan []int)
	var CS1 chan []int = make(chan []int)
	var SC0 chan []int = make(chan []int)
	var SC1 chan []int = make(chan []int)

	// Create link b/w sender & receiver
	go Client(0, CS0, SC0)
	go Client(1, CS1, SC1)
	go Server(CS0, CS1, SC0, SC1)

	var input string
	fmt.Scanln(&input)
}
