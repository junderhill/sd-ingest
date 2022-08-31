package ingest

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/fs"
	"path/filepath"
	"sd-ingest/util"
	"strings"
)

var source string
var destSuffix string

func init() {
	CmdIngest.Flags().StringVar(&source, "source", "", "Set the source directory to copy files from")
	CmdIngest.MarkFlagRequired("source")
	CmdIngest.Flags().StringVar(&destSuffix, "dest-suffix", "", "(Optional) Suffix to add to destination directories.")

	viper.BindPFlag("source", CmdIngest.Flags().Lookup("source"))
	viper.BindPFlag("dest-suffix", CmdIngest.Flags().Lookup("dest-suffix"))
}

var CmdIngest = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest from SD card",
	Run: func(cmd *cobra.Command, args []string) {
		allFiles := GetFilesFromSource()

		includes := []string{ //todo: pull this from the yaml config
			"arw",
			"jpg",
			"dng",
			"mp4",
		}

		filteredFilenames := FilterFiles(allFiles, includes)

		if viper.GetBool("verbose") {
			fmt.Printf("Filtered files: %s\n", filteredFilenames)
		}

		files := make([]util.File, 0)
		for _, f := range filteredFilenames {
			files = append(files, *util.NewFile(f))
		}

		if viper.GetBool("verbose") {
			fmt.Printf("Converted files: %+v\n", files)
		}

	},
}

func FilterFiles(files []string, includes []string) []string {
	filtered := make([]string, 0)
	for _, v := range files {

		for _, ext := range includes {
			if strings.HasSuffix(strings.ToLower(v), ext) {
				filtered = append(filtered, v)
				break
			}
		}
	}

	return filtered
}

func GetFilesFromSource() []string {
	sourceDirectory := viper.GetString("source")
	if viper.GetBool("verbose") {
		fmt.Printf("Source Directory: %s\n", sourceDirectory)
	}

	var files []string

	filepath.WalkDir(sourceDirectory, func(path string, d fs.DirEntry, err error) error {

		if viper.GetBool("verbose") {
			fmt.Printf("Found %s\n", path)
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files
}
