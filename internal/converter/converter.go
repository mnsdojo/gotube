package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConvertFormat(inputFile, format string, audioOnly bool) error {
	outputFile := getOutputFilename(inputFile, format)

	args := []string{"-i", inputFile}

	if audioOnly {
		switch format {
		case "mp3":
			args = append(args, "-vn", "-acodec", "libmp3lame", "-b:a", "192k")
		case "m4a":
			args = append(args, "-vn", "-acodec", "aac", "-b:a", "192k")
		default:
			return fmt.Errorf("unsupported audio format: %s", format)
		}
	} else {
		switch format {
		case "mp4":
			args = append(args, "-c:v", "libx264", "-crf", "23", "-c:a", "aac", "-b:a", "128k")
		case "webm":
			args = append(args, "-c:v", "libvpx-vp9", "-crf", "30", "-b:v", "0", "-b:a", "128k", "-c:a", "libopus")
		default:
			return fmt.Errorf("unsupported video format: %s", format)
		}
	}

	args = append(args, outputFile)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Converting %s to %s...\n", inputFile, outputFile)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("conversion failed: %v", err)
	}

	fmt.Printf("Conversion complete. Output file: %s\n", outputFile)
	return nil
}

func getOutputFilename(inputFile, format string) string {
	dir := filepath.Dir(inputFile)
	filename := filepath.Base(inputFile)
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	return filepath.Join(dir, nameWithoutExt+"."+format)
}
