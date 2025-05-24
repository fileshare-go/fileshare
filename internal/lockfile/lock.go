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
	FileName         string  `json:"filename"`
	Sha256           string  `json:"sha256"`
	ChunkSize        int64   `json:"chunkSize"`
	TotalChunkNumber int32   `json:"totalChunkNumber"`
	ChunkList        []int32 `json:"chunks"`
}

func ReadLockFile(lockDirectory string) (*LockFile, error) {
	lockPath := GetLockPath(lockDirectory)
	bytes, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, err
	}

	var lock LockFile
	if err := json.Unmarshal(bytes, &lock); err != nil {
		return nil, err
	}

	return &lock, nil
}

func (l *LockFile) SaveLock(lockDirectory string) error {
	bytes, err := json.Marshal(l)
	if err != nil {
		return err
	}
	file, err := os.Create(GetLockPath(lockDirectory))
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	return err
}

func (l *LockFile) UpdateLock(other *LockFile) {
	l.ChunkList = mergeList(l.ChunkList, other.ChunkList)
}

func (l *LockFile) RemainingChunks() []int32 {
	if l.TotalChunkNumber == 0 {
		return []int32{}
	}

	return missingElementsInSortedList(l.TotalChunkNumber, l.ChunkList)
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

func missingElementsInSortedList(total int32, subList []int32) []int32 {
	if len(subList) == int(total) {
		return []int32{}
	}

	i := int32(0)
	j := 0
	result := []int32{}

	for i < total && j < len(subList) {
		if i == subList[j] {
			i++
			j++
		} else if i < subList[j] {
			result = append(result, i)
			i++
		} else {
			return []int32{}
		}
	}

	for i < total {
		result = append(result, i)
		i++
	}

	return result
}
