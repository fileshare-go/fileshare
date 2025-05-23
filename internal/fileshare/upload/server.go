package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunk"
	pb "github.com/chanmaoganda/fileshare/proto/upload"
	"github.com/sirupsen/logrus"
)

type UploadServer struct {
	pb.UnimplementedUploadServiceServer
}

func (s *UploadServer) PreUpload(_ context.Context, task *pb.UploadTask) (*pb.UploadSummary, error) {
	logrus.Debugf("Upload task [filename: %s, file size: %d, sha256: %s]", task.Meta.Filename, task.FileSize, task.Meta.Sha256)

	chunkSummary := chunk.DealChunkSize(task.FileSize)

	chunkList := make([]int32, 0)
	for index := range chunkSummary.Number {
		chunkList = append(chunkList, index)
	}

	logrus.Debugf("Chunk Summary [chunk number: %d, chunk size: %d]", chunkSummary.Number, chunkSummary.Size)

	return &pb.UploadSummary{
		Meta: task.Meta,
		ChunkNumber: chunkSummary.Number,
		ChunkSize:   chunkSummary.Size,
		ChunkList: chunkList,
	}, nil
}

func (UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	logrus.Debug("Starting Upload Process!")

	chunkList := make([]int32, 0)
	once := sync.Once{}
	
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Error(err)
			return stream.SendAndClose(&pb.UploadStatus{
				Status: pb.Status_ERROR,
			})
		}

		logrus.Debugf("filename: %s, chunk index: %d, chunk size: %d", chunk.Meta.Filename, chunk.GetIndex(), len(chunk.GetData()))

		once.Do(func () {
			logrus.Debug("Creating directory for ", chunk.Meta.Sha256)
			if err := os.Mkdir(chunk.Meta.Sha256, 0755); err != nil {
				logrus.Warn(err)
			}
		})

		chunkList = append(chunkList, chunk.Index)

		if err := SaveChunk(chunk); err != nil {
			logrus.Error(err)
			return stream.SendAndClose(&pb.UploadStatus{
				Status: pb.Status_ERROR,
			})
		}
	}

	stream.SendAndClose(&pb.UploadStatus{
		Status: pb.Status_OK,
	})

	logrus.Debug("Ending Upload Process!")
	return nil
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
