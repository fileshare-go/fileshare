package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunk"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	"github.com/chanmaoganda/fileshare/internal/lockfile"
	pb "github.com/chanmaoganda/fileshare/proto/upload"
	"github.com/sirupsen/logrus"
)

type UploadServer struct {
	pb.UnimplementedUploadServiceServer
	Settings *config.Settings
}

func (s *UploadServer) PreUpload(_ context.Context, task *pb.UploadTask) (*pb.UploadSummary, error) {
	logrus.Debugf("Upload task [filename: %s, file size: %d, sha256: %s]", task.Meta.Filename, task.FileSize, task.Meta.Sha256)

	chunkSummary := chunk.DealChunkSize(task.FileSize)

	chunkList := make([]int32, 0)
	for index := range chunkSummary.Number {
		chunkList = append(chunkList, index)
	}

	logrus.Debugf("Chunk Summary [chunk number: %d, chunk size: %d]", chunkSummary.Number, chunkSummary.Size)

	return &pb.UploadSummary{
		Meta:        task.Meta,
		ChunkNumber: chunkSummary.Number,
		ChunkSize:   chunkSummary.Size,
		ChunkList:   chunkList,
	}, nil
}

func (s *UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	logrus.Debug("Starting Upload Process!")

	chunkList := make([]int32, 0)
	once := sync.Once{}
	var meta pb.FileMeta

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Error(err)
			return stream.SendAndClose(&pb.UploadStatus{
				Status: pb.Status_ERROR,
			})
		}

		logrus.Debugf("filename: %s, chunk index: %d, chunk size: %d", chunk.Meta.Filename, chunk.GetIndex(), len(chunk.GetData()))

		once.Do(func() {
			// create folder and record meta info
			initUpload(chunk, &meta)
		})

		chunkList = append(chunkList, chunk.Index)

		if err := SaveChunk(chunk); err != nil {
			logrus.Error(err)
			return stream.SendAndClose(&pb.UploadStatus{
				Status: pb.Status_ERROR,
			})
		}
	}

	uploadStatus := pb.UploadStatus{
		Meta:      &meta,
		Status:    pb.Status_OK,
		ChunkList: chunkList,
	}

	if err := saveLockFile(meta.Sha256, &uploadStatus); err != nil {
		logrus.Error(err)
	}

	stream.SendAndClose(&uploadStatus)

	logrus.Debug("Ending Upload Process!")
	return nil
}

func initUpload(chunk *pb.FileChunk, meta *pb.FileMeta) {
	dirName := chunk.Meta.Sha256

	meta.Filename = chunk.Meta.Filename
	meta.Sha256 = chunk.Meta.Sha256

	logrus.Debug("Creating directory for ", dirName)

	if fileutil.FileExists(dirName) {
		return
	}

	if err := os.Mkdir(dirName, 0755); err != nil {
		logrus.Errorf("While creating %s, %s", dirName, err.Error())
	}
}

func SaveChunk(chunk *pb.FileChunk) error {
	// Create or truncate the file
	chunkFileName := fmt.Sprintf("%s/%d", chunk.Meta.Sha256, chunk.Index)
	file, err := os.Create(chunkFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write bytes to the file
	_, err = file.Write(chunk.Data)
	if err != nil {
		return err
	}

	return nil
}

func saveLockFile(lockDirectory string, status *pb.UploadStatus) error {
	lockPath := lockfile.GetLockPath(lockDirectory)
	lock := lockfile.LockFile{
		LockPath:  lockPath,
		FileName:  status.Meta.Filename,
		Sha256:    status.Meta.Sha256,
		ChunkList: status.ChunkList,
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
