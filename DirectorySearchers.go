package ghostls

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func RecursiveSearchDir(filepath string) {
	var Directories []string
	var fileArray []string
	files, err := os.ReadDir(filepath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") && !DisplayHidden {
			continue
		}
		fileinfo, err := os.Stat(filepath + "/" + file.Name())
		if err != nil {
			fmt.Println("STAT ERROR")
			log.Fatal(err)
		}
		if fileinfo.IsDir() {
			Directories = append(Directories, file.Name())
		}
		fileArray = append(fileArray, file.Name())
	}
	if DisplayHidden {
		fileArray = append(fileArray, ".", "..")
	}

	//* sort the arrays and proceed
	if !Timesort {
		if ReverseOrder {
			RevBubbleSort(fileArray)
			RevBubbleSort(Directories)
		} else {
			BubbleSort(fileArray)
			BubbleSort(Directories)
		}
	} else {
		if ReverseOrder {
			fileArray = SortByCreationTime(filepath, fileArray, true)
			Directories = SortByCreationTime(filepath, Directories, true)
		} else {
			fileArray = SortByCreationTime(filepath, fileArray, false)
			Directories = SortByCreationTime(filepath, Directories, false)
		}
	}

	for _, v := range fileArray {
		todisplay := ""
		filestat, err := os.Stat(filepath + "/" + v)
		if err != nil {
			fmt.Println("FILEARRAY ERR")
			log.Fatal(err)
		}
		permissions, err := GetFilePermissions(filepath + "/" + v)
		if permissions == "" {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal(err)
		}
		if !LongFormat && !DashO {
			if filestat.IsDir() || permissions == "rwx-rwx-r-x" {
				todisplay = BlueFormat(v)
				fmt.Print(todisplay + " ")
			} else {
				fmt.Print(filestat.Name() + " ")
			}
		} else if LongFormat || DashO {
			LongFormatDisplay(filepath + "/" + v)
		}
	}
	fmt.Println()
	for _, dir := range Directories {
		OrangePrintln(dir)
		RecursiveSearchDir(filepath + "/" + dir)
	}
}

func NormalSearchDir(filepath string) {
	var fileArray []string
	files, err := os.ReadDir(filepath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") && !DisplayHidden {
			continue
		}
		if err != nil {
			fmt.Println("STAT ERROR")
			log.Fatal(err)
		}
		fileArray = append(fileArray, file.Name())
	}
	if DisplayHidden {
		fileArray = append(fileArray, ".", "..")
	}

	//* sort the arrays and proceed
	if !Timesort {
		if ReverseOrder {
			RevBubbleSort(fileArray)
		} else {
			BubbleSort(fileArray)
		}
	} else {
		if ReverseOrder {
			fileArray = SortByCreationTime(filepath, fileArray, true)
		} else {
			fileArray = SortByCreationTime(filepath, fileArray, false)
		}
	}

	for _, v := range fileArray {
		todisplay := ""
		filestat, err := os.Stat(filepath + "/" + v)
		if err != nil {
			fmt.Println("FILEARRAY ERR")
			log.Fatal(err)
		}
		permissions, err := GetFilePermissions(filepath + "/" + v)
		if permissions == "" {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal(err)
		}
		if !LongFormat && !DashO {
			if filestat.IsDir() || permissions == "rwx-rwx-r-x" {
				todisplay = BlueFormat(v)
				fmt.Print(todisplay + " ")
			} else {
				fmt.Print(filestat.Name() + " ")
			}
		} else if LongFormat || DashO {
			LongFormatDisplay(filepath + "/" + v)
		}
	}
	fmt.Println()
}
