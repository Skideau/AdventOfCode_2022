package main

import(
	"fmt"
	"log"
	"bufio"
	"os"
	"strings"
	"strconv"
)

type area struct {
	lowerBound int
	upperBound int
}

type myCounter int

func (pairCounter *myCounter) DetectTotallyOverlappedPairArea(pairBounds [2]area) {
	firstIn := pairBounds[0].lowerBound >= pairBounds[1].lowerBound && pairBounds[0].upperBound <= pairBounds[1].upperBound
	secondIn := pairBounds[1].lowerBound >= pairBounds[0].lowerBound && pairBounds[1].upperBound <= pairBounds[0].upperBound

	if firstIn || secondIn {
		*pairCounter++
	}
}

func (pairCounter *myCounter) DetectPartiallyOverlappedPairArea(pairBounds [2]area) {
	firstLowerIn := pairBounds[0].lowerBound >= pairBounds[1].lowerBound && pairBounds[0].lowerBound <= pairBounds[1].upperBound
	firstUpperIn := pairBounds[0].upperBound >= pairBounds[1].lowerBound && pairBounds[0].upperBound <= pairBounds[1].upperBound
	secondLowerIn := pairBounds[1].lowerBound >= pairBounds[0].lowerBound && pairBounds[1].lowerBound <= pairBounds[0].upperBound
	secondUpperIn := pairBounds[1].upperBound >= pairBounds[0].lowerBound && pairBounds[1].upperBound <= pairBounds[0].upperBound

	anyBoundIn := firstLowerIn || firstUpperIn || secondLowerIn || secondUpperIn

	if anyBoundIn {
		*pairCounter++
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	findPairsBounds := func (c rune) bool {
		return c == ','
	}

	findElveArea := func (c rune) bool {
		return c == '-'
	}

	var totallyOverlappedPairArea myCounter
	var partiallyOverlappedPairArea myCounter

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()

		var pairArea [2]area
		pairsBounds := strings.FieldsFunc(currentLine, findPairsBounds)
		for pairIdx, pairBound := range pairsBounds {
			bounds := strings.FieldsFunc(pairBound, findElveArea)
			pairArea[pairIdx].lowerBound, _ = strconv.Atoi(bounds[0])
			pairArea[pairIdx].upperBound, _ = strconv.Atoi(bounds[1])
		}
		totallyOverlappedPairArea.DetectTotallyOverlappedPairArea(pairArea)
		partiallyOverlappedPairArea.DetectPartiallyOverlappedPairArea(pairArea)
	}
	
	fmt.Printf("\n=====PART1=====\n")
	fmt.Printf("Number of assignment pairs fully contained in the other: %v\n", totallyOverlappedPairArea)
	fmt.Printf("\n=====PART2=====\n")
	fmt.Printf("Number of assignment pairs partially overlapping the other: %v\n", partiallyOverlappedPairArea)
}