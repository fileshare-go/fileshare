package testing

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"google.golang.org/protobuf/proto"
)

func BenchmarkTransfer(b *testing.B) {
	file, err := os.Open("../.download/kafka_2.13-4.0.0.tgz")
	if err != nil {
		b.Error(err)
	}

	// Create a CSV file to record results
	csvFile, err := os.Create("benchmark_result.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"ChunkSizeKB", "ProtoLen", "JsonLen", "ProtoOverChunkRatio", "JsonOverProtoRatio"})


	for size := 1; size < 16; size += 1 {
		chunkSize := 1024 * 512 * int64(size)
		data := make([]byte, chunkSize)

		_, _ = file.Seek(0, 0)
		_, _ = file.Read(data)
		chunk := &pb.FileChunk{
			Sha256:     "7b852e938bc09de10cd96eca3755258c7d25fb89dbdd76305717607e1835e2aa",
			ChunkIndex: int32(size),
			Data:       data,
		}
		jsonData, _ := json.Marshal(chunk)
		protoData, _ := proto.Marshal(chunk)

		protoOverChunk := float64(len(protoData) - int(chunkSize)) / float64(chunkSize)
		fmt.Printf("chunk size %d, proto len %d, proto compared to chunk size overwhelming percentage %f\n", len(jsonData), len(protoData), protoOverChunk)

		jsonOverProto := float64(len(jsonData)-len(protoData)) / float64(len(protoData))
		fmt.Printf("json len %d, proto len %d, json compared to proto overwhelming percentage %f\n", len(jsonData), len(protoData), jsonOverProto)

		// Write to CSV
		writer.Write([]string{
			fmt.Sprintf("%d", chunkSize/1024),
			fmt.Sprintf("%d", len(protoData)),
			fmt.Sprintf("%d", len(jsonData)),
			fmt.Sprintf("%.6f", protoOverChunk),
			fmt.Sprintf("%.6f", jsonOverProto),
		})
	}
}
