package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/kkdai/youtube/v2"
	"github.com/mnsdojo/gotube/internal/converter"
	"github.com/mnsdojo/gotube/internal/utils"
	"github.com/schollz/progressbar/v3"
)

type Config struct {
	Format      string
	AudioOnly   bool
	OutputDir   string
	Concurrency int
	Quality     string
	Verbose     bool
}

func DownloadVideo(url string, config Config) error {
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		return err
	}

	format := utils.SelectFormat(video.Formats, config.Quality, config.AudioOnly)
	if format == nil {
		return fmt.Errorf("no suitable format found")
	}

	stream, size, err := client.GetStream(video, format)
	if err != nil {
		return err
	}
	defer stream.Close()

	outputPath := filepath.Join(config.OutputDir, utils.SanitizeFilename(video.Title)+"."+config.Format)
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	bar := progressbar.DefaultBytes(
		size,
		"Downloading",
	)

	// Create a buffered reader and copy with progress update
	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			bar.Add(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	if config.Format != "mp4" || config.AudioOnly {
		return converter.ConvertFormat(outputPath, config.Format, config.AudioOnly)
	}
	return nil
}

func DownloadPlaylist(url string, config Config) error {
	client := youtube.Client{}
	playlist, err := client.GetPlaylist(url)
	if err != nil {
		return err
	}

	sem := make(chan bool, config.Concurrency)
	var wg sync.WaitGroup

	for _, video := range playlist.Videos {
		wg.Add(1)
		sem <- true
		go func(video *youtube.PlaylistEntry) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := DownloadVideo(video.ID, config); err != nil && config.Verbose {
				fmt.Printf("Error downloading video %s: %v\n", video.ID, err)
			}
		}(video)
	}

	wg.Wait()
	return nil
}


