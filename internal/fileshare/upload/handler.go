package upload

import (
	"io"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunker"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	"github.com/chanmaoganda/fileshare/internal/lockfile"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	stream           pb.UploadService_UploadServer
	once             sync.Once
	chunkList        []int32
	meta             *pb.FileMeta
	totalChunkNumber int32
}

func NewHandler(stream pb.UploadService_UploadServer) *Handler {
	return &Handler{
		stream: stream,
		meta: &pb.FileMeta{},
	}
}

func (h *Handler) recordInformation(chunk *pb.FileChunk) {
	h.meta.Filename = chunk.Meta.Filename
	h.meta.Sha256 = chunk.Meta.Sha256

	h.totalChunkNumber = chunk.Total
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

		logrus.Debugf("filename: %s, total chunk: %d, chunk index: %d, chunk size: %d", chunk.Meta.Filename, chunk.GetTotal(), chunk.GetIndex(), len(chunk.GetData()))

		// create folder, record total chunk number and meta info
		h.once.Do(func() {
			h.recordInformation(chunk)
		})

		h.chunkList = append(h.chunkList, chunk.Index)

		if err := chunker.SaveChunk(chunk); err != nil {
			return err
		}
	}
	return h.SaveLockFile()
}

// close the stream, saving current status to lockfile
func (h *Handler) CloseWithErr(err error) error {
	logrus.Error("[handler] close with err: ", err)

	if err := h.SaveLockFile(); err != nil {
		logrus.Error(err)
	}

	return h.stream.SendAndClose(&pb.UploadStatus{
		Status:    pb.Status_ERROR,
		Meta:      h.meta,
		ChunkList: h.chunkList,
	})
}

// update lockfile if exists, else saves the lockfile
func (h *Handler) SaveLockFile() error {
	lockDirectory := h.meta.Sha256

	lockPath := lockfile.GetLockPath(lockDirectory)
	lock := lockfile.LockFile{
		FileName:         h.meta.Filename,
		Sha256:           lockDirectory,
		ChunkList:        h.chunkList,
		TotalChunkNumber: h.totalChunkNumber,
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
	if chunker.ValidateChunks(h.meta.Filename, h.meta.Sha256) {
		logrus.Debugf("[validate] %s validated! sha256 is %s", h.meta.Filename, h.meta.Sha256)
	} else {
		status = pb.Status_ERROR
		logrus.Warnf("[validate] %s not validated!", h.meta.Filename)
	}

	uploadStatus := pb.UploadStatus{
		Meta:      h.meta,
		Status:    status,
		ChunkList: h.chunkList,
	}
	if err := h.stream.SendAndClose(&uploadStatus); err != nil {
		logrus.Error("[validate] err: ", err)
	}
}
