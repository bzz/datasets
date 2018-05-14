package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/src-d/datasets/PublicGitArchive/pga/pga"
	pb "gopkg.in/cheggaaa/pb.v1"
)

// getCmd represents the get command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "checks status of all repositories in the index",
	Long: `Checks the status of all repositories using md5
	Does not dowload anything`,
	RunE: func(cmd *cobra.Command, args []string) error {
		source := urlFS("http://pga.sourced.tech/")

		dest, err := FileSystemFromFlags(cmd.Flags())
		if err != nil {
			return err
		}

		stdin, err := cmd.Flags().GetBool("stdin")
		if err != nil {
			return err
		}

		var filenames []string
		if stdin {
			fmt.Fprintln(os.Stderr, "getting siva files by name from stdin")
			fmt.Fprintln(os.Stderr, "filter flags will be ignored")
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("could not read from standard input: %v", err)
			}
			for _, filename := range strings.Split(string(b), "\n") {
				filename = strings.TrimSpace(filename)
				if filename == "" {
					continue
				}
				filenames = append(filenames, filename)
			}
		} else {
			f, err := getIndex()
			if err != nil {
				return fmt.Errorf("could not open index file: %v", err)
			}
			defer f.Close()

			index, err := pga.IndexFromCSV(f)
			if err != nil {
				return err
			}

			filter, err := filterFromFlags(cmd.Flags())
			if err != nil {
				return err
			}

			index = pga.WithFilter(index, filter)

			for {
				r, err := index.Next()
				if err == io.EOF {
					break
				} else if err != nil {
					return err
				}
				filenames = append(filenames, r.Filenames...)
			}
		}

		logrus.Infof("Checking status for %d\n", len(filenames))
		maxDownloads, err := cmd.Flags().GetInt("jobs")
		if err != nil {
			return err
		}

		return checkFilenames(dest, source, filenames, maxDownloads)
	},
}

func checkFilenames(dest, source FileSystem, filenames []string, maxDownloads int) error {
	tokens := make(chan bool, maxDownloads)
	for i := 0; i < maxDownloads; i++ {
		tokens <- true
	}

	done := make(chan int)
	for _, filename := range filenames {
		filename := filepath.Join("siva", "latest", filename[:2], filename)
		go func() {
			<-tokens
			defer func() { tokens <- true }()

			if !upToDate(dest, source, filename) {
				logrus.Infof("File is corrupted or not up to date %s\n", filename)
				done <- 1
			} else {
				done <- 0
			}
		}()
	}

	corrupted := 0
	bar := pb.StartNew(len(filenames))
	for i := 1; ; i++ {
		corrupted += <-done
		bar.Set(i)
		bar.Update()
		if i == len(filenames) {
			logrus.Infof("%d of %d corrupted in '%s'", corrupted, len(filenames), dest.Abs(""))
			return nil
		}
	}
}

func init() {
	RootCmd.AddCommand(statusCmd)
	flags := statusCmd.Flags()
	addFilterFlags(flags)
	flags.StringP("output", "o", ".", "path where the siva files are stored")
	flags.BoolP("stdin", "i", false, "take list of siva files from standard input")
	flags.IntP("jobs", "j", 10, "number of concurrent gets allowed")
}
