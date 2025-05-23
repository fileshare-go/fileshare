package upload

import (
	"context"
	"io"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunk"
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
	logrus.Debug("Starting Upload Process!")

	chunkList := make([]int32, 0)

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

		chunkList = append(chunkList, chunk.Index)

		logrus.Debug("chunk: ", chunk)
	}

	stream.SendAndClose(&pb.UploadStatus{
		Status: pb.Status_OK,
	})

	logrus.Debug("Ending Upload Process!")
	return nil
}
