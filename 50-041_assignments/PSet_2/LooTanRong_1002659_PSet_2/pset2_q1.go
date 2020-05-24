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
		 w/o Ricart & Agrawala's optimisation
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
	delay := time.Duration(idx)
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

// Listen to messages
func Listen(id string, channel chan []string) []string {
	select {
	case msg := <-channel:
		fmt.Println(msg, ">> node", id)
		return msg

	default:
		return []string{}
	}
}

// Request to enter critical section
func Request(id string, startTime int64, reqChans [3]chan []string) bool {
	requestMsg := []string{id, "request"}
	//idx, _ := strconv.Atoi(id)

	// periodically request
	DelayById(id)

	// append timestamp to reqMsg
	currentTime := TimeStamp(id, startTime)
	requestMsg = append(requestMsg, currentTime)

	// broadcast to all
	fmt.Println("Broadcasting REQUEST", requestMsg, "from node", id)
	for i := 0; i < 3; i++ {
		reqChans[i] <- requestMsg
	}

	fmt.Println("Request broadcast completed! Node", id, "waiting to enter CRITICAL SECTION!")
	return false
}

// check for vote
func VoteCounter(id string, voteCount int, voteMsg []string) int {
	// if message is received
	if len(voteMsg) > 0 {
		fmt.Println(voteMsg, ">> node", id)
		voteCount++
		fmt.Println("Request voted by", voteMsg[0], "\nNode", id, "vote count", voteCount)
	}
	return voteCount
}

// enter critical section & broadcast release message
func CriticalSection(id string, majority int, startTime int64, voteCount int, numCrit int, canRequest bool, releaseChans [3]chan []string) (bool, int) {
	releaseMsg := []string{id, "release"}

	// enters critical section if
	if voteCount >= numCrit*2 {
		fmt.Println("MAJORITY votes for request by node", id, "\nENTERING critical section!")

		// Arbitary delay
		time.Sleep(time.Second)

		// Release votes
		fmt.Println("EXITING critical section ... node", id, "broadcasting to RELEASE votes")

		// append timestamp to release message
		currentTime := TimeStamp(id, startTime)
		releaseMsg = append(releaseMsg, currentTime)

		for i := 0; i < 3; i++ {
			releaseChans[i] <- releaseMsg
		}

		// check time
		timeNow := time.Now().UnixNano() / int64(time.Millisecond)
		timeDifference := timeNow - startTime

		fmt.Println("RELEASE broadcast completed by Node", id, "time taken: ", timeDifference)
		canRequest = true
		numCrit++ // reset vote counter
		fmt.Println(numCrit, "critical section preping. Node", id)
	}
	return canRequest, numCrit
}

// add received request into queue
func QueueRequest(id string, qCount int, requestMsg []string, queue [][]string) [][]string {
	// when a request if received, add to queue
	if len(requestMsg) > 0 {
		fmt.Println("Receive request", requestMsg, ">> node", id)
		queue = append(queue, requestMsg)

		// sort queue according to timeline
		queue = SortQueue(queue)
		fmt.Println(queue[qCount:], "is the current queue of node", id)
		fmt.Println(queue, "is the total queue of node", id)
	}
	return queue
}

// Vote for request
func Vote(id string, queue [][]string, qCount int, startTime int64, voteChans [3]chan []string, releaseC chan []string) int {
	voteMsg := []string{id, "vote"}

	// vote if there are new request in queue
	if len(queue) > qCount {
		fmt.Println(queue[qCount:], "is the current queue of node", id)
		voteFor := queue[qCount][0]
		voteForInt, _ := strconv.Atoi(voteFor)
		voteMsg = append(voteMsg, voteFor)

		RandomDelay(id, 20, 30)

		fmt.Println("VOTING for ", voteFor, "by node", id)

		// append timestamp to reqMsg
		currentTime := TimeStamp(id, startTime)
		voteMsg = append(voteMsg, currentTime)

		voteChans[voteForInt] <- voteMsg
		qCount++
		fmt.Println("Voting by node", id, "completed! Increase queue counter to", qCount)

		var waitingForRelease bool = true

		// wait for release
		for waitingForRelease {
			releaseMsg := Listen(id, releaseC)
			if len(releaseMsg) > 0 {
				fmt.Println(releaseMsg, "release message received by node", id)
				waitingForRelease = false
			}
		}
		fmt.Println("Release for voting ... node", id)
	}
	return qCount
}

func Client(id string,
	n int,
	startTime int64,
	requestChans [3]chan []string,
	voteChans [3]chan []string,
	releaseChans [3]chan []string,
	requestC chan []string,
	voteC chan []string,
	releaseC chan []string) {
	// define constants
	majority := 1
	const numEvents int = 4

	// counters
	var voteCount int
	var qCount int = 0
	var numCrit int = 1

	// locks
	var canRequest bool = true

	// messages
	var requestMsg []string
	var voteMsg []string

	// request queue
	var queue [][]string

	// create concurrent events
	var wg sync.WaitGroup
	wg.Add(numEvents)

	for i := 0; i < numEvents; i++ {
		go func(i int) {
			defer wg.Done()

			switch {
			// Request, enter & exit critical section
			case i == 0:
				for {
					if canRequest {
						// Periodically request
						canRequest = Request(id, startTime, requestChans)
					} else {
						// Enter & exit critical section
						canRequest, numCrit = CriticalSection(id, majority, startTime, voteCount, numCrit, canRequest, releaseChans)
					}
				}

			// Listen for votes & update vote counter
			case i == 1:
				for {
					voteMsg = Listen(id, voteC)
					voteCount = VoteCounter(id, voteCount, voteMsg)
				}

			// Listen & queue request
			case i == 2:
				for {
					requestMsg = Listen(id, requestC)
					queue = QueueRequest(id, qCount, requestMsg, queue)
				}

			// Vote for 1st request in queue & wait for release
			case i == 3:
				for {
					qCount = Vote(id, queue, qCount, startTime, voteChans, releaseC)
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
	var voteChans [n]chan []string
	var releaseChans [n]chan []string

	for i := 0; i < n; i++ {
		fmt.Println("Creating channel", i)
		requestChans[i] = make(chan []string)
		voteChans[i] = make(chan []string)
		releaseChans[i] = make(chan []string)
	}

	// Create nodes
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			id := strconv.Itoa(i)
			go Client(id,
				n,
				startTime,
				requestChans,
				voteChans,
				releaseChans,
				requestChans[i],
				voteChans[i],
				releaseChans[i])
		}(i)
	}

	wg.Wait()

	var input string
	fmt.Scanln(&input)
}
