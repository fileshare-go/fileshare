package chunker

import (
	"fmt"
	"io"
	"os"

	"github.com/chanmaoganda/fileshare/internal/lockfile"
	"github.com/chanmaoganda/fileshare/internal/sha256"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
)

func MakeChunk(file *os.File, fileName, sha256 string, chunkSize int64, totalChunkNumber, chunkIndex int32) *pb.FileChunk {
	data := make([]byte, chunkSize)
	file.Seek(chunkSize*int64(chunkIndex), 0)
	n, err := file.Read(data)
	if err != nil {
		logrus.Error(err)
	}

	return &pb.FileChunk{
		Meta: &pb.FileMeta{
			Filename: fileName,
			Sha256:   sha256,
		},
		Total: totalChunkNumber,
		Index: chunkIndex,
		Data:  data[:n],
	}
}

func SaveChunk(chunk *pb.FileChunk) error {
	// Create or truncate the file
	chunkFileName := fmt.Sprintf("%s/%d", chunk.Meta.Sha256, chunk.Index)
	file, err := os.Create(chunkFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write bytes to the file
	_, err = file.Write(chunk.Data)
	if err != nil {
		return err
	}

	return nil
}

func ValidateChunks(filename, Sha256 string) bool {
	lock, err := lockfile.ReadLockFile(Sha256)
	if err != nil {
		logrus.Error("[validate]", err)
		return false
	}

	filePath := fmt.Sprintf("%s/%s", Sha256, filename)
	out, err := os.Create(filePath)
	if err != nil {
		logrus.Error("[validate]", err)
		return false
	}

	for _, index := range lock.ChunkList {
		in, err := os.Open(fmt.Sprintf("%s/%d", Sha256, index))
		if err != nil {
			logrus.Error("[validate]", err)
			return false
		}

		// Copy file contents to output file
		_, err = io.Copy(out, in)
		if err != nil {
			logrus.Error("[validate]", err)
			return false
		}
		in.Close()
	}
	out.Close()

	checkSum, err := sha256.CalculateSHA256(filePath)
	if err != nil {
		logrus.Error("[validate]", err)
		return false
	}

	return checkSum == Sha256
}
