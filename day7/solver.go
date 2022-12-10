package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const MAX_FOLDER_SIZE int = 100000

type command struct {
	name string
	args []string
}

type fs_file struct {
	name      string
	extension string
	size      int
}

type fs_folder struct {
	name       string
	path       []string
	folderList []string
	fileList   []fs_file
	size       int
}

/*
 * fs_folder methods
**/

func (folder *fs_folder) CheckAvailableFolder(newFolder string) bool {
	for _, listedFolder := range folder.folderList {
		if newFolder == listedFolder {
			return true
		}
	}
	return false
}

func (folder *fs_folder) navigateFS(direction string) (newPath []string, rootPath bool) {
	newPath = folder.path
	rootPath = false

	if direction == ".." {
		if len(folder.path) > 1 {
			// Go up one folder
			newPath = folder.path[:len(folder.path)-1]
			// Try to update folder size ? Checking nested folders if they are already listed in folderSizeMap
		} else {
			rootPath = true
			fmt.Println("Current path is at root level.")
		}
		return
	}

	if direction[0] == '/' {
		// If Absolute path --> build new path from direction instruction
		newPath = append(append(newPath, "/"), strings.FieldsFunc(direction, func(r rune) bool { return r == '/' })...)
		return
	}

	if folder.CheckAvailableFolder(direction) {
		// If folder exists in current location --> go to it
		newPath = append(folder.path, direction)
		return
	}

	fmt.Println("[ERROR] File system navigation did not find any good rule to navigate")
	return
}

func (folder *fs_folder) UpdateContent(commandLine []string) {
	if commandLine[0] == "dir" {
		// Adding sufolder
		folder.folderList = append(folder.folderList, commandLine[1])
	} else {
		// Adding files
		fileDetails := strings.FieldsFunc(commandLine[1], func(r rune) bool { return r == '.' })
		if len(fileDetails) < 2 {
			fileDetails = append(fileDetails, "")
		}
		fileSize, _ := strconv.Atoi(commandLine[0])
		folder.size += fileSize
		folder.fileList = append(folder.fileList, fs_file{fileDetails[0], fileDetails[1], fileSize})
	}
}

func (folder *fs_folder) CheckSubfoldersSize(folderSizeMap *map[string]int) bool {
	_, currentFolderFound := (*folderSizeMap)[stringifyPath(folder.path)]
	if len(folder.folderList) > 0 && !currentFolderFound {
		subFolderSizeArray := []int{}
		cumulatedSize := 0
		for _, subFolderName := range folder.folderList {
			subFolderSize, subfolderFound := (*folderSizeMap)[stringifyPath(append(folder.path, subFolderName))]
			if subfolderFound {
				subFolderSizeArray = append(subFolderSizeArray, subFolderSize)
				cumulatedSize += subFolderSize
			} else {
				break
			}
		}
		if len(subFolderSizeArray) == len(folder.folderList) {
			(*folderSizeMap)[stringifyPath(folder.path)] = cumulatedSize + folder.size
			return true
		}
	}
	return false
}

func (folder *fs_folder) ShowContent() {
	fmt.Printf("%s:\n", stringifyPath(folder.path))
	fmt.Printf("\tFolders: %v\n", folder.folderList)
	fmt.Printf("\tFiles: %v\n", folder.fileList)
}

func (folder *fs_folder) ShowSubFolders(folderSizeMap *map[string]int) {
	fmt.Printf("%s - Subfolders: [", folder.name)
	for _, subFolderName := range folder.folderList {
		subFolderSize, subFolderFound := (*folderSizeMap)[stringifyPath(append(folder.path, subFolderName))]
		if subFolderFound {
			fmt.Printf(" {%s ; %d} ", subFolderName, subFolderSize)
		} else {
			fmt.Printf(" {%s ; _} ", subFolderName)
		}
	}
	fmt.Println(" ]")
}

/*
 * fs_file methods
**/

/*
 * free methods
**/

func stringifyPath(pathArray []string) (pathStr string) {
	pathStr = "/"
	for pathCounter := 1; pathCounter < len(pathArray); pathCounter++ {
		pathStr += pathArray[pathCounter] + "/"
	}
	return
}

func handleNewCommand(newInstruction []string, currentCommand *command, currentPath *[]string, currentFolder *fs_folder, folderMap *map[string]fs_folder, folderSizeMap *map[string]int) {
	if currentCommand.name == "ls" {
		// Update folder map after ls command is finished
		if len(currentFolder.folderList) == 0 {
			// Folder has no subFolders --> Update folderSizeMap using nested files
			_, folderExists := (*folderSizeMap)[stringifyPath(currentFolder.path)]
			if !folderExists {
				(*folderSizeMap)[stringifyPath(currentFolder.path)] = currentFolder.size
			}
		}
		(*folderMap)[stringifyPath(*currentPath)] = *currentFolder
	}

	currentCommand.name = newInstruction[1]
	currentCommand.args = newInstruction[2:]

	if currentCommand.name == "cd" {
		if len(*currentPath) == 0 {
			// Initializing currentPath
			*currentPath = append([]string{}, currentCommand.args[0])
		} else {
			// Updating currentPath
			lastFolder, _ := (*folderMap)[stringifyPath(*currentPath)]
			newPath, _ := lastFolder.navigateFS(currentCommand.args[0])
			if len(newPath) < len(*currentPath) {
				newFolder := (*folderMap)[stringifyPath(newPath)]
				_ = newFolder.CheckSubfoldersSize(folderSizeMap)
			}
			*currentPath = newPath
		}
		currentFolder.name = currentCommand.args[0]
		currentFolder.path = *currentPath
		currentFolder.size = 0
		currentFolder.folderList = make([]string, 0)
		currentFolder.fileList = make([]fs_file, 0)
	}
}

/*
 * MAIN
**/

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Key must be folder absolute path (handle possible duplicate)
	folderMap := map[string]fs_folder{}
	folderSizeMap := map[string]int{}
	var currentFolder fs_folder
	var currentPath []string
	var currentCommand command

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine := scanner.Text()
		instructionElements := strings.Fields(currentLine)
		// println(currentLine)
		if instructionElements[0] == "$" {
			handleNewCommand(instructionElements, &currentCommand, &currentPath, &currentFolder, &folderMap, &folderSizeMap)
		} else if currentCommand.name == "ls" {
			currentFolder.UpdateContent(instructionElements)
		}
	}
	_, exist := folderMap[stringifyPath(currentPath)]
	if !exist {
		// Updates last folder checked
		folderMap[stringifyPath(currentPath)] = currentFolder
		_, folderFound := folderSizeMap[stringifyPath(currentPath)]
		if !folderFound && len(currentFolder.folderList) == 0 {
			folderSizeMap[stringifyPath(currentPath)] = currentFolder.size
		}
	}

	rootFound := len(currentPath) == 1
	for !rootFound {
		// Updates parents' folder size
		currentPath, rootFound = currentFolder.navigateFS("..")
		currentFolder = folderMap[stringifyPath(currentPath)]
		_ = currentFolder.CheckSubfoldersSize(&folderSizeMap)
	}

	fmt.Printf("Len(folderMap) = %d\n", len(folderMap))
	fmt.Printf("Len(folderSizeMap) = %d\n", len(folderSizeMap))

	fmt.Printf("\n=====PART1=====\n")
	var smallFolderCumulatedSize int = 0
	for _, currentFolderSize := range folderSizeMap {
		if currentFolderSize <= MAX_FOLDER_SIZE {
			smallFolderCumulatedSize += currentFolderSize
		}
	}
	fmt.Printf("Cumulated folder size = %d\n", smallFolderCumulatedSize)

	fmt.Printf("\n=====PART2=====\n")
	const FS_TOTAL_SIZE int = 70000000
	const FREE_SPACE_NEEDED int = 30000000
	occupiedSpace, _ := folderSizeMap["/"]
	unusedSpace := FS_TOTAL_SIZE - occupiedSpace
	minFolderSizeToDelete := FREE_SPACE_NEEDED - unusedSpace
	fmt.Printf("Root folder takes %d space.\nUnused space = %d\nA folder of size at least %d needs to be deleted.\n", occupiedSpace, unusedSpace, minFolderSizeToDelete)
	fittingFoldersPaths := []string{}
	minFolderSize := FREE_SPACE_NEEDED
	for folderPath, folderSize := range folderSizeMap {
		if folderSize >= minFolderSizeToDelete {
			fittingFoldersPaths = append(fittingFoldersPaths, folderPath)
			if folderSize < minFolderSize {
				minFolderSize = folderSize
			}
		}
	}
	fmt.Printf("Found %d folders matching this need.\n", len(fittingFoldersPaths))
	for _, folderPath := range fittingFoldersPaths {
		folderSize := folderSizeMap[folderPath]
		fmt.Printf("\t'%s' occupying %d space\n", folderPath, folderSize)
	}
	fmt.Printf("Smallest folder takes %d space.\n", minFolderSize)
}
