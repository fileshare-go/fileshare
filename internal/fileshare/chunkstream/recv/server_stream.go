package recv

import (
	"io"
	"time"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkstream"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"google.golang.org/grpc/peer"
)

type ServerRecvStream struct {
	chunkstream.Core
	Stream pb.UploadService_UploadServer
}

func NewServerRecvStream(settings *config.Settings, manager *dbmanager.DBManager, stream pb.UploadService_UploadServer) *ServerRecvStream {
	return &ServerRecvStream{
		Core: chunkstream.Core{
			Settings: settings,
			Manager:  manager,
		},
		Stream: stream,
	}
}

func (s *ServerRecvStream) RecvStreamChunks() error {
	var chunk *pb.FileChunk
	var err error

	for chunk, err = s.RecvChunk(); err == nil; chunk, err = s.RecvChunk() {
		if saveStatus := s.SaveChunkToDisk(chunk); !saveStatus {
			break
		}
	}

	// if err is not EOF, then this err should be handled
	if err != io.EOF {
		return err
	}

	// merge missing chunks that has been uploaded
	s.FileInfo.UpdateChunks(s.ChunkList)
	return nil
}

func (s *ServerRecvStream) RecvChunk() (*pb.FileChunk, error) {
	return s.Stream.Recv()
}

func (s *ServerRecvStream) PeerAddress() string {
	peer, ok := peer.FromContext(s.Stream.Context())
	if ok {
		return peer.Addr.String()
	}
	return "unknown"
}

func (s *ServerRecvStream) MakeRecord() *model.Record {
	return &model.Record{
		Sha256:         s.FileInfo.Sha256,
		InteractAction: "upload",
		ClientIp:       s.PeerAddress(),
		Time:           time.Now(),
	}
}

func (s *ServerRecvStream) CloseStream(validate bool) error {
	s.Manager.UpdateFileInfo(&s.FileInfo)
	s.Manager.CreateRecord(s.MakeRecord())

	status := s.genUploadStatus(validate)

	return s.Stream.SendAndClose(status)
}

func (s *ServerRecvStream) genUploadStatus(validate bool) *pb.UploadStatus {
	var statusCode pb.Status
	if validate {
		statusCode = pb.Status_OK
	} else {
		statusCode = pb.Status_ERROR
	}

	return &pb.UploadStatus{
		Meta: &pb.FileMeta{
			Filename: s.FileInfo.Filename,
			Sha256:   s.FileInfo.Sha256,
			FileSize: s.FileInfo.FileSize,
		},
		Status:    statusCode,
		ChunkList: s.ChunkList,
	}
}
