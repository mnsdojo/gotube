package utils

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func IsPlaylist(url string) bool {
	return strings.Contains(url, "playlist?list=")
}

func SanitizeFilename(filename string) string {
	reg := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := reg.ReplaceAllString(filename, "")
	sanitized = strings.TrimSpace(sanitized)
	if len(sanitized) > 200 {
		sanitized = sanitized[:200]
	}
	return sanitized
}

func SelectFormat(formats youtube.FormatList, quality string, audioOnly bool) *youtube.Format {
	if audioOnly {
		return selectAudioFormat(formats)
	}
	return selectVideoFormat(formats, quality)
}

func selectAudioFormat(formats youtube.FormatList) *youtube.Format {
	audioFormats := formats.Type("audio")
	if len(audioFormats) == 0 {
		return nil
	}
	// Select the audio format with the highest bitrate
	best := &audioFormats[0]
	for i := 1; i < len(audioFormats); i++ {
		if audioFormats[i].Bitrate > best.Bitrate {
			best = &audioFormats[i]
		}
	}
	return best
}

func selectVideoFormat(formats youtube.FormatList, quality string) *youtube.Format {
	videoFormats := formats.Type("video").WithAudioChannels()
	if len(videoFormats) == 0 {
		return nil
	}

	switch quality {
	case "highest":
		return &videoFormats[0]
	case "lowest":
		return &videoFormats[len(videoFormats)-1]
	default:
		return findByQuality(videoFormats, quality)
	}
}

func findByQuality(formats youtube.FormatList, quality string) *youtube.Format {
	targetHeight := parseQuality(quality)
	if targetHeight == 0 {
		return &formats[0] // Return highest quality if parsing fails
	}

	var bestMatch *youtube.Format
	smallestDiff := int(^uint(0) >> 1) // Initialize with max int

	for i, format := range formats {
		height := parseQuality(format.Quality)
		diff := abs(height - targetHeight)
		if diff < smallestDiff {
			smallestDiff = diff
			bestMatch = &formats[i]
		}
	}

	return bestMatch
}

func parseQuality(quality string) int {
	// Remove "p" suffix if present
	quality = strings.TrimSuffix(quality, "p")

	// Try to parse as integer
	height, err := strconv.Atoi(quality)
	if err != nil {
		return 0
	}
	return height
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
