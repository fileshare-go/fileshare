package send

import (
	"time"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/core"
	"github.com/chanmaoganda/fileshare/internal/core/chunkstream"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/chunkio"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/chanmaoganda/fileshare/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type ServerSendStream struct {
	chunkstream.Core
	Stream pb.DownloadService_DownloadServer
	Task   *pb.DownloadTask
}

func NewServerSendStream(task *pb.DownloadTask, stream pb.DownloadService_DownloadServer) chunkstream.StreamSendCore {
	return &ServerSendStream{
		Core:   chunkstream.Core{},
		Stream: stream,
		Task:   task,
	}
}

func (s *ServerSendStream) SendStreamChunks() error {
	if len(s.Task.ChunkList) == 0 || s.ValidateTask() {
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
	s.SetupAndRecordInfo(chunk)
	return s.Stream.Send(chunk)
}

func (s *ServerSendStream) CloseStream() error {
	service.Mgr().InsertRecord(s.MakeRecord())

	logrus.Debug("Closing server sending stream")
	return nil
}

func (s *ServerSendStream) ValidateTask() bool {
	if err := service.Mgr().SelectFileInfo(&s.FileInfo); err != nil {
		logrus.Error(err)
		return false
	}

	for _, chunkIdx := range s.Task.ChunkList {
		if chunkIdx >= 0 && chunkIdx < s.FileInfo.ChunkNumber {
			return false
		}
	}

	return true
}

func (s *ServerSendStream) MakeRecord() *model.Record {
	return &model.Record{
		Sha256:         s.FileInfo.Sha256,
		InteractAction: core.DownloadAction,
		ClientIp:       s.PeerAddress(),
		Os:             s.PeerOs(),
		Time:           time.Now(),
	}
}

func (s *ServerSendStream) PeerAddress() string {
	peer, ok := peer.FromContext(s.Stream.Context())
	if ok {
		return peer.Addr.String()
	}
	return "unknown"
}

func (s *ServerSendStream) PeerOs() string {
	md, ok := metadata.FromIncomingContext(s.Stream.Context())
	if !ok {
		return "unknown"
	}

	if osInfo, ok := md["os"]; ok && len(osInfo) != 0 {
		return osInfo[0]
	}
	return "unknown"
}

func (s *ServerSendStream) LoadChunk(chunkIdx int32) *pb.FileChunk {
	chunkData := chunkio.ReadChunk(config.Cfg().CacheDirectory, s.Task.Meta.Sha256, chunkIdx)

	return &pb.FileChunk{
		Sha256:     s.Task.Meta.Sha256,
		ChunkIndex: chunkIdx,
		Data:       chunkData,
	}
}

func (s *ServerSendStream) LoadEmptyChunk() *pb.FileChunk {
	return &pb.FileChunk{
		Sha256:     s.Task.Meta.Sha256,
		ChunkIndex: 0,
		Data:       []byte{},
	}
}
