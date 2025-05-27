package send

import (
	"os"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
)

type ClientSendStream struct {
	Stream pb.UploadService_UploadClient
	Task   *pb.UploadTask
	File   *os.File
}

func NewClientSendStream(task *pb.UploadTask, filePath string, stream pb.UploadService_UploadClient) *ClientSendStream {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Panic(err)
	}

	return &ClientSendStream{
		Stream: stream,
		Task:   task,
		File:   file,
	}
}

func (s *ClientSendStream) SendStreamChunks() error {
	if len(s.Task.ChunkList) == 0 {
		return s.SendChunk(s.LoadEmptyChunk())
	}

	for _, idx := range s.Task.ChunkList {
		if err := s.SendChunk(s.LoadChunk(idx)); err != nil {
			return err
		}
	}
	return nil
}

func (s *ClientSendStream) SendChunk(chunk *pb.FileChunk) error {
	return s.Stream.Send(chunk)
}

func (s *ClientSendStream) LoadChunk(chunkIdx int32) *pb.FileChunk {
	return chunkio.MakeChunk(s.File, s.Task.Meta.Sha256, s.Task.ChunkSize, chunkIdx)
}

func (s *ClientSendStream) LoadEmptyChunk() *pb.FileChunk {
	return &pb.FileChunk{
		Sha256:     s.Task.Meta.Sha256,
		ChunkIndex: 0,
		Data:       []byte{},
	}
}
