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
			info, e := file.Info()
			if e != nil {
				log.Fatal("Lstat ERR: " + e.Error())
			}
			if info.Mode()&os.ModeSymlink != 0 {
				continue
			}
			if DisplayHidden {
				continue
			} else {
				fmt.Println("STAT ERROR")
				log.Fatal(err)
			}
		}
		if fileinfo.IsDir() {
			Directories = append(Directories, file.Name())
		}
		fileArray = append(fileArray, file.Name())
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
	maxLength := 0
	for _, file := range fileArray {
		length := len(file)
		if length > maxLength {
			maxLength = length
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
		if err != nil {
			log.Fatal(err)
		}
		// check formatting
		if !LongFormat && !DashO {
			if filestat.IsDir() || permissions[0] == 'd' {
				todisplay = BlueFormat(v)
				fmt.Print(todisplay + " ")
			} else {
				extension := getExtension(string(v))
				fmt.Println(extension)
				todisplay = getColorizedFileType(extension, string(v))
				fmt.Print(todisplay + " ")
			}
		} else if LongFormat || DashO {
			LongFormatDisplay(filepath + "/" + v)
		}
	}
	for _, dir := range Directories {
		OrangePrintln(dir)
		mainPath := filepath + "/" + dir
		RecursiveSearchDir(mainPath)
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
			info, e := file.Info()
			if e != nil {
				log.Fatal("Lstat ERR: " + e.Error())
			}
			if info.Mode()&os.ModeSymlink != 0 {
				continue
			}
			fmt.Println("STAT ERROR")
			log.Fatal(err)
		}
		fileArray = append(fileArray, file.Name())
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

	maxLength := 0
	for _, file := range fileArray {
		length := len(file)
		if length > maxLength {
			maxLength = length
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
			// padding := maxLength - len(string(v)) + 4
			if filestat.IsDir() || permissions == "rwx-rwx-r-x" {
				todisplay = BlueFormat(v)
				fmt.Print(todisplay + " ")
			} else {
				extension := getExtension(string(v))
				todisplay = getColorizedFileType(extension, string(v))
				fmt.Print(todisplay + " ")
			}
		} else if LongFormat || DashO {
			LongFormatDisplay(filepath + "/" + v)
		}
	}
	fmt.Println()
}
