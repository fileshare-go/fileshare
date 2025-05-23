package chunk

const (
	SMALL  = 1024 * 1024 // 1MB
	MEDIUM = 2 * SMALL   // 2MB
	LARGE  = 4 * SMALL   // 4MB
)

type ChunkSummary struct {
	Size   int64
	Number int32
}

func DealChunkSize(fileSize int64) ChunkSummary {
	chunkSize := SMALL
	chunkNumber := fileSize / int64(chunkSize)

	if fileSize%int64(chunkSize) != 0 {
		chunkNumber += 1
	}

	return ChunkSummary{
		Size:   int64(chunkSize),
		Number: int32(chunkNumber),
	}
}
