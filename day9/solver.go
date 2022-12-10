package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type ropeEnd struct {
	name string
	posX int
	posY int
}

func (ropePart *ropeEnd) MoveDiag(distVect ropeEnd) {
	if !math.Signbit(float64(distVect.posX)) {
		ropePart.posX++
	} else {
		ropePart.posX--
	}

	if !math.Signbit(float64(distVect.posY)) {
		ropePart.posY++
	} else {
		ropePart.posY--
	}
	// fmt.Printf("Diag move: %v\n", ropePart)
}

func (ropePart *ropeEnd) MoveStraight(direction string) {
	switch direction {
	case "U": // UP
		ropePart.posY++
	case "R": // RIGHT
		ropePart.posX++
	case "D": // DOWN
		ropePart.posY--
	case "L": // LEFT
		ropePart.posX--
	}
	// fmt.Printf("Straight move: %v\n", ropePart)
}

func (tail *ropeEnd) Follow(head *ropeEnd) (moved bool) {
	moved = false
	if maxDist(head, tail) > 1 {
		moveVect := ropeVect(head, tail)
		if head.posX == tail.posX {
			// Move Straight vertically
			if !math.Signbit(float64(moveVect.posY)) {
				tail.MoveStraight("U")
			} else {
				tail.MoveStraight("D")
			}
		} else if head.posY == tail.posY {
			// MoveStraight horizontally
			if !math.Signbit(float64(moveVect.posX)) {
				tail.MoveStraight("R")
			} else {
				tail.MoveStraight("L")
			}
		} else {
			// Move diagonaly
			tail.MoveDiag(moveVect)
		}
		moved = true
		return
	}
	return
}

func updatePositionHistoric(posHistory *[]ropeEnd, newPos *ropeEnd) {
	valFound := false
	for _, historicVal := range *posHistory {
		valFound = newPos.posX == historicVal.posX && newPos.posY == historicVal.posY
		if valFound {
			break
		}
	}

	if !valFound {
		*posHistory = append(*posHistory, *newPos)
		// fmt.Printf("History updated with value: %v\n", *newPos)
	}
}

// func manhattanDist(head *ropeEnd, tail *ropeEnd) int {
// 	return int(math.Abs(float64((*head).posX-(*tail).posX))) + int(math.Abs(float64((*head).posY-(*tail).posY)))
// }

func maxDist(head *ropeEnd, tail *ropeEnd) int {
	return int(math.Max(math.Abs(float64((*head).posX-(*tail).posX)), math.Abs(float64((*head).posY-(*tail).posY))))
}

func ropeVect(head *ropeEnd, tail *ropeEnd) ropeEnd {
	return ropeEnd{"Vect", head.posX - tail.posX, head.posY - tail.posY}
}

func main() {
	file, err := os.Open("../inputs/day9/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	head := ropeEnd{"Head", 0, 0}
	tail := ropeEnd{"Tail", 0, 0}
	tailHistory := append([]ropeEnd{}, tail)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		moveInstructions := strings.Fields(currentLine)
		headDirection := moveInstructions[0]
		headMoveRepetition, _ := strconv.Atoi(moveInstructions[1])
		// fmt.Printf("Head moving: %d times direction: %v\n", headMoveRepetition, headDirection)
		for moveCounter := 0; moveCounter < headMoveRepetition; moveCounter++ {
			head.MoveStraight(headDirection)
			if tail.Follow(&head) {
				updatePositionHistoric(&tailHistory, &tail)
			}
		}
	}

	println("=====PART1=====")
	fmt.Printf("Tail tooked %d different positions during heads journey\n", len(tailHistory))

	println("\n=====PART2=====")
}
