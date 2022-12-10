package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type grid struct {
	nbRows int
	nbCols int
	rows   []string
	cols   []string
}

type visibility struct {
	north bool
	east  bool
	south bool
	west  bool
}
type tree struct {
	rowIdx      int
	colIdx      int
	height      int
	exposition  visibility
	scenicScore int
}

func (currentTree *tree) checkVisibility(treeGrid *grid) bool {
	// North visibility
	if currentTree.rowIdx > 0 {
		for _, treeHeight_rune := range treeGrid.cols[currentTree.colIdx][:currentTree.rowIdx] {
			treeHeight := int(treeHeight_rune - '0')
			if treeHeight >= currentTree.height {
				currentTree.exposition.north = false
				break
			}
		}
	}

	// East visibility
	if currentTree.colIdx < treeGrid.nbCols-1 {
		for _, treeHeight_rune := range treeGrid.rows[currentTree.rowIdx][currentTree.colIdx+1:] {
			treeHeight := int(treeHeight_rune - '0')
			if treeHeight >= currentTree.height {
				currentTree.exposition.east = false
				break
			}
		}
	}

	// South visibility
	if currentTree.rowIdx < treeGrid.nbRows-1 {
		for _, treeHeight_rune := range treeGrid.cols[currentTree.colIdx][currentTree.rowIdx+1:] {
			treeHeight := int(treeHeight_rune - '0')
			if treeHeight >= currentTree.height {
				currentTree.exposition.south = false
				break
			}
		}
	}

	// West visibility
	if currentTree.colIdx > 0 {
		for _, treeHeight_rune := range treeGrid.rows[currentTree.rowIdx][:currentTree.colIdx] {
			treeHeight := int(treeHeight_rune - '0')
			if treeHeight >= currentTree.height {
				currentTree.exposition.west = false
				break
			}
		}
	}

	// fmt.Printf("Tree[%d][%d] visibility: %v\n", currentTree.rowIdx, currentTree.colIdx, currentTree.exposition)
	return currentTree.exposition.north || currentTree.exposition.east || currentTree.exposition.south || currentTree.exposition.west
}

func (currentTree *tree) CheckSurroundingTrees(treeGrid *grid) {
	// fmt.Printf("Tree[%d][%d]\n", currentTree.rowIdx, currentTree.colIdx)
	// North visibility
	northTreeScore := 0
	if currentTree.rowIdx > 0 {
		var treeRange string = treeGrid.cols[currentTree.colIdx][:currentTree.rowIdx]
		for treeCounter := len(treeRange) - 1; treeCounter >= 0; treeCounter-- {
			northTreeScore++
			if int(treeRange[treeCounter]-'0') >= currentTree.height {
				break
			}
		}
	}

	// East visibility
	eastTreeScore := 0
	if currentTree.colIdx < treeGrid.nbCols-1 {
		for _, treeHeight_rune := range treeGrid.rows[currentTree.rowIdx][currentTree.colIdx+1:] {
			eastTreeScore++
			if int(treeHeight_rune-'0') >= currentTree.height {
				break
			}
		}
	}

	// South visibility
	southTreeScore := 0
	if currentTree.rowIdx < treeGrid.nbRows-1 {
		for _, treeHeight_rune := range treeGrid.cols[currentTree.colIdx][currentTree.rowIdx+1:] {
			southTreeScore++
			if int(treeHeight_rune-'0') >= currentTree.height {
				break
			}
		}
	}

	// West visibility
	westTreeScore := 0
	if currentTree.colIdx > 0 {
		var treeRange string = treeGrid.rows[currentTree.rowIdx][:currentTree.colIdx]
		for treeCounter := len(treeRange) - 1; treeCounter >= 0; treeCounter-- {
			westTreeScore++
			if int(treeRange[treeCounter]-'0') >= currentTree.height {
				break
			}
		}
	}
	// fmt.Printf("Scenic score: [%d, %d, %d, %d]\n", northTreeScore, eastTreeScore, southTreeScore, westTreeScore)
	currentTree.scenicScore = northTreeScore * eastTreeScore * southTreeScore * westTreeScore
}

func (currentTree *tree) Show() {
	fmt.Printf("Tree[%d][%d]: %d\n", currentTree.rowIdx, currentTree.colIdx, currentTree.height)
	fmt.Printf("Visibility: %v\n", currentTree.exposition)
	fmt.Printf("Scenic score: %d\n", currentTree.scenicScore)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rowCounter := 0
	treeGrid := grid{0, 0, make([]string, 0), make([]string, 0)}
	trees := make([]tree, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		treeGrid.rows = append(treeGrid.rows, currentLine)

		if treeGrid.nbCols == 0 {
			treeGrid.nbCols = len(currentLine)
			for colCounter := 0; colCounter < treeGrid.nbCols; colCounter++ {
				treeGrid.cols = append(treeGrid.cols, "")
			}
		}
		for charIdx, char := range currentLine {
			treeGrid.cols[charIdx] += string(char)
			trees = append(trees, tree{rowCounter, charIdx, int(char - '0'), visibility{true, true, true, true}, 0})
		}
		rowCounter++
	}
	treeGrid.nbRows = rowCounter

	// DEBUG PRINTS
	// println("ROWS")
	// for rowIdx, row := range treeGrid.rows {
	// 	fmt.Printf("[%d] - %s\n", rowIdx, row)
	// }

	// println("COLS")
	// for colIdx, col := range treeGrid.cols {
	// 	fmt.Printf("[%d] - %s\n", colIdx, col)
	// }

	println("\n=====PART1=====")
	visibleTreeCounter := 0
	for _, tree := range trees {
		if tree.checkVisibility(&treeGrid) {
			visibleTreeCounter++
		}
	}
	fmt.Printf("Found %d visible trees.\n", visibleTreeCounter)

	println("\n=====PART2=====")
	maxScenicScore := 0
	mss_tree := tree{}
	for _, tree := range trees {
		tree.CheckSurroundingTrees(&treeGrid)
		if tree.scenicScore > maxScenicScore {
			maxScenicScore = tree.scenicScore
			mss_tree = tree
		}
	}
	println("Max scenic score reached by")
	mss_tree.Show()
}
