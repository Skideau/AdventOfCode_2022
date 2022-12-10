package main

import (
	"fmt"
	"math"
)

type ropeEnd struct {
	posX int
	posY int
}

func (ropePart *ropeEnd) Move(direction rune, repetition int) {
	switch direction {
	case 'U': // UP
	case 'R': // RIGHT
	case 'D': // DOWN
	case 'L': // LEFT
	}
}

func manhattanDist(head *ropeEnd, tail *ropeEnd) int {
	return int(math.Abs(float64((*head).posX-(*tail).posX))) + int(math.Abs(float64((*head).posY-(*tail).posY)))
}

func maxDist(head *ropeEnd, tail *ropeEnd) int {
	return int(math.Max(math.Abs(float64((*head).posX-(*tail).posX)), math.Abs(float64((*head).posY-(*tail).posY))))
}

func ropeVect(head *ropeEnd, tail *ropeEnd) ropeEnd {
	return ropeEnd{head.posX - tail.posX, head.posY - tail.posY}
}

func main() {
	// file, err := os.Open("input.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	currentLine := scanner.Text()
	// }
	head := ropeEnd{0, 0}
	tail := ropeEnd{0, 0}
	for y := -2; y <= 2; y++ {
		head.posY = y
		for x := -2; x <= 2; x++ {
			head.posX = x
			fmt.Printf("Head%v, maxDist=%d, manhDist=%d, vect=%v\n", head, maxDist(&head, &tail), manhattanDist(&head, &tail), ropeVect(&head, &tail))
		}
	}
}
