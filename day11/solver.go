package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type monkey struct {
	id                  int
	items               []int
	operation           func(int) int
	testLvl             int
	targetMonkeyId      [2]int
	nbItemInvestigation int
}

func (mk *monkey) throw(monkeyArr *[]monkey) {
	itemToThrow := mk.items[0]
	mk.items = mk.items[1:]
	targetMonkey := 0
	if math.Mod(float64(itemToThrow), float64(mk.testLvl)) == 0 {
		targetMonkey = mk.targetMonkeyId[0]
	} else {
		targetMonkey = mk.targetMonkeyId[1]
	}
	(*monkeyArr)[targetMonkey].items = append((*monkeyArr)[targetMonkey].items, itemToThrow)
}

func (mk *monkey) investigateItems(monkeyArr *[]monkey) {
	// Monkey inspects first item of the list
	for _, item := range mk.items {
		item = mk.operation(item)
		// Monkey gets bored
		item = int(math.Floor(float64(item) / 3))
		mk.items[0] = item
		// Item is tested and throwed
		mk.throw(monkeyArr)
		mk.nbItemInvestigation++
	}
}

func createMonkey(monkeyDesc []string) (newMonkey monkey) {
	newMonkey = monkey{}
	newMonkey.nbItemInvestigation = 0
	for _, newLine := range monkeyDesc {
		lineArr := strings.Fields(newLine)
		switch lineArr[0] {
		case "Monkey":
			newMonkey.id, _ = strconv.Atoi(strings.TrimFunc(lineArr[1], func(r rune) bool { return !unicode.IsDigit(r) }))
		case "Starting":
			strippedItemLine := strings.FieldsFunc(newLine, func(r rune) bool { return r == ':' })
			itemStrList := strings.FieldsFunc(strippedItemLine[1], func(r rune) bool { return r == ',' })
			for _, item := range itemStrList {
				newItem, _ := strconv.Atoi(item[1:])
				newMonkey.items = append(newMonkey.items, newItem)
			}
		case "Operation:":
			strippedOperationLine := strings.FieldsFunc(newLine, func(r rune) bool { return r == '=' })
			operationElements := strings.Fields(strippedOperationLine[1])
			newMonkey.operation = fctBuilder_operation(operationElements[1], operationElements[2])
		case "Test:":
			newMonkey.testLvl, _ = strconv.Atoi(lineArr[len(lineArr)-1])
		case "If":
			if lineArr[1] == "true:" {
				newMonkey.targetMonkeyId[0], _ = strconv.Atoi(lineArr[len(lineArr)-1])
			} else {
				newMonkey.targetMonkeyId[1], _ = strconv.Atoi(lineArr[len(lineArr)-1])
			}
		}
	}
	return
}

func fctBuilder_operation(operator string, operand string) func(int) int {
	if operand != "old" {
		secondOperator, _ := strconv.Atoi(operand)
		switch operator {
		case "+":
			return func(worriedLvl int) int { return worriedLvl + secondOperator }
		case "*":
			return func(worriedLvl int) int { return worriedLvl * secondOperator }
		default:
			fmt.Println("[ERROR] Function builder: Operand not listed ")
			return nil
		}
	} else {
		switch operator {
		case "+":
			return func(worriedLvl int) int { return worriedLvl + worriedLvl }
		case "*":
			return func(worriedLvl int) int { return worriedLvl * worriedLvl }
		default:
			fmt.Println("[ERROR] Function builder: Operand not listed ")
			return nil
		}
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	monkeyDesc := make([][]string, 1)
	monkeyCounter := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		if currentLine != "" {
			monkeyDesc[monkeyCounter] = append(monkeyDesc[monkeyCounter], currentLine)
		} else {
			monkeyDesc = append(monkeyDesc, []string{})
			monkeyCounter++
		}
	}

	monkeys := []monkey{}
	for _, desc := range monkeyDesc {
		monkeys = append(monkeys, createMonkey(desc))
	}

	for roundCounter := 0; roundCounter < 20; roundCounter++ {
		for monkeyIdx, monkey := range monkeys {
			monkey.investigateItems(&monkeys)
			monkeys[monkeyIdx] = monkey
		}
		// Debug prints for each round results
		// fmt.Printf("Round [%d] results:\n", roundCounter)
		// for _, monkey := range monkeys {
		// 	fmt.Printf("\tMonkey[%d] items: %v\n", monkey.id, monkey.items)
		// }
	}

	println("=====PART1=====")
	investigationStats := []int{}
	fmt.Println("Monkeys' items investigation statistics:")
	for _, monkey := range monkeys {
		fmt.Printf("\tMonkey[%d] inespected items %d times\n", monkey.id, monkey.nbItemInvestigation)
		investigationStats = append(investigationStats, monkey.nbItemInvestigation)
	}
	sort.Ints(investigationStats)
	monkeyBusiness := investigationStats[len(investigationStats)-1] * investigationStats[len(investigationStats)-2]
	fmt.Printf("Monkey business gives: %d x %d = %d\n", investigationStats[len(investigationStats)-1], investigationStats[len(investigationStats)-2], monkeyBusiness)

	println("\n=====PART2=====")
}
