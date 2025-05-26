package download

import (
	"fmt"
	"io"
	"os"
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

type StreamHandler struct {
	stream    pb.DownloadService_DownloadClient
	once      sync.Once
	Settings  *config.Settings
	Manager   *dbmanager.DBManager
	fileInfo  model.FileInfo
	chunkList []int32
}

func NewHandler(settings *config.Settings, stream pb.DownloadService_DownloadClient, manager *dbmanager.DBManager) *StreamHandler {
	return &StreamHandler{
		Settings: settings,
		stream:   stream,
		Manager:  manager,
	}
}

func (h *StreamHandler) Recv() error {
	for {
		chunk, err := h.stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		h.saveChunkToDisk(chunk)
	}

	h.fileInfo.UpdateChunks(h.chunkList)
	return nil
}

func (h *StreamHandler) saveChunkToDisk(chunk *pb.FileChunk) bool {
	debugprint.DebugChunk(chunk)

	h.recordFileInfoAndCreateCacheDir(chunk)

	// we need to handle if chunk has no data actually
	// or the situation that, task does not require any chunk
	// but for recording meta, send a chunk without actual data
	if len(chunk.Data) == 0 {
		logrus.Debugf("[Download] This chunk [%s] is empty, maybe it is for send file meta instead", chunk.Sha256[:8])
		return false
	}

	h.chunkList = append(h.chunkList, chunk.ChunkIndex)

	if err := chunkio.SaveChunk(h.Settings.CacheDirectory, chunk); err != nil {
		logrus.Error(err)
		return false
	}

	return true
}

// record chunk info for the first chunk
func (h *StreamHandler) recordFileInfoAndCreateCacheDir(chunk *pb.FileChunk) {
	h.once.Do(func() {
		// select from database
		h.fileInfo.Sha256 = chunk.Sha256
		h.Manager.SelectFileInfo(&h.fileInfo)

		folder := fmt.Sprintf("%s/%s", h.Settings.CacheDirectory, chunk.Sha256)
		if !fileutil.FileExists(folder) {
			if err := os.Mkdir(folder, 0755); err != nil {
				logrus.Error(err)
			}
		}
	})
}

// close the stream, saving current status to lockfile
func (h *StreamHandler) CloseWithErr(err error) error {
	logrus.Error("[handler] close with err: ", err)

	if !h.Manager.UpdateFileInfo(&h.fileInfo) {
		logrus.Warn("FileInfo save failed")
	}

	return h.stream.CloseSend()
}

func (h *StreamHandler) ValidateAndClose() {
	if h.fileInfo.ValidateChunks(h.Settings.CacheDirectory, h.Settings.DownloadDirectory) {
		logrus.Debugf("[Validate] %s validated! sha256 is %s", h.fileInfo.Filename, h.fileInfo.Sha256)
	} else {
		logrus.Warnf("[validate] %s not validated!", h.fileInfo.Filename)
	}

	h.Manager.UpdateFileInfo(&h.fileInfo)

	if err := h.stream.CloseSend(); err != nil {
		logrus.Error(err)
	}
}
