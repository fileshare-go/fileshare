package send

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkstream"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
)

type ServerSendStream struct {
	chunkstream.Core
	Stream pb.DownloadService_DownloadServer
	Task   *pb.DownloadTask
}

func NewServerSendStream(settings *config.Settings, manager *dbmanager.DBManager, task *pb.DownloadTask, stream pb.DownloadService_DownloadServer) *ServerSendStream {
	return &ServerSendStream{
		Core: chunkstream.Core{
			Settings: settings,
			Manager:  manager,
		},
		Stream: stream,
		Task:   task,
	}
}

func (s *ServerSendStream) SendStreamChunks() error {
	if len(s.Task.ChunkList) == 0 {
		return s.SendChunk(s.LoadEmptyChunk())
	}

	for _, idx := range s.Task.ChunkList {
		if err := s.SendChunk(s.LoadChunk(idx)); err != nil {
			return err
		}
	}
	return nil
}

func (s *ServerSendStream) SendChunk(chunk *pb.FileChunk) error {
	return s.Stream.Send(chunk)
}

func (s *ServerSendStream) LoadChunk(chunkIdx int32) *pb.FileChunk {
	byteSlice := chunkio.UploadChunk(s.Settings.CacheDirectory, s.Task.Meta.Sha256, chunkIdx)

	return &pb.FileChunk{
		Sha256:     s.Task.Meta.Sha256,
		ChunkIndex: chunkIdx,
		Data:       byteSlice,
	}
}

func (s *ServerSendStream) LoadEmptyChunk() *pb.FileChunk {
	return &pb.FileChunk{
		Sha256:     s.Task.Meta.Sha256,
		ChunkIndex: 0,
		Data:       []byte{},
	}
}
