package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func monitorFuncPerf(funcName string) func() {
	funcStart := time.Now()
	return func() {
		fmt.Printf("Function '%s' took: %v\n", funcName, time.Since(funcStart))
	}
}

func main() {
	defer monitorFuncPerf("main")()

	file, err := os.Open("../inputs/day13/testInput.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
	}
}
