package download

import (
	"io"
	"os"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/fileshare/debugprint"
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	stream    pb.DownloadService_DownloadClient
	once      sync.Once
	Manager   *dbmanager.DBManager
	fileInfo  model.FileInfo
	chunkList []int32
}

func NewHandler(stream pb.DownloadService_DownloadClient, manager *dbmanager.DBManager) *Handler {
	return &Handler{
		stream:  stream,
		Manager: manager,
	}
}

func (h *Handler) Recv() error {
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

func (h *Handler) saveChunkToDisk(chunk *pb.FileChunk) {
	debugprint.DebugChunk(chunk)

	h.once.Do(func() {
		// select from database
		h.fileInfo.Sha256 = chunk.Sha256
		h.Manager.SelectFileInfo(&h.fileInfo)

		if !fileutil.FileExists(chunk.Sha256) {
			if err := os.Mkdir(chunk.Sha256, 0755); err != nil {
				logrus.Error(err)
			}
		}
	})

	h.chunkList = append(h.chunkList, chunk.ChunkIndex)

	if err := chunkio.SaveChunk(chunk); err != nil {
		logrus.Error(err)
	}
}

// close the stream, saving current status to lockfile
func (h *Handler) CloseWithErr(err error) error {
	logrus.Error("[handler] close with err: ", err)

	h.Manager.UpdateFileInfo(&h.fileInfo)

	return h.stream.CloseSend()
}

func (h *Handler) ValidateAndClose() {
	if h.fileInfo.ValidateChunks() {
		logrus.Debugf("[Validate] %s validated! sha256 is %s", h.fileInfo.Filename, h.fileInfo.Sha256)
	} else {
		logrus.Warnf("[validate] %s not validated!", h.fileInfo.Filename)
	}

	h.Manager.UpdateFileInfo(&h.fileInfo)

	if err := h.stream.CloseSend(); err != nil {
		logrus.Error(err)
	}
}
