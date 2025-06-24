package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetupLogger() {
	levelString := os.Getenv("LOG")

	level, err := logrus.ParseLevel(levelString)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
}
