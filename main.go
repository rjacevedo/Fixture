package main

import (
	"./fixture"
	"fmt"
	"time"
)

func test() {
	for {
		time.Sleep(time.Second)
	}
}
func main() {
	go fixture.RunAnalysis(10, 20, 5, "r.jon.acevedo@gmail.com")

	for i:= 1; i <= 30; i++ {
		go test()
	}

	for {
		fmt.Println("Main Program run")
		time.Sleep(time.Second)
	}
}