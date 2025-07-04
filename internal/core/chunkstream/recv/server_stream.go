package recv

import (
	"io"
	"time"

	"github.com/chanmaoganda/fileshare/internal/core"
	"github.com/chanmaoganda/fileshare/internal/core/chunkstream"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/chanmaoganda/fileshare/internal/service"
)

type ServerRecvStream struct {
	chunkstream.Core
	Stream pb.UploadService_UploadServer
}

func NewServerRecvStream(stream pb.UploadService_UploadServer) chunkstream.StreamRecvCore {
	return &ServerRecvStream{
		Core:   chunkstream.Core{},
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

	// merge missing chunks that has been uploaded
	// update current chunks whether err is nil or not
	s.FileInfo.UpdateChunks(s.ChunkList)

	// if err is not EOF, then this err should be handled
	if err != io.EOF {
		return err
	}
	return nil
}

func (s *ServerRecvStream) RecvChunk() (*pb.FileChunk, error) {
	return s.Stream.Recv()
}

func (s *ServerRecvStream) ValidateRecvChunks() bool {
	return s.Validate()
}

func (s *ServerRecvStream) CloseStream(validate bool) error {
	var err error
	if service.Orm().Save(&s.FileInfo).Error != nil {
		return err
	}
	ctx := s.Stream.Context()
	record := makeRecord(s.FileInfo.Sha256, util.PeerAddress(ctx), util.PeerOs(ctx))

	if service.Orm().Create(record).Error != nil {
		return err
	}

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

func makeRecord(sha256, peerAddress, peerOs string) *model.Record {
	return &model.Record{
		Sha256:         sha256,
		InteractAction: core.DownloadAction,
		ClientIp:       peerAddress,
		Os:             peerOs,
		Time:           time.Now(),
	}
}
