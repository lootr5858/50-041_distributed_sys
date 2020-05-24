/* !!!!! ----- Q2. Bully Algorithm ----- !!!!!
   1. Fixed timeout (simulate the behaviour of detecting fault)
   2. Randomly select a GO routine to be faulty node
   3. Consensus across all GO nodes for newly elected coordinator

   - simulate senarios:
		- best case
		- worse case
	- multiple GO routines start the election process simulatneously
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
/* --- Clients ---
Network of 6 servers
Node 6 is DOWN

Best case:
	Node 5 request for coordinator

Worse case:
	Node 1 request for coordinator

Msg:
	{id, value}
	Value:
		1. Request for coordinator
		2. Reject coordinator
		3. Acknowledge coordinator
		4. Broadcast as coordinator
*/
func Client(id int, send [5]chan []int, receive chan []int) {
	msg := []int{id, 0}

	go ClientSend(id, send, msg)
	go ClientReceive(id, receive)
}

func ClientSend(id int, send [5]chan []int, msg []int) {
	for {
		// Periodically send message
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(5000)))

		for i := range send {
			switch {
			case i != id:
				send[i] <- msg
			}
		}

		msg[1]++
	}
}

func ClientReceive(id int, receive chan []int) {
	for {
		select {
		case msg := <-receive:
			fmt.Println("Client ", id, "receives msg from client ", msg[0], " : ", msg)

		default:
			continue
		}
	}
}

/*  !!! --- EXECUTE --- !!! */
func main() {
	// create a channel of type integer
	const num int = 5
	var channel [num]chan []int

	for i := 0; i < num; i++ {
		fmt.Println("Created channel ", i)
		channel[i] = make(chan []int)
	}

	// Create link b/w clients
	var wg sync.WaitGroup
	wg.Add(num)

	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()

			Client(i, channel, channel[i])
		}(i)
	}

	var input string
	fmt.Scanln(&input)
}
