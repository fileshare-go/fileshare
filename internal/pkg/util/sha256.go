package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func CalculateFileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("Error opening file:", err)
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Warn(err)
		}
	}()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		logrus.Error("Error hashing file:", err)
		return "", err
	}

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}
