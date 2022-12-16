package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type packetData struct {
	dataType   string
	dataLength int
	elems      []string
}

func (pckData *packetData) ParseElements(packetStr string) bool {
	if packetStr == "" {
		fmt.Println("Error when reading packet data")
		return false
	}

	packetStr = packetStr[1 : len(packetStr)-1]
	for _, strElem := range packetStr {
	}
}

type packetPair struct {
	packet   [2]string
	pcktData [2]packetData
}

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

	var packetPairs []packetPair
	nbPacketPair := 0
	currentPacketPair := packetPair{}
	packetCounter := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		if currentLine == "" {
			nbPacketPair++
			packetCounter = 0
			packetPairs = append(packetPairs, currentPacketPair)
		} else {
			currentPacketPair.packet[packetCounter] = currentLine
			packetCounter++
		}

	}

	for pairIdx, pair := range packetPairs {
		fmt.Printf("[%d]:\n", pairIdx)
		fmt.Printf("\t%s\n", pair.packet[0])
		fmt.Printf("\t%s\n", pair.packet[1])
	}
}
