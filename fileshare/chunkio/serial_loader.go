package chunkio

import (
	"os"

	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
)

// serial chunk loader is designed for sending a sorted chunklist, which may be not continuous
//
// if chunklist is continuous, this loader performs like normal Read
// if not continuous, this loader will first seek to next start point and read
type SerialChunkLoader struct {
	*os.File
	sha256    string
	chunkSize int64
	lastIdx   int32
	seekAt    int64
}

func NewSerialChunkLoader(file *os.File, sha256 string, chunkSize int64) *SerialChunkLoader {
	return &SerialChunkLoader{
		File:      file,
		sha256:    sha256,
		chunkSize: chunkSize,
		lastIdx:   -1,
		seekAt:    0,
	}
}

// load a chunk from disk
func (l *SerialChunkLoader) LoadChunk(index int32) *pb.FileChunk {
	if index-l.lastIdx == 1 {
		return l.loadNonSeekingChunk()
	}
	return l.loadSeekingChunk(index - l.lastIdx)
}

// load a non seeking chunk will return next chunk according to chunksize
//
// mention that SeekAt will be updated
func (l *SerialChunkLoader) loadNonSeekingChunk() *pb.FileChunk {
	buffer := make([]byte, l.chunkSize)
	n, err := l.Read(buffer)
	if err != nil {
		logrus.Error(err)
		l.lastIdx += 1
		l.seekAt += l.chunkSize
		return l.EmptyChunk(l.lastIdx)
	}

	l.lastIdx += 1
	return &pb.FileChunk{
		Sha256:     l.sha256,
		ChunkIndex: l.lastIdx,
		Data:       buffer[:n],
	}
}

// load a seeking chunk will first performs seek on reader, then read the next chunk according to chunksize
//
// mention that SeekAt will be updated according to gap passed in
func (l *SerialChunkLoader) loadSeekingChunk(gap int32) *pb.FileChunk {
	buffer := make([]byte, l.chunkSize)
	_, err := l.Seek(int64(gap)*l.chunkSize, int(l.seekAt))
	if err != nil {
		l.lastIdx += gap
		return l.EmptyChunk(l.lastIdx)
	}

	n, err := l.Read(buffer)
	if err != nil {
		l.lastIdx += gap
		return l.EmptyChunk(l.lastIdx)
	}

	l.lastIdx += gap
	l.seekAt += int64(gap) * l.chunkSize
	return &pb.FileChunk{
		Sha256:     l.sha256,
		ChunkIndex: l.lastIdx,
		Data:       buffer[:n],
	}
}

// return a empty chunk to mark meta data for a file
func (l *SerialChunkLoader) EmptyChunk(index int32) *pb.FileChunk {
	return &pb.FileChunk{
		Sha256:     l.sha256,
		ChunkIndex: index,
		Data:       []byte{},
	}
}
