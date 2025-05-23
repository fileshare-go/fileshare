package upload

import (
	"context"
	"fmt"
	"io"

	"github.com/chanmaoganda/fileshare/pkg/fileshare/chunk"
	pb "github.com/chanmaoganda/fileshare/proto/upload"
	"github.com/sirupsen/logrus"
)

type UploadServer struct {
	pb.UnimplementedUploadServiceServer
}

func (s *UploadServer) PreUpload(_ context.Context, task *pb.UploadTask) (*pb.UploadSummary, error) {
	logrus.Debugf("Upload task [filename: %s, file size: %d, sha256: %s]", task.Filename, task.FileSize, task.Sha256)
	chunkSummary := chunk.DealChunkSize(task.FileSize)
	logrus.Debugf("Chunk Summary [chunk size: %d, chunk number: %d]", chunkSummary.Size, chunkSummary.Number)
	return &pb.UploadSummary{
		Filename:    task.Filename,
		Sha256:      task.Sha256,
		ChunkNumber: chunkSummary.Number,
		ChunkSize:   chunkSummary.Size,
	}, nil
}

func (UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	chunkList := make([]int32, 0)
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return stream.SendAndClose(&pb.UploadStatus{
				Filename:  chunk.Filename,
				Sha256:    "",
				Status:    pb.Status_ERROR,
				ChunkList: chunkList,
			})
		}

		chunkList = append(chunkList, chunk.Index)

		fmt.Printf("chunk: %v\n", chunk)
	}

	return nil
}
