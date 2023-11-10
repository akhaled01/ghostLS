package ghostls

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type File struct {
	Path string
	Time time.Time
}

func SortByCreationTime(initialdir string, filePaths []string, reverse bool) []string {
	files := make([]File, len(filePaths))

	for i, path := range filePaths {
		fileInfo, err := os.Stat(initialdir + "/" + path)
		if err != nil {
			log.Fatal(err)
		}

		files[i] = File{
			Path: path,
			Time: fileInfo.ModTime(),
		}
	}

	for i := 0; i < len(files)-1; i++ {
		maxIndex := i

		for j := i + 1; j < len(files); j++ {
			if ReverseOrder {
				if files[j].Time.Before(files[maxIndex].Time) {
					maxIndex = j
				}
			} else {
				if files[j].Time.After(files[maxIndex].Time) {
					maxIndex = j
				}
			}
		}

		files[i], files[maxIndex] = files[maxIndex], files[i]
	}

	sortedPaths := make([]string, len(files))
	for i, file := range files {
		sortedPaths[i] = file.Path
	}
	return sortedPaths
}

func BubbleSort(arr []string) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if strings.ToLower(arr[i]) > strings.ToLower(arr[j]) {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

func RevBubbleSort(arr []string) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if strings.ToLower(arr[i]) < strings.ToLower(arr[j]) {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

//* parsing Binary permissions
func GetFilePermissions(path string) (string, error) {
	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	// Get permission bits
	mode := fileInfo.Mode()
	permissions := mode.Perm()

	// Convert permission bits to string
	permissionString := strconv.FormatUint(uint64(permissions), 8)

	// Pad the permission string to ensure 3 digits
	for len(permissionString) < 3 {
		permissionString = "0" + permissionString
	}

	// Map permission bits to their respective symbols
	permissionSymbols := map[int]string{
		0: "---",
		1: "--x",
		2: "-w-",
		3: "-wx",
		4: "r--",
		5: "r-x",
		6: "rw-",
		7: "rwx",
	}

	// Extract permission symbols for owner, group, and others
	ownerPermissions := permissionSymbols[int(permissions>>6)&7]
	groupPermissions := permissionSymbols[int(permissions>>3)&7]
	otherPermissions := permissionSymbols[int(permissions)&7]
	dirbool := fileInfo.IsDir()
	// Return the formatted permissions
	if dirbool {
		return "d" + fmt.Sprintf("%s-%s-%s", ownerPermissions, groupPermissions, otherPermissions), nil
	} else {
		return fmt.Sprintf("%s-%s-%s", ownerPermissions, groupPermissions, otherPermissions), nil
	}
}

//* syscall to get hard link numbers
func GetHardLinkNum(path string) (string, error) {
	fcount := uint16(0)

	fileinfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if sys := fileinfo.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			fcount = stat.Nlink
		}
	}
	mainnum := strconv.Itoa(int(fcount))

	return mainnum, nil
}

func GetFileOwnerAndGroup(filePath string) (string, string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", "", err
	}

	fileOwner := fileInfo.Sys().(*syscall.Stat_t).Uid
	fileGroup := fileInfo.Sys().(*syscall.Stat_t).Gid

	owner, err := lookupUserById(fileOwner)
	if err != nil {
		return "", "", err
	}

	group, err := lookupGroupById(fileGroup)
	if err != nil {
		return "", "", err
	}

	return owner, group, nil
}

func lookupUserById(uid uint32) (string, error) {
	u, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

func lookupGroupById(gid uint32) (string, error) {
	g, err := user.LookupGroupId(strconv.Itoa(int(gid)))
	if err != nil {
		return "", err
	}
	return g.Name, nil
}

func GetBlockCount(directoryPath string) (int64, error) {
	// Open the directory
	dir, err := os.Open(directoryPath)
	if err != nil {
		return 0, err
	}
	defer dir.Close()

	// Get the directory file information
	dirInfo, err := dir.Stat()
	if err != nil {
		return 0, err
	}

	// Get the underlying syscall.Stat_t structure
	stat := dirInfo.Sys().(*syscall.Stat_t)

	// Return the block count
	return stat.Blocks, nil
}
