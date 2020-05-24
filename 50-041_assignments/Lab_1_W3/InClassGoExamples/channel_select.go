package main

import (
  "fmt"
  "time"
	"math/rand"
)


func sender(id int, clock chan []int) {
	send := []int{1,1,1,1,1,1,1,1,1,1};
  for i := 0; ; i++ {
		// send an integer
		fmt.Println("Channel ", id, " sending ", send)
    clock <- send
		//sleep sporadically
    amt := time.Duration(rand.Intn(1500))
    time.Sleep(time.Millisecond * amt)
		send[id] = send[id]+1
  }
}

func receiver(clock1 chan []int, clock2 chan []int) {
  for {
		// receive message from clock1 or clock2 whichever is ready
		// if both clock1 and clock2 are ready, then break the tie randomly
		select {
      case msg1 := <- clock1:
        fmt.Println("Receiving from channel 1: ", msg1, "\n")
      case msg2 := <- clock2:
        fmt.Println("Receiving from channel 2:", msg2, "\n")
      case <- time.After(time.Second):
	fmt.Println("timeout\n")
    }
  }
}

func main() {
	// create a channel of type integer
  var clock1 chan []int = make(chan []int)
  var clock2 chan []int = make(chan []int)

	// launch two go routines "sender" and "receiver"
  go sender(0, clock1)
  go sender(1, clock2)
  go receiver(clock1, clock2)

  var input string
  fmt.Scanln(&input)
}
