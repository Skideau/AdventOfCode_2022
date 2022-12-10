package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	elves := []int{}
	maxCal := 0
	elvesCounter := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		if len(currentLine) == 0 {
			if elves[elvesCounter] > maxCal {
				maxCal = elves[elvesCounter]
			}
			elves = append(elves, 0)
			elvesCounter++
		} else {
			if len(elves) == 0 {
				elves = append(elves, 0)
			}
			newCal, _ := strconv.Atoi(currentLine)
			elves[elvesCounter] += newCal
		}
	}

	println("\n=====PART1=====")
	fmt.Printf("Max calories carried: %d", maxCal)

	println("\n=====PART1=====")
	sort.Ints(elves)
	fmt.Printf("Three top elves carry: %d calories.\n", elves[elvesCounter]+elves[elvesCounter-1]+elves[elvesCounter-2])
}
