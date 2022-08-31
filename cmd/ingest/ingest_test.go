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
