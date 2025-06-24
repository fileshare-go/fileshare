package recv

import (
	"io"

	"github.com/chanmaoganda/fileshare/internal/core/chunkstream"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/chanmaoganda/fileshare/internal/service"
)

type ClientRecvStream struct {
	chunkstream.Core
	Stream pb.DownloadService_DownloadClient
}

func NewClientRecvStream(stream pb.DownloadService_DownloadClient) chunkstream.StreamRecvCore {
	return &ClientRecvStream{
		Core:   chunkstream.Core{},
		Stream: stream,
	}
}

func (c *ClientRecvStream) RecvStreamChunks() error {
	var chunk *pb.FileChunk
	var err error

	for chunk, err = c.RecvChunk(); err == nil; chunk, err = c.RecvChunk() {
		if saveStatus := c.SaveChunkToDisk(chunk); !saveStatus {
			break
		}
	}

	// merge missing chunks that has been uploaded
	// update current chunks whether err is nil or not
	c.FileInfo.UpdateChunks(c.ChunkList)

	// if err is not EOF, then this err should be handled
	if err != io.EOF {
		return err
	}

	return nil
}

func (c *ClientRecvStream) RecvChunk() (*pb.FileChunk, error) {
	return c.Stream.Recv()
}

func (c *ClientRecvStream) ValidateRecvChunks() bool {
	return c.Validate()
}

func (c *ClientRecvStream) CloseStream(bool) error {
	var err error
	if err = service.Orm().Save(&c.FileInfo).Error; err != nil {
		return err
	}

	return c.Stream.CloseSend()
}
