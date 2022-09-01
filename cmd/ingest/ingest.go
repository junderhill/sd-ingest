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
		includes := GetIncludes()

		filteredFilenames := FilterFiles(allFiles, includes)

		if viper.GetBool("verbose") {
			for _, i := range filteredFilenames {
				fmt.Printf("Found %s\n", i)
			}
		}

		files := make([]util.File, 0)
		for _, f := range filteredFilenames {
			files = append(files, *util.NewFile(f))
		}

		photos, videos := GroupFiles(files)

		wg := sync.WaitGroup{}
		wg.Add(2)

		totalBytesCopied := &util.SafeInt64Count{}
		totalFilesCopied := &util.SafeIntCount{}

		go CopyFiles(photos, viper.GetString("photos.destination"), &wg, totalFilesCopied, totalBytesCopied)
		go CopyFiles(videos, viper.GetString("video.destination"), &wg, totalFilesCopied, totalBytesCopied)
		wg.Wait()

		fmt.Println("Ingest Complete")
		fmt.Printf("Total Files: %d Total Size: %s", totalFilesCopied.Value(), util.ByteCountIEC(totalBytesCopied.Value()))
	},
}

func GetIncludes() []string {
	photoExts := viper.GetStringSlice("photos.formats")
	videoExts := viper.GetStringSlice("video.formats")

	output := make([]string, 0)

	for _, ext := range photoExts {
		output = append(output, strings.ToLower(ext))
	}
	for _, ext := range videoExts {
		output = append(output, strings.ToLower(ext))
	}

	return output
}

func CopyFiles(files map[string][]util.File, destinationPath string, s *sync.WaitGroup, fileCount *util.SafeIntCount, byteCount *util.SafeInt64Count) {
	for key, value := range files {

		dest := GetPath(destinationPath, key, destSuffix)

		if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
			if viper.GetBool("verbose") {
				fmt.Printf("Creating Directory: %s \n", dest)
			}
			err := os.Mkdir(dest, os.ModePerm)
			if err != nil && !os.IsExist(err) {
				fmt.Println(err)
			}
		}

		for _, file := range value {
			destinationFilename := fmt.Sprintf("%s/%s", dest, file.Filename)
			if viper.GetBool("verbose") {
				fmt.Printf("Copying %s to %s \n", file.Filename, destinationFilename)
			}
			_, err := copyFile(file.Path, destinationFilename)
			if err != nil {
				log.Println(err)
			} else {
				fileCount.Increment(1)
				byteCount.Increment(file.Size)
			}
		}
	}

	s.Done()
}

func copyFile(in, out string) (int64, error) {
	i, e := os.Open(in)
	if e != nil {
		return 0, e
	}
	defer i.Close()
	o, e := os.Create(out)
	if e != nil {
		return 0, e
	}
	defer o.Close()
	return o.ReadFrom(i)
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

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files
}
