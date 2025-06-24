package upload

const (
	SMALL  = 1024 * 1024 // 1MB
	MEDIUM = 2 * SMALL   // 2MB
	LARGE  = 4 * SMALL   // 4MB
)

type chunkSummary struct {
	Size   int64
	Number int32
}

func dealChunkSize(fileSize int64) chunkSummary {
	var chunkSize int
	if fileSize < 64*SMALL {
		chunkSize = SMALL
	} else if fileSize < 1024*SMALL {
		chunkSize = MEDIUM
	} else {
		chunkSize = LARGE
	}

	chunkNumber := fileSize / int64(chunkSize)

	if fileSize%int64(chunkSize) != 0 {
		chunkNumber += 1
	}

	return chunkSummary{
		Size:   int64(chunkSize),
		Number: int32(chunkNumber),
	}
}
