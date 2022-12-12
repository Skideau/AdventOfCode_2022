package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

type position struct {
	posX    int
	posY    int
	height  int
	isStart bool
	isGoal  bool
}

func (pos position) ManhDist(targetPos position) int {
	return int(math.Abs(float64(targetPos.posX)-float64(pos.posX)) + math.Abs(float64(targetPos.posY)-float64(pos.posY)))
}

func (pos position) Equal(targetPos position) bool {
	return (pos.posX == targetPos.posX) && (pos.posY == targetPos.posY)
}

func (pos position) Add(posTransform [2]int) position {
	return position{pos.posX + posTransform[0], pos.posY + posTransform[1], 0, false, false}
}

func (pos position) Reachable(target position) bool {
	return target.height <= pos.height+1
}

func (pos position) String() string {
	return fmt.Sprintf("[%d ; %d] = %d", pos.posX, pos.posY, pos.height)
}

type grid struct {
	heightMap map[int]position
	start     position
	goal      position
	mapWidth  int
	mapHeight int
}

func (areaMap grid) ValidPos(newPos position) bool {
	return (newPos.posX >= 0 && newPos.posX < areaMap.mapWidth) && (newPos.posY >= 0 && newPos.posY < areaMap.mapHeight)
}

func (areaMap grid) FindPos(fakePos position) position {
	mapIdx := fakePos.posX + fakePos.posY*areaMap.mapWidth
	return areaMap.heightMap[mapIdx]
}

func (areaMap *grid) AddHeightData(rowIdx int, colIdx int, heightChar rune) bool {
	if areaMap.mapWidth == 0 {
		fmt.Printf("[ERROR] Map width not set\n")
		return false
	}

	startPos := false
	goalPos := false

	if heightChar == 'S' {
		// Start Position
		startPos = true
		areaMap.start = position{colIdx, rowIdx, 0, startPos, goalPos}
		heightChar = 'a'
	}
	if heightChar == 'E' {
		// Goal Position
		goalPos = true
		areaMap.goal = position{colIdx, rowIdx, int('z'-'0') - int('a'-'0'), startPos, goalPos}
		heightChar = 'z'
	}

	mapIdx := colIdx + rowIdx*areaMap.mapWidth
	positionHeight := int(heightChar-'0') - int('a'-'0')
	areaMap.heightMap[mapIdx] = position{colIdx, rowIdx, positionHeight, startPos, goalPos}

	return true
}

func (areaMap grid) FindReachableNeighbours(source position, currentPos position) (int, []position) {
	neighbours := []position{}
	if currentPos.ManhDist(source) > 1 {
		fmt.Printf("[ERROR] Current position isn't next to source position\n")
		return 0, neighbours
	}
	ADJACENT_POS_TRANSFORMATION := [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	neighbourCounter := 0
	for _, posTrans := range ADJACENT_POS_TRANSFORMATION {
		fakePos := currentPos.Add(posTrans)
		if !source.Equal(fakePos) && areaMap.ValidPos(fakePos) {
			gridPos := areaMap.FindPos(fakePos)
			fmt.Printf("\tChecking gridPos: %v\n", gridPos)
			if currentPos.Reachable(gridPos) {
				neighbours = append(neighbours, gridPos)
				neighbourCounter++
			}
		}
	}
	return neighbourCounter, neighbours
}

func (areaMap grid) SelectBestNeighbour(currentPos position, neighbours []position, posHistory []position) position {
	distArr := make([]int, len(neighbours))
	minDist := 10 * areaMap.start.ManhDist(areaMap.goal)
	var chosenNeigh position = currentPos
	for nIdx, neighbour := range neighbours {
		if areaMap.goal.Equal(neighbour) {
			return neighbour
		}
		alreadyVisited := false
		for _, historicPos := range posHistory {
			if neighbour.Equal(historicPos) {
				alreadyVisited = true
				fmt.Printf("\t\tAlready visited neighbour: %v\n", neighbour)
				break
			}
		}
		if alreadyVisited {
			continue
		}
		newDist := areaMap.goal.ManhDist(neighbour)
		distArr[nIdx] = newDist
		if newDist < minDist {
			minDist = newDist
			chosenNeigh = neighbour
			fmt.Printf("\t\tSelecting neighbour: %v\n", chosenNeigh)
		}
	}
	return chosenNeigh
}

func (areaMap *grid) String() (strMap string) {
	strMap = "{\n"
	strMap += fmt.Sprintf("  Start: [%d, %d] = %d\n", areaMap.start.posY, areaMap.start.posX, areaMap.start.height)
	strMap += fmt.Sprintf("  Goal : [%d, %d] = %d\n", areaMap.goal.posY, areaMap.goal.posX, areaMap.goal.height)
	strMap += "}\n"
	strMap += "{\n"
	for mapIdx := 0; mapIdx < len(areaMap.heightMap); mapIdx++ {
		mapPos := areaMap.heightMap[mapIdx]
		if mapPos.posX == 0 {
			strMap += "  [ "
		}

		newHeight := fmt.Sprint(mapPos.height)
		if len(newHeight) <= 1 {
			newHeight = "0" + newHeight
		}
		strMap += newHeight + " "

		if mapPos.posX == areaMap.mapWidth-1 {
			strMap += "] \n"
		}
	}
	strMap += "}\n"
	return
}

func monitorFuncPerf(funcName string) func() {
	funcStart := time.Now()
	return func() {
		fmt.Printf("Function '%s' took: %v\n", funcName, time.Since(funcStart))
	}
}

func main() {
	defer monitorFuncPerf("main")()
	file, err := os.Open("../inputs/day12/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	areaMap := grid{}
	areaMap.heightMap = map[int]position{}
	lineCounter := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		areaMap.mapWidth = len(currentLine)
		for charIdx, charData := range currentLine {
			areaMap.AddHeightData(lineCounter, charIdx, charData)
		}
		lineCounter++
	}
	areaMap.mapHeight = lineCounter
	fmt.Printf("%s", &areaMap)

	sourcePos := areaMap.start
	currentPos := sourcePos
	posCounter := 0
	posHistory := []position{}
	fmt.Printf("Journey:\n")
	fmt.Printf("  [0] %v\n", currentPos)
	for !currentPos.Equal(areaMap.goal) {
		nbPos, neighbours := areaMap.FindReachableNeighbours(sourcePos, currentPos)
		if nbPos <= 0 {
			fmt.Printf("[ERROR] Journey stopped -> No reachable position\n")
			break
		}
		sourcePos = currentPos
		currentPos = areaMap.SelectBestNeighbour(currentPos, neighbours, posHistory)
		posHistory = append(posHistory, currentPos)
		posCounter++
		fmt.Printf("  [%d] %v\n", posCounter, currentPos)
		if posCounter > 30 {
			break
		}
	}

	fmt.Println("=====PART1=====")
	fmt.Printf("Number of steps required to reach goal: %d\n", posCounter)
}
