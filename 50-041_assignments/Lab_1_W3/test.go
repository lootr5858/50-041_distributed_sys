package main

import "fmt" 
import "time"

func main() {
	var cnt int8

	cnt = 0

	for true {
		cnt += 1

		fmt.Printf("%d Pew Pew Pew!\n", cnt)

		time.Sleep(10 * time.Millisecond)
	}
}
