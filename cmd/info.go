package cmd

import (
	"fmt"

	"github.com/kkdai/youtube/v2"
	"github.com/mnsdojo/gotube/internal/utils"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info [url]",
	Short: "info about yt video/playlist",
	Args:  cobra.ExactArgs(1),
	RunE:  infoUrl,
}

var isInfo bool

func init() {
	infoCmd.Flags().BoolVarP(&isInfo, "info", "i", false, "Show Information about url")

}
func runInfo(url string) error {
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return fmt.Errorf("error fetching video details: %w", err)
	}
	safeTitle := utils.SanitizeFilename(video.Title)
	fmt.Println(" Title:", safeTitle)
	fmt.Println("Author:", video.Author)
	fmt.Println("Duration:", video.Duration)
	for _, format := range video.Formats {
		fmt.Printf("Quality: %s, Type: %s, Resolution: %s\n", format.Quality, format.MimeType, format.QualityLabel)
	}

	return nil
}

func infoUrl(cmd *cobra.Command, args []string) error {
	url := args[0]
	return runInfo(url)
}
