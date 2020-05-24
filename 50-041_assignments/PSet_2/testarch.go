// /*
// 	--- !!! Test code for PSet 2 !!! ---
// 	1. Lamport's Shared Priority queue
// 		w/o Ricart & Agrawala's opt
// 	2. Lamport's Shared Priority queue
// 		w/ Ricart & Agrawala's opt
// 	3. Centralised server protocol
// */

// // Define executable package
// package myarch

// // Include dependencies
// import (
// 	"fmt"
// 	"math"
// 	"math/rand"
// 	"sync"
// 	"time"
// )

// /*
// 	--- !!! Lamport's Shared Priority Queue !!! ---
// 		 w/o Ricart & Agrawala's optimisation
// 	Completely distributed system
// 	NO centralised server

// 	1. ALL requests are broadcasted (all machines are aware)
// 	2. Each machine aware of earlier requests

// --- Clients ---
// 5 nodes
// 1. Node 4 request
// 2. Node 2 request

// Request msg:
// 	{id, timestamp}

// Reply msg:
// 	{req_id, timestamp, reply_id, reply}
// */

// // generate random number
// func GenRand(min int, max int) int {
// 	rand.Seed(time.Now().UnixNano())
// 	ranVal := rand.Intn(max-min+1) + min
// 	return ranVal
// }

// // gnerate random delay
// func ranDelay(id int, min int, max int) {
// 	ranVal := GenRand(min, max) * GenRand(min, max)
// 	fmt.Println("Client", id, ": delayed for ", ranVal, "ms ... ...")
// 	time.Sleep(time.Millisecond * time.Duration(ranVal))
// }

// /*
// 	Generate TIMESTAMP
// 1. determine currentTime (in ms)
// 2. timeDiff = currentTime - startTime
// 3. obtain TimeStamp: round TimeDiff ms --> s
// */
// func GenTS(id int, startTime int64) int {
// 	curTime := time.Now().UnixNano() / int64(time.Millisecond)
// 	//fmt.Println("Client ", id, ": current time = ", curTime)

// 	timeDiff := curTime - startTime
// 	//fmt.Println("Client ", id, ": time from start = ", timeDiff)

// 	ts := math.Round(float64(timeDiff) / 1000)
// 	fmt.Println("Client", id, ": timestamp = ", ts)

// 	return int(ts)
// }

// // Sort queue according to timestamp
// func SortQueue(queue [][]int) [][]int {
// 	var temp []int

// 	if len(queue) > 1 {
// 		fmt.Println(queue, len(queue))
// 		for i := 1; i < len(queue)-1; i++ {
// 			for j := 1; j < len(queue)-1; j++ {
// 				temp = queue[j]

// 				if temp[1] > queue[j+1][1] {
// 					queue[j] = queue[j+1]
// 					queue[j+1] = temp
// 				}
// 			}
// 		}
// 	}
// 	return queue
// }

// // Move queue up by 1
// func reQueue(oQ [][]int) [][]int {
// 	var newQ [][]int

// 	if len(oQ) > 2 {
// 		for i := 2; i < len(oQ); i++ {
// 			newQ = append(newQ, oQ[i])
// 		}
// 	}
// 	return newQ
// }

// /*
// 	REQUEST function
// 1. Periodically generate request
// 	- random delay
// 	- with timestamp
// 2. Broadcast to ALL nodes
// 	- including itself
// 3. Wait for replies
// 	- receive majority votes
// 4. Enters critical section
// 5. Exits critical section
// 	- broadcast release message
// */
// func Request(id int,
// 	n int,
// 	startTime int64,
// 	enReq bool,
// 	reqChans [5]chan []int) bool {
// 	if enReq {
// 		// define request message
// 		reqMsg := []int{id}

// 		// create random delay
// 		ranDelay(id, 30, 100)

// 		// get timestamp
// 		TimeStamp := GenTS(id, startTime)

// 		// add timestamp to request_message
// 		reqMsg = append(reqMsg, TimeStamp)
// 		fmt.Println("Client", id, ": BROADCASTING REQUEST ", reqMsg)

// 		// Broadcast request msg to all
// 		for i := 0; i < n; i++ {
// 			if i != id {
// 				reqChans[i] <- reqMsg
// 				fmt.Println("Client", id, ">> node ", i, ": ", reqMsg)
// 			}
// 		}
// 		enReq = false
// 	}

// 	return enReq
// }

// func ListenMajority(id int,
// 	n int,
// 	count int,
// 	startTime int64,
// 	enReq bool,
// 	relChans [5]chan []int,
// 	voteC chan []int) (bool, int) {
// 	var majority = n / 2

// 	// Listen to votes
// 	select {
// 	case vote := <-voteC:
// 		fmt.Println("Client", id, "<< vote from node ", vote[0], "\n", vote)
// 		count++
// 		fmt.Println("Client", id, "counter increased to ", count)
// 	}

// 	if count > majority {
// 		fmt.Println("CLient", id, "accumulated MAJORITY votes!")
// 		// define release msg
// 		relMsg := []int{id}

// 		// generate timestamp before entering critical section
// 		ts := GenTS(id, startTime)

// 		fmt.Println("Client", id, "received MAJORITY votes!!!\nENTERING critical section at ", ts, "...")
// 		time.Sleep(time.Second * time.Duration(5))

// 		// Exits critical section & broadcast releast message
// 		fmt.Println("Client", id, "EXITS critical section!!!\n Broadcasting release message ...")

// 		// add timestamp to release message
// 		ts = GenTS(id, startTime)
// 		relMsg = append(relMsg, ts)

// 		for i := 0; i < n; i++ {
// 			relChans[i] <- relMsg
// 			fmt.Println("Client", id, ">> node ", i, ": ", relMsg)
// 		}
// 		enReq = true
// 		count = 0
// 	}
// 	return enReq, count
// }

// func CriticSection(id int,
// 	n int,
// 	startTime int64,
// 	reqChans [5]chan []int,
// 	relChans [5]chan []int,
// 	voteC chan []int) {

// 	// define enables for functions
// 	var enReq bool = true

// 	var count int = 0

// 	// Create concurrent threads: critical section + listen for lock
// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	for i := 0; i < 2; i++ {
// 		go func(i int) {
// 			defer wg.Done()

// 			switch {
// 			// Request
// 			case i == 0:
// 				for {
// 					//fmt.Println("Client", id, "enReq", enReq, "enRcv", enRcv, "enCri", enCri)
// 					enReq = Request(id,
// 						n,
// 						startTime,
// 						enReq,
// 						reqChans)
// 				}

// 			// Listen for votes && count for majority votes
// 			case i == 1:
// 				for {
// 					enReq, count = ListenMajority(id, n, count, startTime, enReq, relChans, voteC)
// 				}
// 			}
// 		}(i)
// 	}
// 	wg.Wait()
// }

// /*
// 	VOTE function
// C1a. Constantly listen to requests
// C1b. Add requests to queue
// C2a. Vote for request
// 		- with random delay
// 		- 1st in queue
// C2b. Upon voting, lock till release message received
// */
// func ListenReq(id int,
// 	Queue [][]int,
// 	reqC chan []int) [][]int {
// 	select {
// 	case reqMsg := <-reqC:
// 		fmt.Println("Client", id, "<< request from node", reqMsg[0], reqMsg)

// 		fmt.Println("Client", id, "queue =", Queue)
// 		Queue = append(Queue, reqMsg)
// 		//currentQueue = SortQueue(currentQueue)
// 		fmt.Println("Client", id, ": queue = ", Queue)
// 		return Queue

// 	default:
// 		return Queue
// 	}
// }

// func Vote(id int,
// 	qCount int,
// 	startTime int64,
// 	Queue [][]int,
// 	enVote bool,
// 	enRelease bool,
// 	voteChans [5]chan []int) (bool, bool) {
// 	if enVote {
// 		if len(Queue) > qCount {
// 			voteMsg := []int{id}

// 			// delay sending vote
// 			ranDelay(id, 0, 30)

// 			// generate timestamp
// 			ts := GenTS(id, startTime)
// 			voteMsg = append(voteMsg, ts)
// 			fmt.Println("Client", id, "voteMsg =", voteMsg)

// 			// vote for 1st request in the queue
// 			toSend := Queue[qCount][0]
// 			fmt.Println("Client", id, ">> vote for node ", toSend, voteMsg)
// 			voteChans[toSend] <- voteMsg

// 			// disable vote until release message received
// 			enVote = false
// 			enRelease = true
// 		}
// 	}

// 	return enVote, enRelease
// }

// func ReceiveRelease(id int,
// 	qCount int,
// 	Queue [][]int,
// 	enVote bool,
// 	enRelease bool,
// 	relC chan []int) (bool, bool, int) {
// 	select {
// 	// Received release message
// 	// Append previousQueue
// 	// Move up currentQueue
// 	// enVote, disable receiveRelease
// 	case relMsg := <-relC:
// 		if enRelease {
// 			fmt.Println("Client", id, "<< release message from node ", relMsg[0], relMsg)

// 			// prep for voting the next request in queue
// 			qCount++
// 			fmt.Println("Client", id, "queue count  =", qCount)

// 			// enable vote & disable listening for release msg
// 			enVote = true
// 			enRelease = false

// 			return enVote, enRelease, qCount
// 		} else {
// 			return enVote, enRelease, qCount
// 		}

// 	default:
// 		return enVote, enRelease, qCount
// 	}
// }

// func VoteSystem(id int,
// 	startTime int64,
// 	voteChans [5]chan []int,
// 	reqC chan []int,
// 	relC chan []int) {
// 	// function locks
// 	var enVote bool = true
// 	var enRelease bool = false

// 	// request queues
// 	Queue := [][]int{{id}}
// 	var qCount int = 1

// 	// Create concurrent events
// 	var wg sync.WaitGroup
// 	wg.Add(3)

// 	for i := 0; i < 3; i++ {
// 		go func(i int) {
// 			defer wg.Done()

// 			switch {
// 			// Constantly listening to requests
// 			case i == 0:
// 				for {
// 					Queue = ListenReq(id, Queue, reqC)
// 				}

// 			// Vote for request (1st in queue)
// 			case i == 1:
// 				for {
// 					enVote, enRelease = Vote(id,
// 						qCount,
// 						startTime,
// 						Queue,
// 						enVote,
// 						enRelease,
// 						voteChans)
// 				}

// 			// Listen to release (after voting)
// 			case i == 2:
// 				for {
// 					enVote, enRelease, qCount = ReceiveRelease(id,
// 						qCount,
// 						Queue,
// 						enVote,
// 						enRelease,
// 						relC)
// 				}
// 			}
// 		}(i)
// 	}
// 	wg.Wait()
// }

// func Client(id int,
// 	n int,
// 	startTime int64,
// 	reqChans [5]chan []int,
// 	reqC chan []int,
// 	voteChans [5]chan []int,
// 	voteC chan []int,
// 	relChans [5]chan []int,
// 	relC chan []int) {
// 	go CriticSection(id, n, startTime, reqChans, relChans, voteC)
// 	go VoteSystem(id, startTime, voteChans, reqC, relC)
// }

// // Execute main function
// func myarch() {
// 	/* Define CONSTANTS
// 	- # of nodes
// 	- start time	*/
// 	startTime := time.Now().UnixNano() / int64(time.Millisecond)
// 	fmt.Println("Starting time:", startTime)
// 	const n int = 5

// 	/* Create CHANNELS
// 	- Requests
// 	- Reply			*/
// 	var reqChans [n]chan []int
// 	var voteChans [n]chan []int
// 	var relChans [n]chan []int

// 	for i := 0; i < n; i++ {
// 		fmt.Println("Created channel", i)
// 		reqChans[i] = make(chan []int)
// 		voteChans[i] = make(chan []int)
// 		relChans[i] = make(chan []int)
// 	}

// 	// Create nodes
// 	var wg sync.WaitGroup
// 	wg.Add(n)
// 	for i := 0; i < n; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			go Client(i,
// 				n,
// 				startTime,
// 				reqChans,
// 				reqChans[i],
// 				voteChans,
// 				voteChans[i],
// 				relChans,
// 				relChans[i])
// 		}(i)
// 	}

// 	wg.Wait()

// 	var input string
// 	fmt.Scanln(&input)
// }
