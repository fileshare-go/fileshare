package upload

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

type StreamHandler struct {
	stream    pb.UploadService_UploadServer
	once      sync.Once
	DB        *gorm.DB
	fileInfo  model.FileInfo
	chunkList []int32
}

func NewHandler(stream pb.UploadService_UploadServer, db *gorm.DB) *StreamHandler {
	return &StreamHandler{
		stream:    stream,
		DB:        db,
		chunkList: []int32{},
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

func (h *StreamHandler) saveChunkToDisk(chunk *pb.FileChunk) {
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

func (h *StreamHandler) closeStreamAndSaveInfo(status pb.Status) error {
	uploadStatus := &pb.UploadStatus{
		Status: status,
		Meta: &pb.FileMeta{
			Filename: h.fileInfo.Filename,
			Sha256:   h.fileInfo.Sha256,
			FileSize: h.fileInfo.FileSize,
		},
		ChunkList: h.chunkList,
	}

	h.DB.Save(h.fileInfo)

	return h.stream.SendAndClose(uploadStatus)
}

// close the stream, saving current status to lockfile
func (h *StreamHandler) CloseWithErr(err error) error {
	logrus.Error("[handler] close with err: ", err)

	if err := h.DB.Save(h.fileInfo); err != nil {
		logrus.Error(err)
	}

	return h.closeStreamAndSaveInfo(pb.Status_ERROR)
}

func (h *StreamHandler) ValidateAndClose() {
	status := pb.Status_OK
	if h.fileInfo.ValidateChunks() {
		logrus.Debugf("[validate] %s validated! sha256 is %s", h.fileInfo.Filename, h.fileInfo.Sha256)
	} else {
		status = pb.Status_ERROR
		logrus.Warnf("[validate] %s not validated!", h.fileInfo.Filename)
	}

	if err := h.closeStreamAndSaveInfo(status); err != nil {
		logrus.Error(err)
	}
}
