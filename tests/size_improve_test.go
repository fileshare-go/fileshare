package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"google.golang.org/protobuf/proto"
)

func TestLoaderSize(t *testing.T) {
	file, err := os.Open("fileshare")
	if err != nil {
		t.Error(err)
	}

	for size := 1; size < 8; size += 1 {
		chunkSize := 1024 * 512 * int64(size)
		data := make([]byte, chunkSize)
		file.Read(data)
		chunk := &pb.FileChunk{
			Sha256: "7b852e938bc09de10cd96eca3755258c7d25fb89dbdd76305717607e1835e2aa",
			ChunkIndex: int32(size),
			Data: data,
		}
		jsonData, _ := json.Marshal(chunk)
		protoData, _ := proto.Marshal(chunk)

		percentage := float64(len(jsonData) - len(protoData)) / float64(len(jsonData))
		fmt.Printf("json len %d, proto len %d, percentage %f\n", len(jsonData), len(protoData), percentage)
	}
}