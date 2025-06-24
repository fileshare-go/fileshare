package send

import (
	"os"

	"github.com/chanmaoganda/fileshare/internal/core/chunkstream"
	"github.com/chanmaoganda/fileshare/internal/pkg/chunkio"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
)

type ClientSendStream struct {
	Stream       pb.UploadService_UploadClient
	Task         *pb.UploadTask
	SerialLoader *chunkio.SerialChunkLoader
}

func NewClientSendStream(task *pb.UploadTask, filePath string, stream pb.UploadService_UploadClient) chunkstream.StreamSendCore {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Fatal(err)
	}

	return &ClientSendStream{
		Stream:       stream,
		Task:         task,
		SerialLoader: chunkio.NewSerialChunkLoader(file, task.Meta.Sha256, task.ChunkSize),
	}
}

func (c *ClientSendStream) SendStreamChunks() error {
	if len(c.Task.ChunkList) == 0 {
		return c.SendChunk(c.LoadEmptyChunk())
	}

	for _, idx := range c.Task.ChunkList {
		if err := c.SendChunk(c.LoadChunk(idx)); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientSendStream) SendChunk(chunk *pb.FileChunk) error {
	return c.Stream.Send(chunk)
}

func (c *ClientSendStream) CloseStream() error {
	status, err := c.Stream.CloseAndRecv()
	if err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Debugf("[Upload] Status Info [status: %d]", status.Status)
	return nil
}

func (s *ClientSendStream) LoadChunk(chunkIdx int32) *pb.FileChunk {
	return s.SerialLoader.LoadChunk(chunkIdx)
}

func (s *ClientSendStream) LoadEmptyChunk() *pb.FileChunk {
	return &pb.FileChunk{
		Sha256:     s.Task.Meta.Sha256,
		ChunkIndex: 0,
		Data:       []byte{},
	}
}
