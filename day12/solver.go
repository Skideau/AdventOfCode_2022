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
	weight  int // manhattan dist to grid goal
	height  int
	isStart bool
	isGoal  bool
}

func createPosition(rowIdx int, colIdx int, heightChar rune) position {
	startPos := false
	goalPos := false

	if heightChar == 'S' {
		// Start Position
		startPos = true
		heightChar = 'a'
	}
	if heightChar == 'E' {
		// Goal Position
		goalPos = true
		heightChar = 'z'
	}
	positionHeight := int(heightChar-'0') - int('a'-'0')
	return position{colIdx, rowIdx, 0, positionHeight, startPos, goalPos}
}

func (pos position) ManhDist(targetPos position) int {
	return int(math.Abs(float64(targetPos.posX)-float64(pos.posX)) + math.Abs(float64(targetPos.posY)-float64(pos.posY)))
}

func (pos position) Equal(targetPos position) bool {
	return (pos.posX == targetPos.posX) && (pos.posY == targetPos.posY)
}

func (pos position) Add(posTransform [2]int) position {
	return position{pos.posX + posTransform[0], pos.posY + posTransform[1], 0, 0, false, false}
}

func (target *position) Reachable(source position) (isReachable bool) {
	isReachable = target.height >= source.height && target.height <= source.height+1
	if isReachable && target.height < source.height {
		target.weight *= 10
	}
	return
}

func (pos position) String() string {
	return fmt.Sprintf("[%d ; %d] = {%d ; %d}", pos.posX, pos.posY, pos.height, pos.weight)
}

type grid struct {
	heightMap map[int]position
	start     position
	goal      position
	mapWidth  int
	mapHeight int
}

func (areaMap grid) InBound(newPos position) bool {
	return (newPos.posX >= 0 && newPos.posX < areaMap.mapWidth) && (newPos.posY >= 0 && newPos.posY < areaMap.mapHeight)
}

func (areaMap grid) FindPos(fakePos position) position {
	mapIdx := fakePos.posX + fakePos.posY*areaMap.mapWidth
	return areaMap.heightMap[mapIdx]
}

func (areaMap *grid) UpdateHeightMap(positionArray []position) {
	for _, position := range positionArray {
		goalWeight := position.ManhDist(areaMap.goal)
		position.weight = goalWeight
		mapIdx := position.posX + position.posY*areaMap.mapWidth
		areaMap.heightMap[mapIdx] = position
	}
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
		if !source.Equal(fakePos) && areaMap.InBound(fakePos) {
			gridPos := areaMap.FindPos(fakePos)
			if gridPos.Reachable(currentPos) {
				neighbours = append(neighbours, gridPos)
				neighbourCounter++
			}
		}
	}
	return neighbourCounter, neighbours
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
	strMap += "{\n"
	for mapIdx := 0; mapIdx < len(areaMap.heightMap); mapIdx++ {
		mapPos := areaMap.heightMap[mapIdx]
		if mapPos.posX == 0 {
			strMap += "  [ "
		}

		newWeight := fmt.Sprint(mapPos.weight)
		if len(newWeight) <= 1 {
			newWeight = "0" + newWeight
		}
		strMap += newWeight + " "

		if mapPos.posX == areaMap.mapWidth-1 {
			strMap += "] \n"
		}
	}
	strMap += "}\n"
	return
}

type pathNode struct {
	pathIndex            int
	pos                  position
	nbAvailableNeighbour int
	neighbours           []position
	triedPath            []position
}

func (node *pathNode) GetNextBestNeighbour() (position, bool) {
	fakePos := position{79, 18, 0, 0, false, false}
	// blockingNode := false
	if node.pos.Equal(fakePos) {
		// blockingNode = true
		fmt.Printf("Node %v \n\twith neighbours: %q\n", node.pos, node.neighbours)
	}
	if len(node.neighbours) <= 0 {
		// No more neighbours available
		// fmt.Printf("\t\tNo more neighbours listed for pos %v\n", node.pos)
		return position{}, false
	}

	chosenNeighbourIdx := 0
	chosenNeighbour := node.neighbours[chosenNeighbourIdx]
	for neighbourIdx, neighbour := range node.neighbours {
		if neighbour.weight < chosenNeighbour.weight {
			chosenNeighbourIdx = neighbourIdx
			chosenNeighbour = neighbour
		}
	}
	node.neighbours = append(node.neighbours[:chosenNeighbourIdx], node.neighbours[chosenNeighbourIdx+1:]...)
	node.triedPath = append(node.triedPath, chosenNeighbour)
	// fmt.Printf("\t\tFound %d available neigh and %d tried ones for pos %v\n", len(node.neighbours), len(node.triedPath), node.pos)
	return chosenNeighbour, true
}

type mapPath []pathNode

func (myPath mapPath) IsPosInPath(testPos position) bool {
	posFound := false
	for _, node := range myPath {
		if node.pos.Equal(testPos) {
			// fmt.Printf("\t\tFound in path[%d]: %v\n", len(myPath), node.pos)
			posFound = true
			break
		}
	}
	return posFound
}

func (myPath mapPath) GetLastTreatableNodeIdx() (int, []pathNode) {
	// fmt.Printf("Current path : %v\n", myPath)
	rejectedNodes := []pathNode{}
	for nodeCounter := len(myPath) - 1; nodeCounter >= 0; nodeCounter-- {
		currentNode := myPath[nodeCounter]
		if currentNode.nbAvailableNeighbour > 1 && len(currentNode.triedPath) < len(currentNode.neighbours) {
			// fmt.Printf("\tWent back to node[%d] to pos: %v with %d available neighbours and %d tried ones\n", nodeCounter, currentNode.pos, len(currentNode.neighbours), len(currentNode.triedPath))
			return nodeCounter, rejectedNodes
		} else {
			rejectedNodes = append(rejectedNodes, currentNode)
		}
	}
	return -1, rejectedNodes
}

func (myPath mapPath) BuildPath(areaMap *grid) (newPath mapPath, pathBlocked bool, nodeBlocked bool) {
	nodeBlocked = false
	pathBlocked = false

	currentNodeIdx, rejectedNodes := myPath.GetLastTreatableNodeIdx()
	if currentNodeIdx < 0 {
		// No more treatable node in path
		pathBlocked = true
		return
	}
	for _, rNode := range rejectedNodes {
		gridPos := areaMap.FindPos(rNode.pos)
		gridPos.weight += 5
		mapIdx := gridPos.posX + gridPos.posY*areaMap.mapWidth
		areaMap.heightMap[mapIdx] = gridPos
	}

	newPath = myPath[:currentNodeIdx+1]

	currentNode := myPath[currentNodeIdx]
	sourcePos := currentNode.pos
	currentPos, neighbourFound := currentNode.GetNextBestNeighbour()
	newPath[currentNode.pathIndex] = currentNode

	if !neighbourFound {
		nodeBlocked = true
		return
	}

	for !currentPos.Equal(areaMap.goal) {
		nbNeighbours, neighbours := (*areaMap).FindReachableNeighbours(sourcePos, currentPos)
		if nbNeighbours <= 0 {
			// fmt.Println("\t\tNo reachable neighbours found.")
			nodeBlocked = true
			return
		}
		currentNode = pathNode{len(newPath), currentPos, nbNeighbours, neighbours, make([]position, 0)}
		sourcePos = currentPos
		currentPos, neighbourFound = currentNode.GetNextBestNeighbour()
		for neighbourFound && newPath.IsPosInPath(currentPos) {
			currentPos, neighbourFound = currentNode.GetNextBestNeighbour()
		}
		if !neighbourFound {
			nodeBlocked = true
			return
		}
		newPath = append(newPath, currentNode)
	}

	return
}

func (myPath mapPath) String() (strMap string) {
	strMap = fmt.Sprintf("mapPath[%d]: { ", len(myPath))
	for _, node := range myPath {
		strMap += fmt.Sprintf("[%d ; %d] ", node.pos.posX, node.pos.posY)
	}
	strMap += "}"
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
	var positionArray []position
	lineCounter := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		areaMap.mapWidth = len(currentLine)
		for charIdx, charData := range currentLine {
			newPos := createPosition(lineCounter, charIdx, charData)
			positionArray = append(positionArray, newPos)
			if newPos.isStart {
				areaMap.start = newPos
			} else if newPos.isGoal {
				areaMap.goal = newPos
			}
		}
		lineCounter++
	}
	areaMap.mapHeight = lineCounter
	areaMap.UpdateHeightMap(positionArray)
	// fmt.Printf("%s", &areaMap)

	// Initialize path with source position
	sourceNeighboursNb, sourceNeighbours := areaMap.FindReachableNeighbours(areaMap.start, areaMap.start)
	sourceNode := pathNode{0, areaMap.start, sourceNeighboursNb, sourceNeighbours, make([]position, 0)}
	var gridPath mapPath
	gridPath = append(gridPath, sourceNode)

	// Proceed to build path
	pathBlocked := false
	nodeBlocked := false
	nbBlockedNode := 0
	for !pathBlocked {
		gridPath, pathBlocked, nodeBlocked = gridPath.BuildPath(&areaMap)
		if !nodeBlocked {
			break
		}
		nbBlockedNode++
		fmt.Printf("Path has been reset %d times\n", nbBlockedNode)
		fmt.Printf("\tLast path %v\n", gridPath)
		// if math.Mod(float64(nbBlockedNode), 10) == 0 {
		// 	fmt.Printf("Path has been reset %d times\n", nbBlockedNode)
		// 	fmt.Printf("\tLast path %v\n", gridPath)
		// }
		// fmt.Println("Node blocked, trying another way")
	}

	if pathBlocked {
		fmt.Printf("[ERROR] Path blocked without finding any access to the final goal\n")
	}

	// fmt.Println("=====PART1=====")
	// fmt.Printf("Current gridPath has %d nodes, the last one is %v\n", len(gridPath), gridPath[len(gridPath)-1].pos)
	// if gridPath[len(gridPath)-1].pos.Equal(areaMap.goal) {
	// 	fmt.Printf("Number of steps required to reach goal: %d\n", len(gridPath)-1)
	// }
}
