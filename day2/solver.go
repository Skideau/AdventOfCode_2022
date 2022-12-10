package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func playRoundOne(plays []string) int {
	score := 0
	switch plays[1] {
	case "X":
		// ROCK
		score = 1
		if plays[0] == "A" {
			score += 3
		} else if plays[0] == "C" {
			score += 6
		}
	case "Y":
		// PAPER
		score = 2
		if plays[0] == "B" {
			score += 3
		} else if plays[0] == "A" {
			score += 6
		}
	case "Z":
		// SCISSOR
		score = 3
		if plays[0] == "C" {
			score += 3
		} else if plays[0] == "B" {
			score += 6
		}
	}
	return score
}

func playRoundTwo(plays []string) int {
	score := 0
	switch plays[1] {
	case "X":
		// LOSE
		score = 0
		switch plays[0] {
		case "A":
			//ROCK
			score += 3
		case "B":
			// PAPER
			score += 1
		case "C":
			// SCISSOR
			score += 2
		}
	case "Y":
		// DRAW
		score = 3
		switch plays[0] {
		case "A":
			//ROCK
			score += 1
		case "B":
			// PAPER
			score += 2
		case "C":
			// SCISSOR
			score += 3
		}
	case "Z":
		// WIN
		score = 6
		switch plays[0] {
		case "A":
			//ROCK
			score += 2
		case "B":
			// PAPER
			score += 3
		case "C":
			// SCISSOR
			score += 1
		}
	}
	return score
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	matchScoreOne := 0
	matchScoreTwo := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		matchScoreOne += playRoundOne(strings.Fields(currentLine))
		matchScoreTwo += playRoundTwo(strings.Fields(currentLine))
	}

	println("=====PART1=====")
	fmt.Printf("Final score: %d\n", matchScoreOne)

	println("\n=====PART2=====")
	fmt.Printf("Final score: %d\n", matchScoreTwo)
}
