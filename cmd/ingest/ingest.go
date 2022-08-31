package ingest

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sd-ingest/util"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			fmt.Printf("Grouping...")
		}

		photos, videos := GroupFiles(files)

		if viper.GetBool("verbose") {
			fmt.Println()
			fmt.Printf("Photos %v \n", photos)
			fmt.Printf("Videos %v \n", videos)
		}

		wg := sync.WaitGroup{}
		wg.Add(2)

		CopyFiles(photos, viper.GetString("photos.destination"), &wg)
		CopyFiles(videos, viper.GetString("videos.destination"), &wg)

		wg.Wait()

	},
}

func CopyFiles(files map[string][]util.File, destinationPath string, s *sync.WaitGroup) {

	for key, value := range files {

		dest := GetPath(destinationPath, key, destSuffix)

		if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(dest, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}

		for _, file := range value {
			
			//todo: copy file
		}
	}

	s.Done()
}

func GetPath(basePath string, date string, suffix string) string {
	//todo: look for a package to allow this work cross platform
	newDirectoryName := date

	if suffix != "" {
		newDirectoryName = fmt.Sprintf("%s_%s", date, suffix)
	}

	if strings.HasSuffix(basePath, "/") {
		return fmt.Sprintf("%s%s", basePath, newDirectoryName)
	} else {
		return fmt.Sprintf("%s/%s", basePath, newDirectoryName)
	}
}

func GroupFiles(files []util.File) (map[string][]util.File, map[string][]util.File) {

	//returns 2 maps of file slices
	// map key is the date in a string format '20220830'

	photos := make(map[string][]util.File)
	videos := make(map[string][]util.File)

	for _, v := range files {
		dateStr := v.Timestamp.Format("20060102")

		if v.Type == "video" {
			dateSlice, exists := videos[dateStr]
			if !exists {
				dateSlice = make([]util.File, 0)
			}

			videos[dateStr] = append(dateSlice, v)
		} else {
			dateSlice, exists := photos[dateStr]
			if !exists {
				dateSlice = make([]util.File, 0)
			}

			photos[dateStr] = append(dateSlice, v)
		}
	}

	return photos, videos
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
