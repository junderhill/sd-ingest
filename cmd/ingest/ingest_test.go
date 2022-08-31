package ingest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var includes []string = []string{
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
