package chunkstream

import (
	"os"
	"strings"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/debugprint"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	"github.com/chanmaoganda/fileshare/internal/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
)

type StreamRecvCore interface {
	RecvStreamChunks() error
	RecvChunk() (*pb.FileChunk, error)
	CloseStream(bool) error
}

type StreamSendCore interface {
	SendStreamChunks() error
	SendChunk(*pb.FileChunk) error
}

// type TaskMeta interface {
// 	GetChunkList() []int32
// 	GetMeta() *pb.FileMeta
// 	GetChunkNumber() int32
// }

type Core struct {
	Settings  *config.Settings
	Manager   *dbmanager.DBManager
	FileInfo  model.FileInfo
	Once      sync.Once
	ChunkList []int32
}

func (c *Core) SetupAndRecordInfo(chunk *pb.FileChunk) {
	c.Once.Do(func() {
		// select from database
		c.FileInfo.Sha256 = chunk.Sha256
		c.Manager.SelectFileInfo(&c.FileInfo)

		// create sha256 folder in cache folder
		folder := strings.Join([]string{c.Settings.CacheDirectory, chunk.Sha256}, "/")
		if !fileutil.FileExists(folder) {
			if err := os.Mkdir(folder, 0755); err != nil {
				logrus.Error(err)
			}
		}
	})
}

func (c *Core) SaveChunkToDisk(chunk *pb.FileChunk) bool {
	debugprint.DebugChunk(chunk)
	c.SetupAndRecordInfo(chunk)

	// we need to handle if chunk has no data actually
	// or the situation that, task does not require any chunk
	// but for recording meta, send a chunk without actual data
	if len(chunk.Data) == 0 {
		logrus.Debugf("[Upload] This chunk [%s] is empty, maybe it is for send file meta instead", chunk.Sha256[:8])
		return false
	}

	c.ChunkList = append(c.ChunkList, chunk.ChunkIndex)

	if err := chunkio.SaveChunk(c.Settings.CacheDirectory, chunk); err != nil {
		logrus.Error(err)
		return false
	}

	return true
}

func (c *Core) Validate() bool {
	if c.FileInfo.ValidateChunks(c.Settings.CacheDirectory, c.Settings.DownloadDirectory) {
		logrus.Debugf("[Validate] %s validated! sha256 is %s", debugprint.Render(c.FileInfo.Filename), debugprint.Render(c.FileInfo.Sha256))
		return true
	}

	logrus.Warnf("[validate] %s not validated! sha256 is %s", debugprint.Render(c.FileInfo.Filename), debugprint.Render(c.FileInfo.Sha256))
	return false
}
