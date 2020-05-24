/* ------------------------------------
	!!! --- Ivy Architecture --- !!!
	Basic implementation
	NO fault tolerance

	Central Manager
		- Distribute READ & WRITE access
		- ONLY 1 WRITE access at any point in time

	File metadata
	{page#, copy_sets, owner}
		copy_sets: list of clients with read access to file
		owner: client with write access

	write queue
	{client#, client#, client#, ...}

	messages
	{page#, access}
  ------------------------------------ */

package main

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// required dependencies

/* ------------------------------------
   !!! define datatypes here !!!		*/
type File struct {
	data  string
	page  string
	copy  []string
	owner string
}

type Manager struct {
	id          int
	idStr       string
	fromClients []chan []string
	toClients   []chan []string
	files       []*File
}

type Client struct {
	id         int
	idStr      string
	fromClient chan []string
	toClient   chan []string
}

/*      !!! End of datatype !!!
------------------------------------ */

/* ------------------------------------
   !!! ivy functions here !!!		*/

func DelayById(id int) {
	delayDuration := time.Duration(id) * time.Second
	fmt.Println("Delay", delayDuration, "for client", id)
	time.Sleep(delayDuration)
	fmt.Println("Delay finished for client", id)
}

func (cm *Manager) ManagerListen() string {
	msgs := make([]reflect.SelectCase, len(cm.fromClients))

	var wg sync.WaitGroup
	wg.Add(len(cm.fromClients))
	for i, ch := range cm.fromClients {
		go func(i int) {
			defer wg.Done()
			msgs[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
			// fmt.Println("Central Manager: from Client", i, msgs[i])
		}(i)
	}
	wg.Wait()

	idx, value, _ := reflect.Select(msgs)
	fmt.Println("Central Manager received from client", idx, ":", value)

	return value.String()
}

func (client *Client) ClientRequest() {
	reqMsg := []string{"0", "read"}
	DelayById(client.id)

	fmt.Println("Client", client.id, "requesting", reqMsg)
	client.fromClient <- reqMsg
}

/*      !!! End of functions !!!
------------------------------------ */

// executing program
func main() {
	// define constants
	const numOfClients int = 5

	// create channels
	var fromClients []chan []string
	var toClients []chan []string

	// create clients
	var clients []*Client
	for i := 0; i < numOfClients; i++ {
		fromClients = append(fromClients, make(chan []string))
		toClients = append(toClients, make(chan []string))
		tempClient := Client{
			id:         i,
			idStr:      strconv.Itoa(i),
			fromClient: fromClients[i],
			toClient:   toClients[i],
		}
		clients = append(clients, &tempClient)
	}

	// create files
	var copy0 []string
	file0 := File{
		data:  "haha noob noob de ni",
		page:  "0",
		copy:  copy0,
		owner: "",
	}
	var fileSlice []*File
	fileSlice = append(fileSlice, &file0)

	// create manager
	cm0 := *&Manager{
		id:          0,
		idStr:       "0",
		fromClients: fromClients,
		toClients:   toClients,
		files:       fileSlice,
	}

	var wg sync.WaitGroup
	wg.Add(6)
	for j := 0; j < 6; j++ {
		go func(j int) {
			defer wg.Done()
			if j == 5 {
				for {
					rcvMsg := cm0.ManagerListen()
					fmt.Println("Manager received:", rcvMsg)
				}
			} else if j < 5 {
				for {
					clients[j].ClientRequest()
				}
			}
		}(j)
	}

	wg.Wait()

	var input string
	fmt.Scanln(&input)
}
