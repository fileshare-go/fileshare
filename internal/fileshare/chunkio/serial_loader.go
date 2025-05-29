package chunkio

import (
	"os"

	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
)

type SerialChunkLoader struct {
	*os.File
	Sha256    string
	ChunkSize int64
	LastIdx   int32
	SeekAt    int64
	Data      []byte
}

func NewSerialChunkLoader(file *os.File, sha256 string, chunkSize int64) *SerialChunkLoader {
	return &SerialChunkLoader{
		File:      file,
		Sha256:    sha256,
		ChunkSize: chunkSize,
		LastIdx:   -1,
		SeekAt:    0,
		Data:      make([]byte, chunkSize),
	}
}

func (l *SerialChunkLoader) LoadChunk(index int32) *pb.FileChunk {
	if index-l.LastIdx == 1 {
		return l.loadNonSeekChunk()
	}
	return l.loadSeekChunk(index - l.LastIdx)
}

func (l *SerialChunkLoader) loadNonSeekChunk() *pb.FileChunk {
	n, err := l.Read(l.Data)
	if err != nil {
		logrus.Error(err)
		l.LastIdx += 1
		l.SeekAt += l.ChunkSize
		return l.EmptyChunk(l.LastIdx)
	}

	l.LastIdx += 1
	return &pb.FileChunk{
		Sha256:     l.Sha256,
		ChunkIndex: l.LastIdx,
		Data:       l.Data[:n],
	}
}

func (l *SerialChunkLoader) loadSeekChunk(gap int32) *pb.FileChunk {
	_, err := l.Seek(int64(gap)*l.ChunkSize, int(l.SeekAt))
	if err != nil {
		l.LastIdx += gap
		return l.EmptyChunk(l.LastIdx)
	}

	n, err := l.Read(l.Data)
	if err != nil {
		l.LastIdx += gap
		return l.EmptyChunk(l.LastIdx)
	}

	l.LastIdx += gap
	l.SeekAt += int64(gap) * l.ChunkSize
	return &pb.FileChunk{
		Sha256:     l.Sha256,
		ChunkIndex: l.LastIdx,
		Data:       l.Data[:n],
	}
}

func (l *SerialChunkLoader) EmptyChunk(index int32) *pb.FileChunk {
	return &pb.FileChunk{
		Sha256:     l.Sha256,
		ChunkIndex: index,
		Data:       []byte{},
	}
}
