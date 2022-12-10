package main

import(
	"fmt"
	"log"
	"bufio"
	"os"
)

type myString string
type myByteArray []byte

func (s *myString) mapAllCharacter() map[byte]int {
	charMap := map[byte]int {}
	for _, val := range *s {
		_, ok := charMap[byte(val)]
		if ok {
			charMap[byte(val)]++
		} else {
			charMap[byte(val)] = 1
		}
	}
	return charMap
}

func (s *myString) findMatchingChar(charMap map[byte]int) (bool, byte) {
	for _, val := range *s {
		_, ok := charMap[byte(val)]
		if ok {
			return true, byte(val)
		}
	}
	return false, byte(0)
}

func (s *myString) findAllMatchingChar(charMap map[byte]int) (commonByteCounter int, commonBytes []byte) {
	for _, val := range *s {
		_, ok := charMap[byte(val)]
		if ok {
			commonByteCounter++
			commonBytes = append(commonBytes, byte(val))
		}
	}
	return commonByteCounter, commonBytes
}

func addUnique(byteArr *[]byte, elem byte) (*[]byte, bool) {
	var newByte byte
	if elem > 'Z' {
		newByte = elem - 'a' + 1
	} else {
		newByte = elem - 'A' + 27
	}

	if byteArr == nil {
		someArr := make([]byte, 1)
		someArr[0] = newByte
		return &someArr, true
	}
	for _, val := range *byteArr {
		if val == newByte {
			return byteArr, false
		}
	}
	*byteArr = append(*byteArr, newByte)
	return byteArr, true
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice := []byte {}
	groupItems := [3]myString {}
	allGroupItems := [][3]myString {}


	groupCounter := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		if groupCounter < 3 {
			groupItems[groupCounter] = myString(currentLine)
		} else {
			allGroupItems = append(allGroupItems, groupItems)
			groupCounter = 0
			groupItems[groupCounter] = myString(currentLine)
		}
		groupCounter++
		
		lineLength := len(currentLine)
		subStrings := [2]myString {myString(currentLine[:lineLength/2]), myString(currentLine[lineLength/2:])}
		
		charMap := subStrings[0].mapAllCharacter()
		valFound, commonByte := subStrings[1].findMatchingChar(charMap)
		if valFound {
			var newByte byte
			if commonByte > 'Z' {
				newByte = commonByte - 'a' + 1
			} else {
				newByte = commonByte - 'A' + 27
			}
			byteSlice = append(byteSlice, newByte)
		}
	}
	allGroupItems = append(allGroupItems, groupItems)

	if err := scanner.Err() ; err != nil {
		log.Fatal(err)
	}

	charSum := 0
	for _, val := range byteSlice {
		charSum += int(val)
	}

	fmt.Printf("\n=====PART1=====\n")
	fmt.Printf("Total sum: %v\n", charSum)
	
	fmt.Printf("\n=====PART2=====\n")
	charSum = 0
	for _, subItems := range allGroupItems {
		charMap := subItems[0].mapAllCharacter()
		nbCommonBytes := [2]int {}
		commonBytesArr := [2][]byte {}
		nbCommonBytes[0], commonBytesArr[0] = subItems[1].findAllMatchingChar(charMap)
		nbCommonBytes[1], commonBytesArr[1] = subItems[2].findAllMatchingChar(charMap)

	    var allCommonBytes *[]byte
		for _, firstByte := range commonBytesArr[0] {
			for _, secondByte := range commonBytesArr[1] {
				if firstByte == secondByte {
					allCommonBytes, _ = addUnique(allCommonBytes, secondByte)
				}
			}
		}
		for _, cByte := range *allCommonBytes {
			charSum += int(cByte)
		}
	}
	fmt.Printf("Final sum: %v", charSum)
	
}