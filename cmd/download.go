package cmd

import (
	"github.com/mnsdojo/gotube/internal/downloader"
	"github.com/mnsdojo/gotube/internal/utils"
	"github.com/spf13/cobra"
)

var (
	format      string
	audioOnly   bool
	outputDir   string
	concurrency int
	quality     string
	verbose     bool
)

var downloadCmd = &cobra.Command{
	Use:   "download [URL]",
	Short: "Download a YouTube video or playlist",
	Long: `Download a single YouTube video or an entire playlist.
Supports format conversion and quality selection.`,
	Args: cobra.ExactArgs(1),
	RunE: runDownload,
}

func init() {
	downloadCmd.Flags().StringVarP(&format, "format", "f", "mp4", "Output format (mp4, mp3, webm, etc.)")
	downloadCmd.Flags().BoolVarP(&audioOnly, "audio", "a", false, "Download audio only")
	downloadCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	downloadCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 3, "Number of concurrent downloads for playlists")
	downloadCmd.Flags().StringVarP(&quality, "quality", "q", "highest", "Video quality (highest, 1080p, 720p, 480p, 360p, lowest)")
	downloadCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}

func runDownload(cmd *cobra.Command, args []string) error {
	url := args[0]
	config := downloader.Config{
		Format:      format,
		AudioOnly:   audioOnly,
		OutputDir:   outputDir,
		Concurrency: concurrency,
		Quality:     quality,
		Verbose:     verbose,
	}

	if utils.IsPlaylist(url) {
		return downloader.DownloadPlaylist(url, config)
	}
	return downloader.DownloadVideo(url, config)
}
