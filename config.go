package main

import (
	"fmt"
	"os"
	"github.com/vaughan0/go-ini"
)

var config *ini.File

func init() {
	configPath := fmt.Sprintf("%s%capp.ini", runPath(), os.PathSeparator)
	tmpFile, err := ini.LoadFile(configPath)
	if err != nil {
		panic(err)
	}
	config = &tmpFile
}
