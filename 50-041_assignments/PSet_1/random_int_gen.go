/* !!!!! ----- Random Integer Generator ----- !!!!! */

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
func random_int (min int, max int) int {
  rand.Seed(time.Now().UnixNano())
  value := rand.Intn(max - min + 1) + min
  fmt.Println("Generate: ", value)

  return value
}

func main () {
  for {
    msg := random_int(0, 10)
    fmt.Println("Received: ", msg)

    latency := time.Duration(rand.Intn(500))
    time.Sleep(time.Millisecond * latency)
  }
}
