package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type instruction struct {
	name    string
	value   int
	nbCycle int
}

func main() {
	file, err := os.Open("../inputs/day10/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var registerX []int
	lastRegisterVal := 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		cmdArray := strings.Fields(currentLine)
		var currentInstruction instruction
		if len(cmdArray) > 1 {
			// ADDX instruction
			addVal, _ := strconv.Atoi(cmdArray[1])
			currentInstruction = instruction{cmdArray[0], addVal, 2}
		} else {
			// NOOP instruction
			currentInstruction = instruction{cmdArray[0], 0, 1}
		}

		for cycleCounter := 0; cycleCounter < currentInstruction.nbCycle; cycleCounter++ {
			registerX = append(registerX, lastRegisterVal)
		}
		lastRegisterVal += currentInstruction.value
	}
	registerX = append(registerX, lastRegisterVal)

	println("=====PART1=====")
	fmt.Printf("[#cycle] x register_value = signal_strength\n")
	cycleOfInterest := [6]int{20, 60, 100, 140, 180, 220}
	var cumulatedSignalStrength int = 0
	for _, cycleIdx := range cycleOfInterest {
		cumulatedSignalStrength += cycleIdx * registerX[cycleIdx-1]
		fmt.Printf("\t[%d] x %d = %d\n", cycleIdx, registerX[cycleIdx-1], cycleIdx*registerX[cycleIdx-1])
	}
	fmt.Printf("Cumulated signal strength = %d\n", cumulatedSignalStrength)
}
