package upload

import (
	"io"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	stream    pb.UploadService_UploadServer
	once      sync.Once
	db        *gorm.DB
	fileInfo  *model.FileInfo
	chunkList []int32
}

func NewHandler(stream pb.UploadService_UploadServer) *Handler {
	return &Handler{
		stream:    stream,
		chunkList: []int32{},
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

	return nil
}

func (h *Handler) saveChunkToDisk(chunk *pb.FileChunk) {
	logrus.Debugf("file sha256: %s, chunk index: %d, chunk size: %d", chunk.Sha256, chunk.ChunkIndex, len(chunk.GetData()))

	h.once.Do(func() {
		// select from database
		h.db.First(&h.fileInfo, chunk.Sha256)
	})

	h.chunkList = append(h.chunkList, chunk.ChunkIndex)

	if err := chunkio.SaveChunk(chunk); err != nil {
		logrus.Error(err)
	}
}

// close the stream, saving current status to lockfile
func (h *Handler) CloseWithErr(err error) error {
	logrus.Error("[handler] close with err: ", err)

	if err := h.db.Model(&h.fileInfo).Updates(&h.fileInfo); err != nil {
		logrus.Error(err)
	}

	return h.stream.SendAndClose(&pb.UploadStatus{
		Status: pb.Status_ERROR,
		Meta: &pb.FileMeta{
			Filename: h.fileInfo.Filename,
			Sha256:   h.fileInfo.Sha256,
			FileSize: h.fileInfo.FileSize,
		},
		ChunkList: h.chunkList,
	})
}

func (h *Handler) ValidateAndClose() {
	status := pb.Status_OK
	if h.fileInfo.ValidateChunks() {
		logrus.Debugf("[validate] %s validated! sha256 is %s", h.fileInfo.Filename, h.fileInfo.Sha256)
	} else {
		status = pb.Status_ERROR
		logrus.Warnf("[validate] %s not validated!", h.fileInfo.Filename)
	}

	uploadStatus := pb.UploadStatus{
		Meta: &pb.FileMeta{
			Filename: h.fileInfo.Filename,
			Sha256:   h.fileInfo.Sha256,
			FileSize: h.fileInfo.FileSize,
		},
		Status:    status,
		ChunkList: h.chunkList,
	}
	if err := h.stream.SendAndClose(&uploadStatus); err != nil {
		logrus.Error("[validate] err: ", err)
	}

	h.db.Model(&h.fileInfo).Updates(&h.fileInfo)
}
