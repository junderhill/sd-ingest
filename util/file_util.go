package util

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type File struct {
	Filename  string
	Path      string
	Type      string
	Timestamp time.Time
}

var videoSuffix = []string{
	"mp4",
	"mov",
	"avi",
}
var photoSuffix = []string{
	"jpg",
	"jpeg",
	"arw",
	"dng",
}

func NewFile(path string) *File {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Print(err)
	}
	lowerName := strings.ToLower(fileInfo.Name())
	var fileType string

	for _, ext := range photoSuffix {
		if strings.HasSuffix(lowerName, ext) {
			fileType = "photo"
			break
		}
	}
	for _, ext := range videoSuffix {
		if strings.HasSuffix(lowerName, ext) {
			fileType = "video"
			break
		}
	}

	return &File{
		Path:      path,
		Type:      fileType,
		Filename:  fileInfo.Name(),
		Timestamp: fileInfo.ModTime(),
	}
}
