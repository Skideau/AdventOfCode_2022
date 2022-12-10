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
	name  string
	posX  int
	posY  int
	moved bool
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

func (tail *ropeEnd) Follow(head *ropeEnd) {
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
		tail.moved = true
		return
	}
	tail.moved = false
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

func maxDist(head *ropeEnd, tail *ropeEnd) int {
	return int(math.Max(math.Abs(float64((*head).posX-(*tail).posX)), math.Abs(float64((*head).posY-(*tail).posY))))
}

func ropeVect(head *ropeEnd, tail *ropeEnd) ropeEnd {
	return ropeEnd{"Vect", head.posX - tail.posX, head.posY - tail.posY, false}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Part 1
	// knots := [2]ropeEnd{{"Head", 0, 0, true}, {"Tail", 0, 0, false}}
	// tailHistory := append([]ropeEnd{}, knots[1])

	// Part 2
	knots := [10]ropeEnd{
		{"Head", 0, 0, true},
		{"Tail_1", 0, 0, false},
		{"Tail_2", 0, 0, false},
		{"Tail_3", 0, 0, false},
		{"Tail_4", 0, 0, false},
		{"Tail_5", 0, 0, false},
		{"Tail_6", 0, 0, false},
		{"Tail_7", 0, 0, false},
		{"Tail_8", 0, 0, false},
		{"Tail_9", 0, 0, false}}
	tailHistory := append([]ropeEnd{}, knots[len(knots)-1]) // Keep track of last tail : Tail_9

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		moveInstructions := strings.Fields(currentLine)
		headDirection := moveInstructions[0]
		headMoveRepetition, _ := strconv.Atoi(moveInstructions[1])
		// fmt.Printf("Head moving: %d times direction: %v\n", headMoveRepetition, headDirection)
		for moveCounter := 0; moveCounter < headMoveRepetition; moveCounter++ {
			knots[0].MoveStraight(headDirection)
			for tailCounter := 1; tailCounter < len(knots); tailCounter++ {
				knots[tailCounter].Follow(&knots[tailCounter-1])
			}
			if knots[len(knots)-1].moved {
				updatePositionHistoric(&tailHistory, &knots[len(knots)-1])
			}
		}
	}

	// println("=====PART1=====")
	// fmt.Printf("Tail tooked %d different positions during heads journey\n", len(tailHistory))

	println("\n=====PART2=====")
	fmt.Printf("Tail tooked %d different positions during heads journey\n", len(tailHistory))
}
