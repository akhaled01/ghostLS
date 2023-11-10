package ghostls

import (
	"fmt"
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
			fmt.Println("ERROR:", err)
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

// * parsing Binary permissions
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
		return "-" + fmt.Sprintf("%s-%s-%s", ownerPermissions, groupPermissions, otherPermissions), nil
	}
}

// * syscall to get hard link numbers
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
	// blockSize, err := GetFileSystemBlockSize(directoryPath)
	// if err != nil {
	// 	// Error getting the file system block size
	// 	return 0, fmt.Errorf("error getting file system block size for %s: %w", directoryPath, err)
	// }
	const blockSize = 512

	var totalSize int64
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		// Error reading the directory
		return 0, fmt.Errorf("error reading directory %s: %w", directoryPath, err)
	}

	for _, file := range files {
		fullPath := directoryPath
		if !strings.HasSuffix(directoryPath, "/") {
			fullPath += "/"
		}
		fullPath += file.Name()

		fileSize, err := GetFileSizeConsideringSymlink(fullPath)
		if err != nil {
			// Log the error and continue with the next file
			fmt.Printf("error getting size for %s: %v\n", fullPath, err)
			continue
		}

		// fmt.Printf("File: %s, Size: %d\n", fullPath, fileSize)

		// Round up the file size to the nearest block

		// Inside the loop
		fileBlocks := (fileSize + blockSize - 1) / blockSize
		// fmt.Printf("File: %s, Size: %d, Blocks: %d\n", fullPath, fileSize, fileBlocks)
		totalSize += fileBlocks
	}

	// The block count is the total size divided by the block size
	blockCount := totalSize
	return blockCount, nil
}

func GetFileSystemBlockSize(path string) (int64, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return 0, err
	}
	return int64(stat.Bsize), nil
}

func GetFileSizeConsideringSymlink(path string) (int64, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		// It's a symlink, attempt to get the file it points to
		resolvedPath, err := os.Readlink(path)
		if err != nil {
			return 0, err
		}

		// If resolvedPath is not absolute, construct the full path manually
		if !strings.HasPrefix(resolvedPath, "/") {
			dir := path[:strings.LastIndex(path, "/")+1]
			resolvedPath = dir + resolvedPath
		}

		resolvedFileInfo, err := os.Stat(resolvedPath)
		if err != nil {
			// If the target of the symlink does not exist, handle accordingly
			// Option 1: Return the size of the symlink itself
			return fileInfo.Size(), nil
			// Option 2: Return zero size
			// return 0, nil
		}
		return resolvedFileInfo.Size(), nil
	}

	return fileInfo.Size(), nil
}
