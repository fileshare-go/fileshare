package chunkstream

import (
	"os"
	"strings"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/chunkio"
	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/chanmaoganda/fileshare/internal/service"
	"github.com/sirupsen/logrus"
)

// core functions for recv functionality
type StreamRecvCore interface {
	// recv all chunks
	RecvStreamChunks() error
	// recv one chunk, used in RecvStreamChunks
	RecvChunk() (*pb.FileChunk, error)
	// validate all received chunks with checksum256
	ValidateRecvChunks() bool
	// close stream and save states
	CloseStream(bool) error
}

// core functions for send functionality
type StreamSendCore interface {
	// send all chunks
	SendStreamChunks() error
	// send one chunk, used in SendStreamChunks
	SendChunk(*pb.FileChunk) error
	// close stream and save states
	CloseStream() error
}

type Core struct {
	FileInfo  model.FileInfo
	Once      sync.Once
	ChunkList []int32
}

// the first time core recv a chunk, record the FileInfo and then mkdir for cache folder with checksum
func (c *Core) SetupAndRecordInfo(chunk *pb.FileChunk) {
	c.Once.Do(func() {
		// select from database
		c.FileInfo.Sha256 = chunk.Sha256
		service.Orm().Find(&c.FileInfo)

		// create sha256 folder in cache folder
		folder := strings.Join([]string{config.Cfg().CacheDirectory, chunk.Sha256}, "/")
		if !util.FileExists(folder) {
			if err := os.Mkdir(folder, 0755); err != nil {
				logrus.Error(err)
			}
		}
	})
}

func (c *Core) SaveChunkToDisk(chunk *pb.FileChunk) bool {
	util.DebugChunk(chunk)
	c.SetupAndRecordInfo(chunk)

	// we need to handle if chunk has no data actually
	// or the situation that, task does not require any chunk
	// but for recording meta, send a chunk without actual data
	if len(chunk.Data) == 0 {
		logrus.Debugf("[Upload] This chunk [%s] is empty, maybe it is for send file meta instead", chunk.Sha256[:8])
		return false
	}

	c.ChunkList = append(c.ChunkList, chunk.ChunkIndex)

	if err := chunkio.SaveChunk(config.Cfg().CacheDirectory, chunk); err != nil {
		logrus.Error(err)
		return false
	}

	return true
}

// call FileInfo to validate chunks within chunklist
func (c *Core) Validate() bool {
	if c.FileInfo.ValidateChunks(config.Cfg().CacheDirectory, config.Cfg().DownloadDirectory) {
		logrus.Debugf("[Validate] %s validated! sha256 is %s", util.Render(c.FileInfo.Filename), util.Render(c.FileInfo.Sha256))
		return true
	}

	logrus.Warnf("[validate] %s not validated! sha256 is %s", util.Render(c.FileInfo.Filename), util.Render(c.FileInfo.Sha256))
	return false
}
