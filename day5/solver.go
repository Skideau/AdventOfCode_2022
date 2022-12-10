package main

import(
	"fmt"
	"log"
	"bufio"
	"os"
	"strings"
	"strconv"
)

const EMPTY_CRATE byte = byte(0)

type crateStack struct {
	nbCrates int
	crates []byte
}

type cranePlayGround struct {
	nbColumn int
	cratesColumn []crateStack
	pickedCrate byte
}

type instruction struct {
	nbCratesToMove int
	originColumn int
	destinationColumn int
}
/**
 * CRANE PLAYGROUND METHODS
**/
func (crane *cranePlayGround) MoveCrateFor9000(startColumn int, endColumn int) bool {
	if startColumn >= crane.nbColumn || endColumn >= crane.nbColumn {
		print("\t[ERROR] Crane - MoveCrate - column index wrong.\n")
		return false
	}

	popedCrate, poped := crane.cratesColumn[startColumn].PopCrate()
	if poped {
		crane.pickedCrate = popedCrate
		crane.cratesColumn[endColumn].PushCrate(crane.pickedCrate)
		return true
	}
	return false
}

func (crane *cranePlayGround) MoveCrateFor9001(nbCratesToMove int, startColumn int, endColumn int) bool {
	if startColumn >= crane.nbColumn || endColumn >= crane.nbColumn {
		print("\t[ERROR] Crane - MoveCrate - column index wrong.\n")
		return false
	}

	popedCrates, poped := crane.cratesColumn[startColumn].PopNCrates(nbCratesToMove)
	if poped {
		crane.pickedCrate = popedCrates[0]
		crane.cratesColumn[endColumn].PushNCrates(popedCrates)
		return true
	}
	return false
}

func (crane *cranePlayGround) ExecuteInstructionFor9000(inst instruction) bool {
	for crateCounter := 0 ; crateCounter < inst.nbCratesToMove ; crateCounter++ {
		if !crane.MoveCrateFor9000(inst.originColumn - 1, inst.destinationColumn - 1) {
			return false
		}
	}
	return true
}

func (crane *cranePlayGround) ExecuteInstructionFor9001(inst instruction) bool {
	if !crane.MoveCrateFor9001(inst.nbCratesToMove, inst.originColumn - 1, inst.destinationColumn - 1) {
		return false
	}
	return true
}

func (crane *cranePlayGround) GetTopCrates() (result string) {
	for columnCounter := 0 ; columnCounter < crane.nbColumn ; columnCounter++ {
		result += crane.cratesColumn[columnCounter].GetTopCrate()
	}
	return
}

func (crane *cranePlayGround) Show() (result string) {
	for columnCounter := 0 ; columnCounter < crane.nbColumn ; columnCounter++ {
		result += "[" + fmt.Sprintf("%v", columnCounter) + "] : " + crane.cratesColumn[columnCounter].Show()
	}
	return
}

/**
 * CRATE STACK METHODS
**/
func (cs *crateStack) PushCrate(newCrate byte) {
	if len(cs.crates) <= cs.nbCrates {
		cs.crates = append(cs.crates, newCrate)
	} else {
		cs.crates[cs.nbCrates] = newCrate
	}
	cs.nbCrates++
}

func (cs *crateStack) PushNCrates(newCrates []byte) {
	for crateCounter := len(newCrates) - 1 ; crateCounter >= 0 ; crateCounter-- {
		cs.PushCrate(newCrates[crateCounter])
	}
}

func (cs *crateStack) PopCrate() (popedCrate byte, success bool) {
	if len(cs.crates) > 0 && cs.nbCrates > 0 && len(cs.crates) >= cs.nbCrates {
		cs.nbCrates--
		popedCrate = cs.crates[cs.nbCrates]
		cs.crates[cs.nbCrates] = EMPTY_CRATE
		success = true
		return
	}
	popedCrate = EMPTY_CRATE
	success = false
	return
}

func (cs *crateStack) PopNCrates(nbCratesToPop int) (popedCrates []byte, success bool) {
	for crateCounter := 0 ; crateCounter < nbCratesToPop ; crateCounter++ {
		popedCrate, ok := cs.PopCrate()
		if ok {
			popedCrates = append(popedCrates, popedCrate)
		}
	}
	if len(popedCrates) == nbCratesToPop {
		success = true
		return
	}
	success = false
	return
}

func (cs crateStack) GetTopCrate() string {
	return string(cs.crates[cs.nbCrates - 1])
}

func (cs crateStack) Show() (crateStringified string) {
	crateStringified = "_" + fmt.Sprintf("%v", cs.nbCrates) + "_ {"
	for crateIdx := 0 ; crateIdx < cs.nbCrates ; crateIdx++ {
		crateStringified += " [" + string(cs.crates[crateIdx]) + "] "
	}
	crateStringified += "}\n"
	return
}

func initPlayground(crane *cranePlayGround) {
	crane.nbColumn = 9

	newCrateStack := crateStack {}
	newCrateStack.PushCrate(byte('W'))
	newCrateStack.PushCrate(byte('M'))
	newCrateStack.PushCrate(byte('L'))
	newCrateStack.PushCrate(byte('F'))
	crane.cratesColumn = append(crane.cratesColumn, newCrateStack)

	newCrateStack = crateStack {}
	newCrateStack.PushCrate(byte('B'))
	newCrateStack.PushCrate(byte('Z'))
	newCrateStack.PushCrate(byte('V'))
	newCrateStack.PushCrate(byte('M'))
	newCrateStack.PushCrate(byte('F'))
	crane.cratesColumn = append(crane.cratesColumn, newCrateStack)
	
	newCrateStack = crateStack {}
	newCrateStack.PushCrate(byte('H'))
	newCrateStack.PushCrate(byte('V'))
	newCrateStack.PushCrate(byte('R'))
	newCrateStack.PushCrate(byte('S'))
	newCrateStack.PushCrate(byte('L'))
	newCrateStack.PushCrate(byte('Q'))
	crane.cratesColumn = append(crane.cratesColumn, newCrateStack)

	newCrateStack = crateStack {}
	newCrateStack.PushCrate(byte('F'))
	newCrateStack.PushCrate(byte('S'))
	newCrateStack.PushCrate(byte('V'))
	newCrateStack.PushCrate(byte('Q'))
	newCrateStack.PushCrate(byte('P'))
	newCrateStack.PushCrate(byte('M'))
	newCrateStack.PushCrate(byte('T'))
	newCrateStack.PushCrate(byte('J'))
	crane.cratesColumn = append(crane.cratesColumn, newCrateStack)

	newCrateStack = crateStack {}
	newCrateStack.PushCrate(byte('L'))
	newCrateStack.PushCrate(byte('S'))
	newCrateStack.PushCrate(byte('W'))
	crane.cratesColumn = append(crane.cratesColumn, newCrateStack)

	newCrateStack = crateStack {}
	newCrateStack.PushCrate(byte('F'))
	newCrateStack.PushCrate(byte('V'))
	newCrateStack.PushCrate(byte('P'))
	newCrateStack.PushCrate(byte('M'))
	newCrateStack.PushCrate(byte('R'))
	newCrateStack.PushCrate(byte('J'))
	newCrateStack.PushCrate(byte('W'))
	crane.cratesColumn = append(crane.cratesColumn, newCrateStack)

	newCrateStack = crateStack {}
	newCrateStack.PushCrate(byte('J'))
	newCrateStack.PushCrate(byte('Q'))
	newCrateStack.PushCrate(byte('C'))
	newCrateStack.PushCrate(byte('P'))
	newCrateStack.PushCrate(byte('N'))
	newCrateStack.PushCrate(byte('R'))
	newCrateStack.PushCrate(byte('F'))
	crane.cratesColumn = append(crane.cratesColumn, newCrateStack)

	newCrateStack = crateStack {}
	newCrateStack.PushCrate(byte('V'))
	newCrateStack.PushCrate(byte('H'))
	newCrateStack.PushCrate(byte('P'))
	newCrateStack.PushCrate(byte('S'))
	newCrateStack.PushCrate(byte('Z'))
	newCrateStack.PushCrate(byte('W'))
	newCrateStack.PushCrate(byte('R'))
	newCrateStack.PushCrate(byte('B'))
	crane.cratesColumn = append(crane.cratesColumn, newCrateStack)

	newCrateStack = crateStack {}
	newCrateStack.PushCrate(byte('B'))
	newCrateStack.PushCrate(byte('M'))
	newCrateStack.PushCrate(byte('J'))
	newCrateStack.PushCrate(byte('C'))
	newCrateStack.PushCrate(byte('G'))
	newCrateStack.PushCrate(byte('H'))
	newCrateStack.PushCrate(byte('Z'))
	newCrateStack.PushCrate(byte('W'))
	crane.cratesColumn = append(crane.cratesColumn, newCrateStack)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var instructions []instruction

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		lineArr := strings.Fields(currentLine)
		if len(lineArr) < 1 || lineArr[0] != "move" {
			continue
		}

		var newInstruction instruction
		newInstruction.nbCratesToMove, _ = strconv.Atoi(lineArr[1])
		newInstruction.originColumn, _ = strconv.Atoi(lineArr[3])
		newInstruction.destinationColumn, _ = strconv.Atoi(lineArr[5])
		instructions = append(instructions, newInstruction)

	}
	print("\n=====PART1=====\n")

	crane9000 := cranePlayGround {}
	initPlayground(&crane9000)
	print("START STATE\n")
	print(crane9000.Show())
	
	for instructionIdx, currentInstruction := range instructions {
		if !crane9000.ExecuteInstructionFor9000(currentInstruction) {
			fmt.Printf("\tFailed to execute instruction '%d'\n", instructionIdx)
			break
		}
	}
	print("\nEND STATE\n")
	print(crane9000.Show())
	print("\n")
	print(crane9000.GetTopCrates())
	
	print("\n=====PART2=====\n")
	crane9001 := cranePlayGround {}
	initPlayground(&crane9001)
	print("START STATE\n")
	print(crane9001.Show())
	for instructionIdx, currentInstruction := range instructions {
		if !crane9001.ExecuteInstructionFor9001(currentInstruction) {
			fmt.Printf("\tFailed to execute instruction '%d'\n", instructionIdx)
			break
		}
	}
	print("\nEND STATE\n")
	print(crane9001.Show())
	print("\n")
	print(crane9001.GetTopCrates())
}