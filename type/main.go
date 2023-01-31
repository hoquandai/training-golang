package main

import (
	"fmt"
)

func main() {
	var greeting interface{} = "hello world"
	greetingStr := greeting.(int)

	fmt.Println(greetingStr)
}
