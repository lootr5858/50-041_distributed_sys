// // --- !!! For testing codes !!! ---

// // Define executable package
// package main

// // Include dependencies
// import (
// 	"fmt"
// 	"math/rand"
// 	"strconv"
// 	"time"
// )

// // test function
// // generate random number
// func GenRand(min int, max int) int {
// 	rand.Seed(time.Now().UnixNano())
// 	ranVal := rand.Intn(max-min+1) + min
// 	return ranVal
// }

// func genlist() []string {
// 	var temp []string

// 	for j := 0; j < 2; j++ {
// 		val := GenRand(0, 10)
// 		vals := strconv.Itoa(val)
// 		fmt.Println(val, vals)
// 		temp = append(temp, vals)
// 	}

// 	return temp
// }

// func returnEmpty(id string) []string {
// 	if id == "0" {
// 		return []string{}
// 	} else {
// 		return genlist()
// 	}
// }

// // Sort queue according to timestamp
// func SortQueue(queue [][]string) [][]string {

// 	if len(queue) > 1 {
// 		for i := 0; i < len(queue)-1; i++ {
// 			for j := 0; j < len(queue)-1; j++ {
// 				temp := queue[j]
// 				tempTs, _ := strconv.Atoi(temp[1])
// 				nextTs, _ := strconv.Atoi(queue[j+1][1])

// 				if tempTs > nextTs {
// 					queue[j] = queue[j+1]
// 					queue[j+1] = temp
// 				}
// 			}
// 		}
// 	}
// 	return queue
// }

// // executing test code
// func main() {
// 	queue := [][]string{{"0", "8"}, {"1", "5"}, {"2", "7"}}
// 	fmt.Println(queue)
// 	queue = SortQueue(queue)
// 	fmt.Println(queue)
// }
