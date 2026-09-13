package services

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"youtube_converter/common"
)

const CookieFile = "./cookies.txt"

func DownloadVideo(link, format, output, taskID string) error {
	if err := ensureCookies(); err != nil {
		return err
	}

	args := buildArgs(format, output, link)
	cmd := exec.Command("yt-dlp", args...)
	log.Printf("Executing command: %s", cmd.String())

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	go handleOutput(stdout, taskID)
	stderrOutput := make(chan string, 1)
	go collectErrors(stderr, stderrOutput)

	if err := cmd.Wait(); err != nil {
		details := <-stderrOutput
		if details != "" {
			return fmt.Errorf("yt-dlp failed: %s", details)
		}
		return fmt.Errorf("yt-dlp failed: %w", err)
	}
	<-stderrOutput
	return nil
}

func ensureCookies() error {
	if info, err := os.Stat(CookieFile); err == nil && info.Size() > 0 {
		return nil
	}

	encoded := strings.TrimSpace(os.Getenv("YOUTUBE_COOKIES_B64"))
	if encoded == "" {
		return nil
	}

	cookies, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("invalid YOUTUBE_COOKIES_B64: %w", err)
	}

	if err := os.WriteFile(CookieFile, cookies, 0600); err != nil {
		return fmt.Errorf("failed to write YouTube cookies: %w", err)
	}
	return nil
}

func buildArgs(format, output, link string) []string {
	args := []string{
		"-o", output,
		"--no-playlist",
		"--extractor-args", "youtube:player_client=android,web_safari",
		"--concurrent-fragments", "32",
		"--progress",
		"--newline",
	}

	if info, err := os.Stat(CookieFile); err == nil && info.Size() > 0 {
		args = append(args, "--cookies", CookieFile)
	}

	if format == "mp3" {
		args = append(args, "--extract-audio", "--audio-format", "mp3", "--audio-quality", "0")
	} else if format == "m4a" {
		args = append(args, "--format", "bestaudio[ext=m4a]/bestaudio")
	} else {
		args = append(args, "--format", fmt.Sprintf("bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"))
	}

	args = append(args, link)
	return args
}

func handleOutput(stdout io.ReadCloser, taskID string) {
	defer func(stdout io.ReadCloser) {
		err := stdout.Close()
		if err != nil {
			log.Printf("Error closing stdout: %v", err)
		}
	}(stdout)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "[download]") {
			progress := parseProgress(line)
			common.Broadcast <- common.ProgressMessage{TaskID: taskID, Percentage: progress}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from stdout: %v", err)
	}
}

func collectErrors(stderr io.ReadCloser, output chan<- string) {
	defer stderr.Close()
	scanner := bufio.NewScanner(stderr)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			log.Printf("[YT-DLP ERROR] %s", line)
			lines = append(lines, line)
		}
	}

	details := strings.Join(lines, " ")
	if len(details) > 900 {
		details = details[len(details)-900:]
	}
	output <- details
}

func parseProgress(line string) float64 {
	re := regexp.MustCompile(`(\d+\.\d+)%`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		progress, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			return progress
		}
		log.Printf("Error parsing progress: %v", err)
	}
	return 0.0
}
