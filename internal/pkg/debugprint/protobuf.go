package debugprint

import (
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

var Render = color.FgCyan.Render

func DebugUploadTask(task *pb.UploadTask) {
	logrus.Debugf("[Task] filename: %s, sha256: %s, file size %s",
		Render(task.Meta.Filename), Render(task.Meta.Sha256[:8]), Render(task.Meta.FileSize),
	)
	logrus.Debugf("[Task] chunk number: %s, chunk size: %s, needed chunks: %s",
		Render(task.ChunkNumber), Render(task.ChunkSize), Render(task.ChunkList),
	)
}

func DebugDownloadTask(task *pb.DownloadTask) {
	logrus.Debugf("[Task] filename: %s, sha256: %s, file size %s, chunk number: %s",
		Render(task.Meta.Filename), Render(task.Meta.Sha256[:8]), Render(task.Meta.FileSize),
		Render(task.ChunkNumber))
}

func DebugDownloadSummary(summary *pb.DownloadSummary) {
	logrus.Debugf("[Summary] filename: %s, sha256: %s, file size %s, chunk number: %s, chunk size: %s",
		Render(summary.Meta.Filename), Render(summary.Meta.Sha256[:8]), Render(summary.Meta.FileSize),
		Render(summary.ChunkNumber), Render(summary.ChunkSize))
}

func DebugChunk(chunk *pb.FileChunk) {
	logrus.Debugf("[Chunk] file sha256: %s, chunk index: %d, chunk size: %d",
		Render(chunk.Sha256[:8]), chunk.ChunkIndex, len(chunk.Data))
}

func DebugMeta(meta *pb.FileMeta) {
	logrus.Debugf("[Meta] File meta [filename: %s, sha256: %s]", meta.Filename, meta.Sha256[:8])
}
