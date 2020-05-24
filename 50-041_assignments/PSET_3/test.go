package main

import "fmt"

type file struct {
	data  []string
	page  string
	copy  []string
	owner string
}

func main() {
	const page0 string = "0"
	file0 := file{page: page0}
	fmt.Println(file0)
}
