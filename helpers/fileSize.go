package main

import (
	"fmt"
	"os"
)

func main() {

	fileName := os.Args[1]

	info, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("File size: %d bytes\n", info.Size())
}
