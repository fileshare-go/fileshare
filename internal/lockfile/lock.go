package lockfile

import (
	"encoding/json"
	"fmt"
	"os"
)

const LOCK_FILE_NAME = "lock.json"

func GetLockPath(lockDirectory string) string {
	return fmt.Sprintf("%s/%s", lockDirectory, LOCK_FILE_NAME)
}

type LockFile struct {
	LockPath    string  `json:"lockPath"`
	FileName    string  `json:"filename"`
	Sha256      string  `json:"sha256"`
	ChunkNumber int     `json:"chunkNumber"`
	ChunkList   []int32 `json:"chunks"`
}

func ReadLockFile(lockDirectory string) (*LockFile, error) {
	lockPath := fmt.Sprintf("%s/%s", lockDirectory, LOCK_FILE_NAME)
	bytes, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, err
	}

	var lock LockFile
	if err := json.Unmarshal(bytes, &lock); err != nil {
		return nil, err
	}

	lock.LockPath = lockPath
	return &lock, nil
}

func (l *LockFile) SaveLock(lockDirectory string) error {
	bytes, err := json.Marshal(l)
	if err != nil {
		return err
	}
	file, err := os.Create(l.LockPath)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	return err
}

func (l *LockFile) UpdateLock(other *LockFile) {
	l.ChunkList = mergeList(l.ChunkList, other.ChunkList)
	l.ChunkNumber = len(l.ChunkList)
}

func mergeList(list1, list2 []int32) []int32 {
	result := []int32{}
	i, j := 0, 0

	for i < len(list1) && j < len(list2) {
		val := int32(0)
		if list1[i] < list2[j] {
			val = list1[i]
			i++
		} else if list1[i] > list2[j] {
			val = list2[j]
			j++
		} else {
			val = list1[i]
			i++
			j++
		}
		if len(result) == 0 || result[len(result)-1] != val {
			result = append(result, val)
		}
	}

	for i < len(list1) {
		if len(result) == 0 || result[len(result)-1] != list1[i] {
			result = append(result, list1[i])
		}
		i++
	}

	for j < len(list2) {
		if len(result) == 0 || result[len(result)-1] != list2[j] {
			result = append(result, list2[j])
		}
		j++
	}

	return result
}
