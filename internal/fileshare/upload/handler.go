package upload

import (
	"io"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunker"
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	"github.com/chanmaoganda/fileshare/internal/lockfile"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Handler struct {
	stream    pb.UploadService_UploadServer
	once      sync.Once
	db        *gorm.DB
	fileInfo  *model.File
	chunkList []int32
}

func NewHandler(stream pb.UploadService_UploadServer) *Handler {
	return &Handler{
		stream:    stream,
		chunkList: []int32{},
	}
}

func (h *Handler) RecvAndSaveLock() error {
	for {
		chunk, err := h.stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		h.saveChunk(chunk)
	}
	return h.SaveLockFile()
}

func (h *Handler) saveChunk(chunk *pb.FileChunk) {
	logrus.Debugf("file sha256: %s, chunk index: %d, chunk size: %d", chunk.Sha256, chunk.ChunkIndex, len(chunk.GetData()))

	h.once.Do(func() {
		h.db.First(&h.fileInfo, chunk.Sha256)	
	})

	h.chunkList = append(h.chunkList, chunk.ChunkIndex)

	if err := chunker.SaveChunk(chunk); err != nil {
		logrus.Error(err)
	}
}

// close the stream, saving current status to lockfile
func (h *Handler) CloseWithErr(err error) error {
	logrus.Error("[handler] close with err: ", err)

	if err := h.SaveLockFile(); err != nil {
		logrus.Error(err)
	}

	return h.stream.SendAndClose(&pb.UploadStatus{
		Status: pb.Status_ERROR,
		Meta: &pb.FileMeta{
			Filename: h.fileInfo.Filename,
			Sha256: h.fileInfo.Sha256,
			FileSize: h.fileInfo.FileSize,
		},
		ChunkList: h.chunkList,
	})
}

// update lockfile if exists, else saves the lockfile
func (h *Handler) SaveLockFile() error {
	lockDirectory := h.fileInfo.Sha256

	lockPath := lockfile.GetLockPath(lockDirectory)
	lock := lockfile.LockFile{
		FileName:         h.fileInfo.Filename,
		Sha256:           lockDirectory,
		ChunkSize:        h.fileInfo.ChunkSize,
		ChunkList:        h.chunkList,
		TotalChunkNumber: h.fileInfo.ChunkNumber,
	}

	if !fileutil.FileExists(lockPath) {
		return lock.SaveLock(lockDirectory)
	}

	oldLock, err := lockfile.ReadLockFile(lockDirectory)
	if err != nil {
		return err
	}

	lock.UpdateLock(oldLock)
	return lock.SaveLock(lockDirectory)
}

func (h *Handler) ValidateAndClose() {
	status := pb.Status_OK
	if chunker.ValidateChunks(h.fileInfo.Filename, h.fileInfo.Sha256) {
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
}
