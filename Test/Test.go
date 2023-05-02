package main

import (
	"fmt"
	"os"
)

func main() {
	_, err := os.Open("F:\\a.txt")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("ok")
}
