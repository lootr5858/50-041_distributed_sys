/*
	--- !!! Test code for PSet 2 !!! ---
	1. Lamport's Shared Priority queue
		w/o Ricart & Agrawala's opt
	2. Lamport's Shared Priority queue
		w/ Ricart & Agrawala's opt
	3. Centralised server protocol
*/

// Define executable package
package main

// Include dependencies
import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

/*	--- !!! Lamport's Shared Priority Queue !!! ---
		 w/ Ricart & Agrawala's optimisation
	Completely distributed system
	NO centralised server

	1. ALL requests are broadcasted (all machines are aware)
	2. Each machine aware of earlier requests
*/

// generate random number
func RandomNum(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Intn(max-min+1) + min
	return randomValue
}

// generate random delay
func RandomDelay(id string, min int, max int) {
	randomValue := RandomNum(min, max)
	delay := time.Duration(randomValue * randomValue)
	fmt.Println(delay, "DELAY: node", id)
	time.Sleep(time.Millisecond * delay)
}

func DelayById(id string) {
	idx, _ := strconv.Atoi(id)
	delay := time.Duration(idx + 1)
	fmt.Println(delay, "DELAY by id: node", id)
	time.Sleep(time.Second * delay)
}

// gnerate timestamp
func TimeStamp(id string, startTime int64) string {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	timeStamp := int(currentTime - startTime)

	fmt.Println(timeStamp, "TIMESTAMP of node", id)
	return strconv.Itoa(timeStamp)
}

// Sort queue according to timestamp
func SortQueue(queue [][]string) [][]string {
	if len(queue) > 1 {
		for i := 0; i < len(queue)-1; i++ {
			for j := 0; j < len(queue)-1; j++ {
				temp := queue[j]
				tempTs, _ := strconv.Atoi(temp[1])
				nextTs, _ := strconv.Atoi(queue[j+1][1])

				if tempTs > nextTs {
					queue[j] = queue[j+1]
					queue[j+1] = temp
				}
			}
		}
	}
	return queue
}

func ReQueue(oQ [][]string) [][]string {
	var newQ [][]string

	if len(oQ) > 0 {
		for i := 0; i < len(oQ); i++ {
			newQ = append(newQ, oQ[i])
		}
	}
	return newQ
}

// Listen to messages
func Listen(id string, requestC chan []string, replyC chan []string) []string {
	select {
	case msg := <-requestC:
		fmt.Println(msg, ">> node", id)
		return msg

	case msg := <-replyC:
		fmt.Println(msg, ">> node", id)
		return msg

	default:
		return []string{}
	}
}

func checkqueue(id string, startTime int64, requestMsg []string, msg []string, queue [][]string, count int, removeItself bool, replyChans [3]chan []string) (bool, [][]string, int) {
	voteMsg := []string{id, "vote"}

	if removeItself == true {
		queue = ReQueue(queue)
		count = 0
		removeItself = false
		fmt.Println("Removed its own request from queue, node", id, "counter = ", count)
	}

	if len(msg) > 0 {
		if msg[1] == "request" {
			queue = append(queue, msg)
			queue = SortQueue(queue)
			fmt.Println("Added request", msg, ">> queue", queue, "in node", id)

			otherID, _ := strconv.Atoi(msg[0])

			// append timestamp to reqMsg
			currentTime := TimeStamp(id, startTime)
			voteMsg = append(voteMsg, currentTime)

			fmt.Println(otherID, "request replied by node", id)
			replyChans[otherID] <- voteMsg

		} else if msg[1] == "vote" {
			count++
			fmt.Println(count, "counter increased, node", id)
		} else if msg[1] == "release" {
			count++
			queue = ReQueue(queue)
			fmt.Println(queue, "new queue, node", id, "counter = ", count)
		}
	}
	return removeItself, queue, count
}

// Request
func Request(id string, queue [][]string, startTime int64, requestChans [3]chan []string) (bool, []string) {
	requestMsg := []string{id, "request"}
	//idx, _ := strconv.Atoi(id)
	var canRequest bool = false
	var queueItself bool = true

	// periodically request
	DelayById(id)

	// append timestamp to reqMsg
	currentTime := TimeStamp(id, startTime)
	requestMsg = append(requestMsg, currentTime)

	// broadcast to all
	fmt.Println("Broadcasting REQUEST", requestMsg, "from node", id)
	for i := 0; i < 3; i++ {
		requestChans[i] <- requestMsg
	}

	fmt.Println("queueItself", queueItself, "node", id)
	fmt.Println("canRequest", canRequest, "node", id)

	fmt.Println("Request broadcast completed! Node", id, "waiting to enter CRITICAL SECTION!")
	time.Sleep(time.Second * time.Duration(2))
	return canRequest, requestMsg
}

func CriticalSection(id string, queue [][]string, count int, startTime int64, replyChans [3]chan []string) (bool, bool) {
	var removeItself bool = false
	var canRequest bool = false

	//fmt.Println(queue, "node", id)
	//fmt.Println(count, "node", id)

	// If node's request is first in queue
	if len(queue) > 0 {
		if queue[0][0] == id {
			if count > 0 {
				// Enters critical section
				fmt.Println("ENTERING critical section, node", id)
				removeItself = true
				canRequest = true

				fmt.Println("removeItslef", removeItself, "node", id)
				fmt.Println("canRequest", canRequest, "node", id)

				time.Sleep(time.Second)

				timeNow := TimeStamp(id, startTime)

				fmt.Println("EXITIng critical section, node", id, "after time", timeNow)

				if len(queue) > 1 {
					releaseMsg := []string{id, "release"}
					for i := 1; i < len(queue); i++ {
						toReply, _ := strconv.Atoi(queue[i][0])
						fmt.Println(releaseMsg, "release >> node", toReply, "by node", id)
						replyChans[toReply] <- releaseMsg
					}
				}

				fmt.Println("RELESE completed node", id)
			}
		}

	}

	return removeItself, canRequest
}

func Client(id string,
	startTime int64,
	requestChans [3]chan []string,
	replyChans [3]chan []string,
	requestC chan []string,
	replyC chan []string) {

	// constants
	const num int = 3

	// locks
	var canRequest bool = true
	//var queueItself bool = false
	var removeItself bool = false

	// counters
	var count int = 0

	// messages
	var receiveMsg []string
	var requestMsg []string

	// queue
	var queue [][]string

	// Create concurrent threads
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()

			switch {
			// constantly listening
			case i == 0:
				for {
					receiveMsg = Listen(id, requestC, replyC)
				}

				// checkqueue
			case i == 1:
				for {
					removeItself, queue, count = checkqueue(id, startTime, requestMsg, receiveMsg, queue, count, removeItself, replyChans)
				}

				// request & critical section
			case i == 2:
				for {
					if canRequest {
						canRequest, requestMsg = Request(id, queue, startTime, requestChans)
					} else {
						removeItself, canRequest = CriticalSection(id, queue, count, startTime, replyChans)
					}
				}
			}
		}(i)
	}

	wg.Wait()
}

// execute program
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
	var requestChans [n]chan []string
	var replyChans [n]chan []string

	for i := 0; i < n; i++ {
		fmt.Println("Creating channel", i)
		requestChans[i] = make(chan []string)
		replyChans[i] = make(chan []string)
	}

	// Create nodes
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			id := strconv.Itoa(i)
			go Client(id,
				startTime,
				requestChans,
				replyChans,
				requestChans[i],
				replyChans[i])
		}(i)
	}

	wg.Wait()

	var input string
	fmt.Scanln(&input)
}
