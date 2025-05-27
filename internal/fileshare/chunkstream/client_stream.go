package chunkstream

import (
	"io"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
)

type ClientStream struct {
	Core
	Stream pb.DownloadService_DownloadClient
}

func (c *ClientStream) RecvStreamChunks() error {
	var chunk *pb.FileChunk
	var err error

	for chunk, err = c.RecvChunk(); err == nil; chunk, err = c.RecvChunk() {
		if saveStatus := c.SaveChunkToDisk(chunk); !saveStatus {
			break
		}
	}

	// if err is not EOF, then this err should be handled
	if err != io.EOF {
		return err
	}

	// merge missing chunks that has been uploaded
	c.FileInfo.UpdateChunks(c.ChunkList)
	return nil
}

func NewClientStream(settings *config.Settings, manager *dbmanager.DBManager, stream pb.DownloadService_DownloadClient) *ClientStream {
	return &ClientStream{
		Core: Core{
			Settings: settings,
			Manager: manager,
		},
		Stream: stream,
	}
}

func (c *ClientStream) RecvChunk() (*pb.FileChunk, error) {
	return c.Stream.Recv()
}

func (c *ClientStream) ValidateFile() bool {
	return c.Core.Validate()
}

func (c *ClientStream) CloseStream(bool) error {
	return c.Stream.CloseSend()
}
