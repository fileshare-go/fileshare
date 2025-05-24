package download

import (
	"io"
	"os"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	stream    pb.DownloadService_DownloadClient
	once      sync.Once
	DB        *gorm.DB
	fileInfo  model.FileInfo
	chunkList []int32
}

func NewHandler(stream pb.DownloadService_DownloadClient, db *gorm.DB) *Handler {
	return &Handler{
		stream: stream,
		DB:     db,
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
	logrus.Debugf("file sha256: %s, chunk index: %d, chunk size: %d", chunk.Sha256, chunk.ChunkIndex, len(chunk.GetData()))

	h.once.Do(func() {
		// select from database
		h.DB.Where("sha256 = ?", chunk.Sha256).First(&h.fileInfo)
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

func (h *Handler) closeStreamAndSaveInfo() error {
	h.DB.Save(h.fileInfo)

	return h.stream.CloseSend()
}

// close the stream, saving current status to lockfile
func (h *Handler) CloseWithErr(err error) error {
	logrus.Error("[handler] close with err: ", err)

	if err := h.DB.Save(h.fileInfo); err != nil {
		logrus.Error(err)
	}

	return h.closeStreamAndSaveInfo()
}

func (h *Handler) ValidateAndClose() {
	if h.fileInfo.ValidateChunks() {
		logrus.Debugf("[validate] %s validated! sha256 is %s", h.fileInfo.Filename, h.fileInfo.Sha256)
	} else {
		logrus.Warnf("[validate] %s not validated!", h.fileInfo.Filename)
	}

	if err := h.closeStreamAndSaveInfo(); err != nil {
		logrus.Error(err)
	}
}
