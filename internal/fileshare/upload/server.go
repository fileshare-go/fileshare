package upload

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunker"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	"github.com/chanmaoganda/fileshare/internal/lockfile"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
)

type UploadServer struct {
	pb.UnimplementedUploadServiceServer
	Settings *config.Settings
}

// pre upload receives a task from client, calculate missing chunks and send the summary back
func (s *UploadServer) PreUpload(_ context.Context, task *pb.UploadTask) (*pb.UploadSummary, error) {
	logrus.Debugf("Upload task [filename: %s, file size: %d, sha256: %s]", task.Meta.Filename, task.FileSize, task.Meta.Sha256)

	chunkSummary := chunker.DealChunkSize(task.FileSize)

	chunkList := make([]int32, 0)
	for index := range chunkSummary.Number {
		chunkList = append(chunkList, index)
	}

	required := getMissingChunks(task.Meta.Sha256, chunkSummary.Number)

	logrus.Debugf("Chunk Summary [chunk number: %d, chunk size: %d]", chunkSummary.Number, chunkSummary.Size)

	return &pb.UploadSummary{
		Meta:        task.Meta,
		ChunkNumber: chunkSummary.Number,
		ChunkSize:   chunkSummary.Size,
		ChunkList:   required,
	}, nil
}

// upload receives chunks from client, save lockfile
func (s *UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	logrus.Debug("Starting Upload Process!")

	chunkList := make([]int32, 0)
	once := sync.Once{}
	var meta pb.FileMeta
	var totalChunkNumber int32

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return CloseWithErr(stream, &meta, totalChunkNumber, chunkList, err)
		}

		logrus.Debugf("filename: %s, total chunk: %d, chunk index: %d, chunk size: %d", chunk.Meta.Filename, chunk.GetTotal(), chunk.GetIndex(), len(chunk.GetData()))

		once.Do(func() {
			// create folder, record total chunk number and meta info
			initUpload(chunk, &meta, &totalChunkNumber)
		})

		chunkList = append(chunkList, chunk.Index)

		if err := chunker.SaveChunk(chunk); err != nil {
			return CloseWithErr(stream, &meta, totalChunkNumber, chunkList, err)
		}
	}

	if err := saveLockFile(meta.Sha256, &meta, chunkList, totalChunkNumber); err != nil {
		logrus.Error(err)
	}

	validity := chunker.ValidateChunks(meta.Filename, meta.Sha256)
	status := pb.Status_OK
	if validity {
		logrus.Debugf("[validate] %s validated! sha256 is %s", meta.Filename, meta.Sha256)
	} else {
		status = pb.Status_ERROR
		logrus.Warnf("[validate] %s not validated!", meta.Filename)
	}

	uploadStatus := pb.UploadStatus{
		Meta:      &meta,
		Status:    status,
		ChunkList: chunkList,
	}
	stream.SendAndClose(&uploadStatus)

	logrus.Debug("Ending Upload Process!")
	return nil
}

// close the stream, saving current status to lockfile
func CloseWithErr(stream pb.UploadService_UploadServer, meta *pb.FileMeta, totalChunkNumber int32, chunkList []int32, err error) error {
	logrus.Error(err)

	if err := saveLockFile(meta.Sha256, meta, chunkList, totalChunkNumber); err != nil {
		logrus.Error(err)
	}

	return stream.SendAndClose(&pb.UploadStatus{
		Status:    pb.Status_ERROR,
		Meta:      meta,
		ChunkList: chunkList,
	})
}

// get missing chunks from lockfile if exists, either return the enum of `total`
func getMissingChunks(lockFolder string, total int32) []int32 {
	lockPath := lockfile.GetLockPath(lockFolder)
	if fileutil.FileExists(lockPath) {
		lock, err := lockfile.ReadLockFile(lockFolder)
		if err != nil {
			logrus.Error(err)
			return []int32{}
		}

		return lock.RemainingChunks()
	}
	result := []int32{}
	for i := range total {
		result = append(result, i)
	}
	return result
}

// store meta and totalChunkNumber
func initUpload(chunk *pb.FileChunk, meta *pb.FileMeta, totalChunkNumber *int32) {
	meta.Filename = chunk.Meta.Filename
	meta.Sha256 = chunk.Meta.Sha256

	*totalChunkNumber = chunk.Total

	dirName := chunk.Meta.Sha256

	logrus.Debug("Creating directory for ", dirName)

	if fileutil.FileExists(dirName) {
		return
	}

	if err := os.Mkdir(dirName, 0755); err != nil {
		logrus.Errorf("While creating %s, %s", dirName, err.Error())
	}
}

// update lockfile if exists, else saves the lockfile
func saveLockFile(lockDirectory string, meta *pb.FileMeta, chunkList []int32, totalChunkNumber int32) error {
	lockPath := lockfile.GetLockPath(lockDirectory)
	lock := lockfile.LockFile{
		FileName:         meta.Filename,
		Sha256:           meta.Sha256,
		ChunkList:        chunkList,
		TotalChunkNumber: totalChunkNumber,
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
