package chunkio

import (
	"fmt"
	"os"

	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
)

func SaveChunk(cache_dir string, chunk *pb.FileChunk) error {
	// Create or truncate the file
	chunkFileName := fmt.Sprintf("%s/%s/%d", cache_dir, chunk.Sha256, chunk.ChunkIndex)
	file, err := os.Create(chunkFileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Warn(err)
		}
	}()

	// Write bytes to the file
	_, err = file.Write(chunk.Data)
	if err != nil {
		return err
	}

	return nil
}

func UploadChunk(cache_dir, sha256 string, chunkIndex int32) []byte {
	chunkFileName := fmt.Sprintf("%s/%s/%d", cache_dir, sha256, chunkIndex)
	bytes, err := os.ReadFile(chunkFileName)
	if err != nil {
		logrus.Error(err)
		return []byte{}
	}

	return bytes
}
