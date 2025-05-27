package recv

import (
	"io"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkstream"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
)

type ClientRecvStream struct {
	chunkstream.Core
	Stream pb.DownloadService_DownloadClient
}

func NewClientRecvStream(settings *config.Settings, manager *dbmanager.DBManager, stream pb.DownloadService_DownloadClient) chunkstream.StreamRecvCore {
	return &ClientRecvStream{
		Core: chunkstream.Core{
			Settings: settings,
			Manager:  manager,
		},
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

	// if err is not EOF, then this err should be handled
	if err != io.EOF {
		return err
	}

	// merge missing chunks that has been uploaded
	c.FileInfo.UpdateChunks(c.ChunkList)
	return nil
}

func (c *ClientRecvStream) RecvChunk() (*pb.FileChunk, error) {
	return c.Stream.Recv()
}

func (c *ClientRecvStream) ValidateRecvChunks() bool {
	return c.Validate()
}

func (c *ClientRecvStream) CloseStream(bool) error {
	c.Manager.UpdateFileInfo(&c.FileInfo)

	return c.Stream.CloseSend()
}
