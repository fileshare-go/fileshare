package upload

import (
	"io"
	"os"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/debugprint"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	"github.com/chanmaoganda/fileshare/internal/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StreamHandler struct {
	stream    pb.UploadService_UploadServer
	once      sync.Once
	Manager   *dbmanager.DBManager
	fileInfo  model.FileInfo
	chunkList []int32
}

func NewHandler(stream pb.UploadService_UploadServer, DB *gorm.DB) *StreamHandler {
	return &StreamHandler{
		stream:    stream,
		Manager:   dbmanager.NewDBManager(DB),
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

		if !h.saveChunkToDisk(chunk) {
			break
		}
	}

	h.fileInfo.UpdateChunks(h.chunkList)
	return nil
}

func (h *StreamHandler) saveChunkToDisk(chunk *pb.FileChunk) bool {
	debugprint.DebugChunk(chunk)

	h.onceJob(chunk)

	// we need to handle if chunk has no data actually
	// or the situation that, task does not require any chunk
	// but for recording meta, send a chunk without actual data
	if len(chunk.Data) == 0 {
		logrus.Debugf("[Upload] This chunk [%s] is empty, maybe it is for send file meta instead", chunk.Sha256[:8])
		return false
	}

	h.chunkList = append(h.chunkList, chunk.ChunkIndex)

	if err := chunkio.SaveChunk(chunk); err != nil {
		logrus.Error(err)
		return false
	}

	return true
}

// record chunk info for the first chunk
func (h *StreamHandler) onceJob(chunk *pb.FileChunk) {
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

	h.Manager.UpdateFileInfo(&h.fileInfo)

	return h.stream.SendAndClose(uploadStatus)
}

// close the stream, saving current status to lockfile
func (h *StreamHandler) CloseWithErr(err error) error {
	logrus.Error("[handler] close with err: ", err)

	if !h.Manager.UpdateFileInfo(&h.fileInfo) {
		logrus.Warn("FileInfo save failed")
	}

	return h.closeStreamAndSaveInfo(pb.Status_ERROR)
}

func (h *StreamHandler) ValidateAndClose() {
	status := pb.Status_OK
	if h.fileInfo.ValidateChunks() {
		logrus.Debugf("[Validate] %s validated! sha256 is %s", debugprint.Render(h.fileInfo.Filename), debugprint.Render(h.fileInfo.Sha256))
	} else {
		status = pb.Status_ERROR
		logrus.Warnf("[validate] %s not validated!", debugprint.Render(h.fileInfo.Filename))
	}

	if err := h.closeStreamAndSaveInfo(status); err != nil {
		logrus.Error(err)
	}
}
