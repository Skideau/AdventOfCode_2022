package main

import(
	"log"
	"bufio"
	"os"
	"fmt"
)

func checkSliceUnicity(initStr string, windowSize int) (marker string, bufferSize int) {
	for charCounter := 0 ; charCounter < len(initStr) - windowSize + 1 ; charCounter++ {
		marker = initStr[charCounter:charCounter + windowSize]

		charMap := map[byte]int {}
		for sliceUnitCounter := 0 ; sliceUnitCounter < len(marker) ; sliceUnitCounter++ {
			val := marker[sliceUnitCounter]
			_, ok := charMap[val]
			if ok {
				charMap[val]++
			} else {
				charMap[val] = 1
			}
		}
		duplicateChar := false
		for _, val := range charMap {
			if val > 1 {
				duplicateChar = true
				break
			}
		}
		if !duplicateChar {
			bufferSize = charCounter + windowSize
			break
		}
	}
	return
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentLine string
	for scanner.Scan() {
		currentLine = scanner.Text()
	}

	fmt.Printf("\n=====PART1=====\n")
	packetMarker, packetBufferSize := checkSliceUnicity(currentLine, 4)
	fmt.Printf("Found packet start sequence '%s' with buffer size = %d", packetMarker, packetBufferSize)
	
	fmt.Printf("\n=====PART2=====\n")
	messageMarker, messageBufferSize := checkSliceUnicity(currentLine, 14)
	fmt.Printf("Found message start sequence '%s' with buffer size = %d", messageMarker, messageBufferSize)
}