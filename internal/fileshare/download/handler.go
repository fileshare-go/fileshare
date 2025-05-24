package download

import (
	"sync"

	pb "github.com/chanmaoganda/fileshare/proto/gen"
)

type Handler struct {
	stream           pb.DownloadService_DownloadClient
	once             sync.Once
	meta             *pb.FileMeta
	chunkList        []int32
	totalChunkNumber int32
	chunkSize        int64
}

func NewHandler(stream pb.DownloadService_DownloadClient) *Handler {
	return &Handler{
		stream: stream,
		meta:   &pb.FileMeta{},
	}
}

// func (h *Handler) recordInformation(chunk *pb.FileChunk) {
// 	h.meta.Filename = chunk.FileMeta.Filename
// 	h.meta.Sha256 = chunk.FileMeta.Sha256

// 	h.totalChunkNumber = chunk.ChunkMeta.Total
// }

// func (h *Handler) RecvAndSaveLock() error {
// 	for {
// 		chunk, err := h.stream.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}

// 		logrus.Debugf("filename: %s, total chunk: %d, chunk index: %d, chunk size: %d", chunk.FileMeta.Filename, chunk.ChunkMeta.Total, chunk.ChunkMeta.Index, len(chunk.GetData()))

// 		// create folder, record total chunk number and meta info
// 		h.once.Do(func() {
// 			h.recordInformation(chunk)
// 		})

// 		h.chunkList = append(h.chunkList, chunk.ChunkMeta.Index)

// 		if err := chunker.SaveChunk(chunk); err != nil {
// 			return err
// 		}
// 	}
// 	// return h.SaveLockFile()
// 	return nil
// }

// func (h *Handler) ValidateAndClose() {
// 	status := pb.Status_OK
// 	if chunker.ValidateChunks(h.meta.Filename, h.meta.Sha256) {
// 		logrus.Debugf("[validate] %s validated! sha256 is %s", h.meta.Filename, h.meta.Sha256)
// 	} else {
// 		status = pb.Status_ERROR
// 		logrus.Warnf("[validate] %s not validated!", h.meta.Filename)
// 	}

// 	uploadStatus := pb.UploadStatus{
// 		Meta:      h.meta,
// 		Status:    status,
// 		ChunkList: h.chunkList,
// 	}
// 	h.stream.CloseSend()
// 	// if err := h.stream.SendAndClose(&uploadStatus); err != nil {
// 	// 	logrus.Error("[validate] err: ", err)
// 	// }
// }
