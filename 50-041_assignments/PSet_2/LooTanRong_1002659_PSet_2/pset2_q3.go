/*
	--- !!! Test code for PSet 2 !!! ---
	1. Lamport's Shared Priority queue
		w/o Ricart & Agrawala's opt
	2. Lamport's Shared Priority queue
		w/ Ricart & Agrawala's opt
	3. Centralised server protocol
*/

/*
	!!! Q3 Centralised server protocol !!!
	1. Client send request
	2. Server reply earler received request
	3. Client enters critical section
	4. Client release

	msg{id, num}
		num:
			0. request
			1. acknowledge
			2. release
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// gnerate timestamp
func TimeStamp(id int, startTime int64) int {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	timeStamp := int(currentTime - startTime)

	fmt.Println(timeStamp, "TIMESTAMP of node", id)
	return timeStamp
}

func reQueue(oQ []int) []int {
	var newQ []int

	if len(oQ) > 2 {
		for i := 1; i < len(oQ); i++ {
			newQ = append(newQ, oQ[i])
		}
	}
	return newQ
}

// Listen to message
func Listen(channel chan []int) []int {
	select {
	case receiveMsg := <-channel:
		fmt.Println(receiveMsg)
		return receiveMsg

	default:
		return []int{}
	}
}

func Request(id int, canRequest bool, toServer chan []int) bool {
	requestMsg := []int{id, 0}

	if canRequest {
		fmt.Println(requestMsg, "request >> server, from node", id)
		toServer <- requestMsg
		canRequest = false
	}
	return canRequest
}

func CriticalSection(id int,
	startTime int64,
	canRequest bool,
	canEnter bool,
	toServer chan []int,
	fromServer chan []int) (bool, bool) {

	releaseMsg := []int{id, 2}
	//var ackMsg []int

	fmt.Println("ENTERING critical section by node", id)

	time.Sleep(time.Second * time.Duration(2))

	timeNow := TimeStamp(id, startTime)
	canRequest = true
	canEnter = false

	fmt.Println("Exiting critical section by node", id, "at time", timeNow)
	fmt.Println("Sending release message from node", id)
	toServer <- releaseMsg

	time.Sleep(time.Second * time.Duration(5))

	return canRequest, canEnter
}

func Client(id int, startTime int64, toServer chan []int, fromServer chan []int) {
	// locks
	var canRequest bool = true
	var canEnter bool = false

	// msg
	var ackMsg []int

	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			defer wg.Done()
			switch {
			case i == 0:
				for {
					ackMsg = Listen(fromServer)
					if len(ackMsg) > 0 {
						fmt.Println(canEnter, "canEnter, node", id)
						canEnter = true
						fmt.Println(canEnter, "canEnter, node", id)
						fmt.Println(canRequest, "canRequest, node", id)
					}
				}

			case i == 1:
				for {
					if canRequest {
						fmt.Println("Requesting from node", id)
						canRequest = Request(id, canRequest, toServer)
					}
					if canEnter {
						fmt.Println("Critical section from node", id)
						canRequest, canEnter = CriticalSection(id, startTime, canRequest, canEnter, toServer, fromServer)
					}
				}
			}
		}(i)
	}
	wg.Wait()
}

func queueMsg(rcvMsg []int, queue []int, canAck bool, toServer [3]chan []int) (bool, []int) {
	if len(rcvMsg) > 0 {
		if rcvMsg[1] == 0 {
			fmt.Println("Can server acknowledge?", canAck)
			queue = append(queue, rcvMsg[0])
			fmt.Println("Server add to queue node", rcvMsg[0], "New server queue", queue)
		} else if rcvMsg[1] == 2 {
			fmt.Println("Server released by node", rcvMsg[0], "\nServer moving queue up...")
			queue = reQueue(queue)
			fmt.Println("After requeue, server queue", queue)
			canAck = true
			fmt.Println(canAck, "canAck by server")
		}
	}
	return canAck, queue
}

func SendAck(queue []int, canAck bool, fromServer [3]chan []int) bool {
	ackMsg := []int{0, 1}

	//fmt.Println("Can server ack?", canAck)

	if canAck {
		//fmt.Println("Check server queue", queue)
		// time.Sleep(time.Second)
		// fmt.Println("whats going on???")
		if len(queue) != 0 {
			toSend := queue[0]
			fmt.Println("Server acknowledging node", toSend, "with", ackMsg)
			fromServer[toSend] <- ackMsg
			fmt.Println("Server >> node", toSend, ackMsg)
			//fmt.Println("Server completed ack")
			canAck = false
		} else {
			canAck = true
		}
	}
	return canAck
}

func Server(toServer [3]chan []int, fromServer [3]chan []int) {
	// locks
	var canAck bool = true

	// messages
	//var rcvMsg []int
	var queue []int

	// concurrent events
	const num int = 2
	var wg sync.WaitGroup
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()

			switch {
			//listen to message
			case i == 0:
				for {
					select {
					case rcvMsg := <-toServer[0]:
						fmt.Println("Server received", rcvMsg)
						canAck, queue = queueMsg(rcvMsg, queue, canAck, toServer)
						fmt.Println("CAn server ack?", canAck)

					case rcvMsg := <-toServer[1]:
						fmt.Println("Server received", rcvMsg)
						canAck, queue = queueMsg(rcvMsg, queue, canAck, toServer)
						fmt.Println("CAn server ack?", canAck)

					case rcvMsg := <-toServer[2]:
						fmt.Println("Server received", rcvMsg)
						canAck, queue = queueMsg(rcvMsg, queue, canAck, toServer)
						fmt.Println("CAn server ack?", canAck)

					default:
						time.Sleep(time.Microsecond)
					}
				}

				//queue messages
			case i == 1:
				for {
					canAck = SendAck(queue, canAck, fromServer)
				}
			}
		}(i)
	}
	wg.Wait()
}

func main() {
	/* Define CONSTANTS
	- # of nodes
	- start time	*/
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println("Starting time:", startTime)
	const n int = 3

	/* Create CHANNELS
	- Requests
	- Reply			*/
	var toServer [n]chan []int
	var fromServer [n]chan []int

	for i := 0; i < n; i++ {
		fmt.Println("Creating channel", i)
		toServer[i] = make(chan []int)
		fromServer[i] = make(chan []int)
	}

	go Client(0, startTime, toServer[0], fromServer[0])
	go Client(1, startTime, toServer[1], fromServer[1])
	go Client(2, startTime, toServer[2], fromServer[2])
	go Server(toServer, fromServer)

	var input string
	fmt.Scanln(&input)
}
