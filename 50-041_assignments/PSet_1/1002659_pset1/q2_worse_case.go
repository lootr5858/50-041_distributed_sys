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

Status:
	0. Not coordinator
	1. Requesting as coordinator
	2. Confirm as coordinator
*/
func Client(id int, send [5]chan []int, receive chan []int) {
	var channel chan int = make(chan int)
	var status int = 0

	/* -- Create 2 simultaneous threads --
	1. Initiate bully algo
	2. Constantly listening for request
	*/
	fmt.Println("Client ", id, ": initialising ...")
	go ClientInitiate(id, send, channel)
	go ClientReceive(id, send, receive, channel)
	go ClientStatus(id, send, channel, status)
}

func ClientSendAll(id int, send [5]chan []int, msg []int) {
	fmt.Println("Client ", id, ": broadcasting...")
	for i := range send {
		switch {
		case i != id:
			send[i] <- msg
		}
	}
	fmt.Println("Client ", id, ": finished broadcasting!\n")
}

func ClientSend(id int, send chan []int, msg []int, num int) {
	fmt.Println("Client ", id, ": send to client ", num, " : ", msg)
	send <- msg
}

func ClientReceive(id int, send [5]chan []int, receive chan []int, channel chan int) {
	for {
		select {
		case msg := <-receive:
			fmt.Println("Client ", id, ": receives msg from client ", msg[0], " : ", msg)

			// Execute next action according to msg
			go ClientCheck(id, send, msg, channel)

		case <-time.After(time.Second * time.Duration(3)):
			fmt.Println("Client ", id, ": timeout!")
			continue
		}

		time.Sleep(time.Second)
	}
}

func ClientCheck(id int, send [5]chan []int, rcv []int, channel chan int) {
	switch {
	// received REQUEST
	case rcv[1] == 1:
		switch {
		/* REJECT request & request to be coordinator.
		Receiver has higher id than requester */
		case id > rcv[0]:
			fmt.Println("Client ", id, ": REQUEST REJECTED!!! Receive ", id, " > requester ", rcv[0])
			ClientReject(id, send, rcv[0])
			channel <- 1

		/* Acknowledge request
		Receive has LOWER id than requester */
		case id < rcv[0]:
			fmt.Println("Client ", id, ": Acknowledge request!!! Receive ", id, " < requester ", rcv[0])
			ClientAcknowledge(id, send, rcv[0])
		}

	// receive reject
	case rcv[1] == 2:
		fmt.Println("Client ", id, ": request REJECTEd by client ", rcv[0], "!!!")
		channel <- 0

	// receive acknowledge
	case rcv[1] == 3:
		fmt.Println("Client ", id, ": acknowledged by client ", rcv[0], "!")
		switch {
		// Highest coordinator!!!
		case rcv[0] == 3:
			switch {
			case id == 4:
				fmt.Println("Client 4: broadcasting COORDINATOR status!!!")
				msg := []int{4, 4}
				ClientSendAll(id, send, msg)
				channel <- 2
			}
		}

	case rcv[1] == 4:
		fmt.Println("Client ", id, ": Coordinator is client: ", rcv[0], "!!!")
	}
}

func ClientRequest(id int, send [5]chan []int) {
	msg := []int{id, 1}
	fmt.Println("Client ", id, ": requesting to be coordinator!")
	ClientSendAll(id, send, msg)
}

func ClientReject(id int, send [5]chan []int, num int) {
	msg := []int{id, 2}
	fmt.Println("Client ", id, ": rejects client ", num, "'s coordinator request!")
	ClientSend(id, send[num], msg, num)
}

func ClientAcknowledge(id int, send [5]chan []int, num int) {
	msg := []int{id, 3}
	fmt.Println("Client ", id, ": acknowledge client ", num, "'s coordinator request!")
	ClientSend(id, send[num], msg, num)
}

func ClientInitiate(id int, send [5]chan []int, channel chan int) {
	switch {
	case id == 0:
		fmt.Println("Client ", id, ": detects client 6 DOWN & initiates COORDINATOR ELECTION!")
		channel <- 1

	default:
		fmt.Println("Client ", id, ": waiting ... ...")
		channel <- 0
	}
}

func ClientStatus(id int, send [5]chan []int, channel chan int, status int) {
	for status < 2 {
		select {
		case rcv := <-channel:
			status = rcv

		default:
			continue
		}

		switch {
		case status == 1:
			ClientRequest(id, send)
		}
	}

	fmt.Println("Client ", id, ": becomes COORDINATOR!")
	msg := []int{id, 4}

	for {
		ClientSendAll(id, send, msg)
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
	go Client(0, channel, channel[0])
	go Client(1, channel, channel[1])
	go Client(2, channel, channel[2])
	go Client(3, channel, channel[3])
	go Client(4, channel, channel[4])

	var input string
	fmt.Scanln(&input)
}
