package main

import(
	"log"
	"bufio"
	"os"
	// "fmt"
	// "strings"
	// "strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
	}
}