package main

import (
	"errors"
	"github.com/pkgz/logg"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var out string
var link string
var debug bool

var rootCmd = &cobra.Command{
	Use:     "AOSDownloader",
	Short:   "Apple OpenSource download tool",
	Version: "0.0.1",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			logg.DebugMode()
		}

		if link == "" && len(args) > 0 {
			link = args[0]
		}
		if link == "" {
			return errors.New("please provide project url")
		}

		if out == "" && len(args) >= 2 {
			out = args[1]
		} else if out == "" && len(args) == 1 {
			out = strings.Replace(link, "https://opensource.apple.com/source/", "", 1)
			out = strings.Split(out, "/")[0]
		}
		if out == "" {
			return errors.New("please provide destination path")
		}

		if !strings.Contains(link, "https://opensource.apple.com/source") {
			return errors.New("project url must contain https://opensource.apple.com/source")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("[INFO] Fetching project %s...", link)

		links, n, err := parseURI(strings.TrimSuffix(link, "/"), "/")
		if err != nil {
			log.Printf("[ERROR] parse: %v", err)
			return err
		}

		log.Printf("[INFO] Detect %d files in project", n)

		if err := download(links, out); err != nil {
			log.Printf("[ERROR] download: %v", err)
			return err
		}

		log.Printf("[INFO] Project successfully downloaded to `%s`", out)
		return nil
	},
}

func init() {
	logg.NewGlobal(os.Stdout)
	logg.SetFlags(0)

	rootCmd.Flags().StringVarP(&link, "url", "u", "", "url to project which you want to download")
	rootCmd.Flags().StringVarP(&out, "out", "o", "", "destination path for project")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug mode")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
