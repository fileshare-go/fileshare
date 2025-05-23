package chunk

const (
	SMALL       = 1
	MEDIUM      = 2
	LARGE       = 3
	EXTRA_LARGE = 4
)

type ChunkSummary struct {
	Size   int64
	Number int32
}

func DealChunkSize(fileSize int64) ChunkSummary {
	chunkSize := 1024 * 1024
	chunkNumber := fileSize / int64(chunkSize)

	if fileSize%int64(chunkSize) != 0 {
		chunkNumber += 1
	}

	return ChunkSummary{
		Size:   int64(chunkSize),
		Number: int32(chunkNumber),
	}
}
