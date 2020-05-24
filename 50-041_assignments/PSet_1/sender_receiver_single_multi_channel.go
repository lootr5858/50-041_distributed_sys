/* !!!!! ----- Testing codes for 50.041 PSET 1 ----- !!!!! */

/*  !!! --- PACKAGES --- !!!
      Main: create a executable portion of the script
*/
package main

/*  !!! --- Dependencies --- !!! */
import
(
  "fmt"
  "time"
	"math/rand"
)

/*  !!! --- FUNCTIONS --- !!! */
// Single channel SENDER
func sender_single (id int, channel chan int) {
  // Send message repeatedly (with random interval)
  for i := 0; ;i++ {
    // message: random integer
    rand.Seed(time.Now().UnixNano())
    min := 0
    max := 100
    msg := rand.Intn(max - min + 1) + min
    fmt.Println("Client ", id, " sending: ", msg)

    // Increase clock count
    channel <-msg

    // Generate random sleep time (0 - 1000ms)
    latency := time.Duration(rand.Intn(500))
    time.Sleep(time.Millisecond * latency)
  }
}

// Multi channel SENDER
func sender_multi (id int, channel chan []int) {
  // Send message repeatedly (with random interval)
  msg := []int{0,0,0,0}
  rand.Seed(time.Now().UnixNano())
  min := 0
  max := 100

  for i := 0; ;i++ {
    // message: list of random integers
    for j := 0;  j < 4; j++ {
      msg[j] = rand.Intn(max - min + 1) + min
    }

    fmt.Println("Client ", id, " sending: ", msg)

    // Increase clock count
    channel <-msg

    // Generate random sleep time (0 - 1000ms)
    latency := time.Duration(rand.Intn(2000))
    time.Sleep(time.Millisecond * latency)
  }
}

// Single channel RECEIVER
func receiver_single (id int, channel chan int) {
  for i:= 0; ; i++ {
    //  Constantly listening for message from channel id
    msg := <- channel
    fmt.Println("Client ", id, " receiving: ", msg)
  }
}

// multi channel RECEIVER
func receiver_multi (id int, channel chan []int) {
  for i:= 0; ; i++ {
    //  Constantly listening for message from channel id
    msg := <- channel
    fmt.Println("Client ", id, " receiving: ", msg)
  }
}

/*  !!! --- EXECUTE --- !!! */
func main() {
  // create a channel of type integer
  var channel_0 chan int = make(chan int)
  var channel_1 chan []int = make(chan []int)

  // Create link b/w sender & receiver
  go sender_single   (0, channel_0)
  go receiver_single (0, channel_0)

  go sender_multi (1, channel_1)
  go receiver_multi (1, channel_1)

  var input string
  fmt.Scanln(&input)
}
