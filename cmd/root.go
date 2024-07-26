package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gotube",
	Short: "GoTube - YouTube downloader CLI",
	Long: `GoTube is a powerful CLI tool to download YouTube videos and playlists.
It supports various formats, quality options, and concurrent downloads.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(infoCmd)
}
