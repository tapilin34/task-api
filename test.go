package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "/tasks/56"
	arr := strings.Split(s, "/")
	fmt.Println(arr[2])
}
