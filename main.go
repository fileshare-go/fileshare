package main

import (
	"github.com/chanmaoganda/fileshare/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
