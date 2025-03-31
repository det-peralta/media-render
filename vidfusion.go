package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const version = "VidFusion 1.0.0"

func detectHardware() string {
	cmd := exec.Command("ffmpeg", "-hide_banner", "-hwaccels")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error detecting hardware acceleration:", err)
		return "cpu"
	}

	if strings.Contains(string(output), "cuda") {
		return "nvidia"
	}
	if strings.Contains(string(output), "vulkan") {
		return "amd"
	}
	if strings.Contains(string(output), "qsv") {
		return "intel"
	}

	return "cpu"
}

func convertVideos(videoFiles []string, hardware string) ([]string, error) {
	convertedFiles := []string{}

	for _, file := range videoFiles {
		outputFile := strings.TrimSuffix(file, ".mp4") + "_converted.mp4"
		cmdArgs := []string{"-hide_banner", "-loglevel", "error", "-i", file}

		switch hardware {
		case "nvidia":
			cmdArgs = append(cmdArgs, "-c:v", "h264_nvenc")
		case "amd":
			cmdArgs = append(cmdArgs, "-c:v", "h264_amf")
		case "intel":
			cmdArgs = append(cmdArgs, "-c:v", "h264_qsv")
		default:
			cmdArgs = append(cmdArgs, "-c:v", "libx264")
		}

		cmdArgs = append(cmdArgs, outputFile)
		cmd := exec.Command("ffmpeg", cmdArgs...)
		cmd.Stderr = os.Stderr

		fmt.Printf("Converting %s using %s...\n", file, hardware)
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("error converting %s: %w", file, err)
		}

		convertedFiles = append(convertedFiles, outputFile)
	}

	return convertedFiles, nil
}

func main() {
	fmt.Printf("VidFusion - Version %s\n", version)

	if len(os.Args) < 2 {
		fmt.Println("Please drop video files onto this executable.")
		return
	}

	videoFiles := os.Args[1:]
	if len(videoFiles) == 1 {
		fmt.Println("Only one video file provided. Skipping concatenation.")
		fmt.Printf("Output file: %s\n", videoFiles[0])
		fmt.Println("Press Enter to exit...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}

	fmt.Println("Videos to concatenate:")
	for _, file := range videoFiles {
		fmt.Println("-", file)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to remove audio from the output? (yes/no): ")
	response, _ := reader.ReadString('\n')
	removeAudio := strings.TrimSpace(strings.ToLower(response)) == "yes"

	fmt.Print("Do you want to convert videos before concatenation? (yes/no): ")
	convertResponse, _ := reader.ReadString('\n')
	convertVideosFlag := strings.TrimSpace(strings.ToLower(convertResponse)) == "yes"

	if convertVideosFlag {
		hardware := detectHardware()
		fmt.Printf("Detected hardware: %s\n", hardware)
		converted, err := convertVideos(videoFiles, hardware)
		if err != nil {
			fmt.Println("Error during conversion:", err)
			return
		}
		videoFiles = converted
	}

	outputFile := "output.mp4"
	fmt.Printf("Output file will be: %s\n", outputFile)

	tempFile, err := os.Create("file_list.txt")
	if err != nil {
		fmt.Println("Error creating temporary file list:", err)
		return
	}
	defer os.Remove(tempFile.Name())

	for _, file := range videoFiles {
		if _, err := tempFile.WriteString(fmt.Sprintf("file '%s'\n", file)); err != nil {
			fmt.Println("Error writing to temporary file list:", err)
			return
		}
	}
	tempFile.Close()

	cmdArgs := []string{"-hide_banner", "-loglevel", "error", "-f", "concat", "-safe", "0", "-i", tempFile.Name()}
	if removeAudio {
		cmdArgs = append(cmdArgs, "-an")
	}
	cmdArgs = append(cmdArgs, "-c", "copy", outputFile)

	cmd := exec.Command("ffmpeg", cmdArgs...)
	cmd.Stderr = os.Stderr

	fmt.Println("Running FFmpeg for concatenation without re-encoding...")
	start := time.Now()
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running FFmpeg:", err)
		return
	}
	duration := time.Since(start)

	fmt.Printf("Videos concatenated successfully into %s\n", outputFile)
	fmt.Printf("Time taken: %s\n", duration)

	fmt.Println("Press Enter to exit...")
	reader.ReadString('\n')
}
