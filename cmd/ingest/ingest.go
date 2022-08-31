package ingest

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var source string

func init() {
	CmdIngest.Flags().StringVar(&source, "source", "", "Set the source directory to copy files from")
	CmdIngest.MarkFlagRequired("source")
	viper.BindPFlag("source", CmdIngest.Flags().Lookup("source"))
}

var CmdIngest = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest from SD card",
	Run: func(cmd *cobra.Command, args []string) {
		files := GetFilesFromSource()
		fmt.Printf("Found files: %s\n", files)
	},
}

func GetFilesFromSource() []string {

	sourceDirectory := viper.GetString("source")
	if viper.GetBool("verbose") {
		fmt.Printf("Source Directory: %s\n", sourceDirectory)
	}

	files := []string{}

	filepath.WalkDir(sourceDirectory, func(path string, d fs.DirEntry, err error) error {

		if viper.GetBool("verbose") {
			fmt.Printf("Found %s\n", path)
		}

		files = append(files, path)
		return nil
	})

	return files
}
