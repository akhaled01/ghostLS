package ghostls

import (
    "os"
    "os/exec"
    "strconv"
    "strings"
    "testing"
)

func getExpectedBlockSize(path string) (int64, error) {
    // Using 'stat' command to get the block size
    // The command and its arguments might vary based on the operating system
    cmd := exec.Command("stat", "-f", "%k", path)
    output, err := cmd.Output()
    if err != nil {
        return 0, err
    }

    blockSizeStr := strings.TrimSpace(string(output))
    return strconv.ParseInt(blockSizeStr, 10, 64)
}

func TestGetFileSystemBlockSize(t *testing.T) {
    currentDir, err := os.Getwd()
    if err != nil {
        t.Fatalf("Failed to get current directory: %v", err)
    }

    expectedBlockSize, err := getExpectedBlockSize(currentDir)
    if err != nil {
        t.Fatalf("Failed to get expected block size: %v", err)
    }

    blockSize, err := GetFileSystemBlockSize(currentDir)
    if err != nil {
        t.Fatalf("GetFileSystemBlockSize returned an error: %v", err)
    }

    if blockSize != expectedBlockSize {
        t.Errorf("Expected block size %d, got %d", expectedBlockSize, blockSize)
    }
}
