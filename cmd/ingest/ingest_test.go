package ingest

import (
	"testing"
	"time"

	"sd-ingest/util"

	"github.com/stretchr/testify/assert"
)

var includes = []string{
	"jpg",
	"dng",
}

func TestFilterFilesExcludesFilesNotInIncludeList(t *testing.T) {

	input := []string{
		"/User/Test/somefile.txt",
		"/User/Test/document.docx",
		"/User/Test/archive.tar.gz",
	}

	output := FilterFiles(input, includes)

	assert.Empty(t, output, "Expected non of the input values to be retained")
}

func TestFilterFilesIncludesExpectedFiles(t *testing.T) {
	input := []string{
		"/User/Test/somefile.txt",
		"/User/Test/DSC0001.DNG",
		"/User/Test/document.docx",
		"/User/Test/archive.tar.gz",
		"/User/Test/DSC0002.jpg",
	}

	output := FilterFiles(input, includes)

	assert.Len(t, output, 2, "Expected 2 files to be retained after filtering")
}

func TestFilterFilesRetainsAllIfNonToBeExcluded(t *testing.T) {
	input := []string{
		"/User/Test/DSC0001.DNG",
		"/User/Test/DSC0002.jpg",
		"/User/Test/DSC0004.DNG",
		"/User/Test/DSC0005.DNG",
		"/User/Test/DSC0031.dng",
		"/User/Test/DCS0034.jpg",
		"/User/Test/DSC0041.jpg",
		"/User/Test/DSC0035.JPG",
	}

	output := FilterFiles(input, includes)

	assert.ElementsMatch(t, input, output, "Expected all input elements to be included in filtered result.")
}

func TestGroupFilesGroupsPhotosAndVideosSeperately(t *testing.T) {
	input := []util.File{
		{
			Filename:  "DSC0001.jpg",
			Path:      "/temp/DSC0001.jpg",
			Type:      "photo",
			Timestamp: time.Date(2022, time.August, 30, 12, 30, 00, 00, time.UTC),
		},
		{
			Filename:  "DSC0002.mp4",
			Path:      "/temp/DSC0002.mp4",
			Type:      "video",
			Timestamp: time.Date(2022, time.August, 30, 12, 30, 00, 00, time.UTC),
		},
	}

	photos, videos := GroupFiles(input)

	assert.Len(t, photos, 1)
	assert.Len(t, videos, 1)

	for _, p := range photos {
		assert.Len(t, p, 1)
		assert.Equal(t, p[0].Type, "photo")
	}

	for _, v := range videos {
		assert.Len(t, v, 1)
		assert.Equal(t, v[0].Type, "video")
	}
}

func TestGroupFilesIncludesPhotosFromSameDateInSameMapKey(t *testing.T) {
	input := []util.File{
		{
			Filename:  "DSC0001.jpg",
			Path:      "/temp/DSC0001.jpg",
			Type:      "photo",
			Timestamp: time.Date(2022, time.August, 30, 12, 15, 00, 00, time.UTC),
		},
		{
			Filename:  "DSC0002.jpg",
			Path:      "/temp/DSC0002.jpg",
			Type:      "photo",
			Timestamp: time.Date(2022, time.August, 30, 12, 20, 00, 00, time.UTC),
		},
		{
			Filename:  "DSC0099.jpg",
			Path:      "/temp/DSC0099.jpg",
			Type:      "photo",
			Timestamp: time.Date(2022, time.August, 31, 12, 25, 00, 00, time.UTC),
		},
		{
			Filename:  "DSC0003.jpg",
			Path:      "/temp/DSC0003.jpg",
			Type:      "photo",
			Timestamp: time.Date(2022, time.August, 30, 12, 25, 00, 00, time.UTC),
		},
	}

	photos, _ := GroupFiles(input)

	assert.Len(t, photos["20220830"], 3)
}
