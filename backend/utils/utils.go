package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wader/goutubedl"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"codeberg.org/pluja/whishper/models"
)

func SanitizeFilename(filename string) string {
	// First remove trailing spaces
	filename = strings.TrimSpace(filename)
	// Then remove quotes and dots
	filename = strings.Trim(filename, `"'.`)
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	filename = reg.ReplaceAllString(filename, "_")
	return filename
}

func DownloadMedia(t *models.Transcription) (string, error) {
	if t.SourceUrl == "" {
		log.Debug().Msg("Source URL is empty")
		return "", fmt.Errorf("source URL is empty")
	}

	if t.ID == primitive.NilObjectID {
		log.Debug().Msg("Transcription ID is empty")
		return "", fmt.Errorf("transcription ID is empty")
	}

	goutubedl.Path = "yt-dlp"
	result, err := goutubedl.New(context.Background(), t.SourceUrl, goutubedl.Options{})
	if err != nil {
		log.Debug().Err(err).Msg("Error creating goutubedl")
		return "", err
	}

	baseName := SanitizeFilename(fmt.Sprintf("%v%v%v", t.ID.Hex(), models.FileNameSeparator, result.Info.Title))

	// Download to a file so yt-dlp can mux the video and audio streams with
	// ffmpeg. That mux is a stream copy, not a re-encode. Streaming to stdout
	// cannot combine the two tracks.
	tmpDir, err := os.MkdirTemp("", "whishper-media-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	// YouTube randomly rejects some streams with HTTP 403 (missing PO token),
	// and which ones varies per video and per request. Try progressively
	// simpler requests until one succeeds; the last ones are audio-only.
	output := filepath.Join(tmpDir, "media.%(ext)s")
	attempts := []struct {
		name string
		args []string
	}{
		{"merged video", []string{
			"-f", "bv*[vcodec^=avc1]+ba[ext=m4a]/bv*[ext=mp4]+ba/bv*+ba/b",
			"--merge-output-format", "mp4",
		}},
		{"combined video (tv_simply client)", []string{
			"--extractor-args", "youtube:player_client=tv_simply",
			"-f", "b[ext=mp4]/b",
		}},
		{"audio only", []string{
			"-f", "ba[protocol=https]/ba/b",
			"-x", "--audio-format", "best",
		}},
		{"audio only (web_embedded client)", []string{
			"--extractor-args", "youtube:player_client=web_embedded",
			"-f", "ba/b",
			"-x", "--audio-format", "best",
		}},
	}

	var downloadErr error
	for _, attempt := range attempts {
		if err := resetDir(tmpDir); err != nil {
			return "", err
		}
		args := append([]string{"--no-playlist", "-o", output}, attempt.args...)
		args = append(args, t.SourceUrl)
		downloadErr = runYtDlp(args)
		if downloadErr == nil {
			log.Info().Msgf("Downloaded %v as %v", t.SourceUrl, attempt.name)
			break
		}
		log.Warn().Err(downloadErr).Msgf("Download attempt %q failed", attempt.name)
	}
	if downloadErr != nil {
		log.Debug().Err(downloadErr).Msg("Error downloading media")
		return "", downloadErr
	}

	src, err := downloadedMediaPath(tmpDir)
	if err != nil {
		return "", err
	}
	filename := baseName + filepath.Ext(src)
	dest := filepath.Join(os.Getenv("UPLOAD_DIR"), filename)
	if err := copyFile(src, dest); err != nil {
		log.Debug().Err(err).Msg("Error saving downloaded media")
		return "", err
	}

	return filename, nil
}

func resetDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.MkdirAll(dir, 0o700)
}

func runYtDlp(args []string) error {
	cmd := exec.CommandContext(context.Background(), "yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	return ytDlpError(output, err)
}

func ytDlpError(output []byte, err error) error {
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "ERROR:") {
			return errors.New(strings.TrimSpace(strings.TrimPrefix(line, "ERROR:")))
		}
	}
	if len(lines) > 0 && lines[len(lines)-1] != "" {
		return errors.New(lines[len(lines)-1])
	}
	return err
}

func downloadedMediaPath(dir string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "media.*"))
	if err != nil {
		return "", err
	}
	var candidates []string
	for _, match := range matches {
		if strings.HasSuffix(match, ".part") || strings.HasSuffix(match, ".ytdl") {
			continue
		}
		candidates = append(candidates, match)
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("downloaded file not found")
	}
	for _, candidate := range candidates {
		if strings.HasSuffix(candidate, ".mp4") {
			return candidate, nil
		}
	}
	return candidates[0], nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func SendTranscriptionRequest(t *models.Transcription, body *bytes.Buffer, writer *multipart.Writer) (*models.WhisperResult, error) {
	// 1. PREPARE URL PARAMETERS (Query String)
	// We use Query Params because FastAPI expects them there by default.
	baseUrl := fmt.Sprintf("http://%v/transcribe/", os.Getenv("ASR_ENDPOINT"))

	params := url.Values{}
	params.Add("model_size", t.ModelSize)
	params.Add("task", t.Task)
	params.Add("language", t.Language)
	params.Add("device", t.Device) // This ensures 'cuda' is passed cleanly

	if t.BeamSize > 0 {
		params.Add("beam_size", fmt.Sprintf("%d", t.BeamSize))
	}
	if t.InitialPrompt != "" {
		params.Add("initial_prompt", t.InitialPrompt)
	}
	if len(t.Hotwords) > 0 {
		params.Add("hotwords", strings.Join(t.Hotwords, ","))
	}
	// VAD parameters
	if t.VadFilter {
		params.Add("vad_filter", "true")
		if t.VadThreshold != nil {
			params.Add("vad_threshold", fmt.Sprintf("%f", *t.VadThreshold))
		}
		if t.VadMinSpeechDurationMS != nil {
			params.Add("vad_min_speech_duration_ms", fmt.Sprintf("%d", *t.VadMinSpeechDurationMS))
		}
		if t.VadMinSilenceDurationMS != nil {
			params.Add("vad_min_silence_duration_ms", fmt.Sprintf("%d", *t.VadMinSilenceDurationMS))
		}
	}

	// Attach parameters to URL
	fullUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	// 2. CRITICAL: CLOSE THE WRITER
	// This writes the final boundary to the 'body' buffer.
	// If you don't do this, the server gets a corrupted body and returns EOF.
	err := writer.Close()
	if err != nil {
		log.Debug().Err(err).Msg("Error closing multipart writer")
		return nil, err
	}

	// 3. SEND REQUEST
	req, err := http.NewRequest("POST", fullUrl, body)
	if err != nil {
		log.Debug().Err(err).Msg("Error creating request to transcription service")
		return nil, err
	}

	// Set the Content-Type with the boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error sending request")
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading response body")
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("Response from %v: %v", fullUrl, string(b))
		log.Error().Err(err).Msgf("Invalid response status %v:", resp.StatusCode)
		return nil, errors.New("invalid status")
	}

	var asrResponse *models.WhisperResult
	if err := json.Unmarshal(b, &asrResponse); err != nil {
		log.Error().Err(err).Msg("Error decoding response")
		log.Error().Msgf("ASR Response: %+v\n", b)
		return nil, err
	}

	return asrResponse, nil
}

// SendTranscriptionRequestStream sends the transcription request to the
// streaming endpoint and consumes the NDJSON event stream. For every progress
// event it invokes onProgress with a value between 0.0 and 1.0. When the
// transcription service needs to download the model first, it invokes
// onModelDownload with the model name. It returns the final WhisperResult once
// the stream is complete.
func SendTranscriptionRequestStream(t *models.Transcription, body *bytes.Buffer, writer *multipart.Writer, onProgress func(progress float64), onModelDownload func(model string)) (*models.WhisperResult, error) {
	baseUrl := fmt.Sprintf("http://%v/transcribe-stream/", os.Getenv("ASR_ENDPOINT"))

	params := url.Values{}
	params.Add("model_size", t.ModelSize)
	params.Add("task", t.Task)
	params.Add("language", t.Language)
	params.Add("device", t.Device)

	if t.BeamSize > 0 {
		params.Add("beam_size", fmt.Sprintf("%d", t.BeamSize))
	}
	if t.InitialPrompt != "" {
		params.Add("initial_prompt", t.InitialPrompt)
	}
	if len(t.Hotwords) > 0 {
		params.Add("hotwords", strings.Join(t.Hotwords, ","))
	}
	if t.VadFilter {
		params.Add("vad_filter", "true")
		if t.VadThreshold != nil {
			params.Add("vad_threshold", fmt.Sprintf("%f", *t.VadThreshold))
		}
		if t.VadMinSpeechDurationMS != nil {
			params.Add("vad_min_speech_duration_ms", fmt.Sprintf("%d", *t.VadMinSpeechDurationMS))
		}
		if t.VadMinSilenceDurationMS != nil {
			params.Add("vad_min_silence_duration_ms", fmt.Sprintf("%d", *t.VadMinSilenceDurationMS))
		}
	}

	fullUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	if err := writer.Close(); err != nil {
		log.Debug().Err(err).Msg("Error closing multipart writer")
		return nil, err
	}

	req, err := http.NewRequest("POST", fullUrl, body)
	if err != nil {
		log.Debug().Err(err).Msg("Error creating request to transcription service")
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/x-ndjson")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error sending request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		log.Error().Msgf("Response from %v: %v", fullUrl, string(b))
		return nil, errors.New("invalid status")
	}

	// streamEvent matches the NDJSON lines emitted by the transcription API.
	type streamEvent struct {
		Type     string                `json:"type"`
		Progress float64               `json:"progress"`
		Model    string                `json:"model"`
		Result   *models.WhisperResult `json:"result"`
		Error    string                `json:"error"`
	}

	var finalResult *models.WhisperResult
	scanner := bufio.NewScanner(resp.Body)
	// Allow long NDJSON lines (the final result line can be large).
	scanner.Buffer(make([]byte, 1024*1024), 64*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var ev streamEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			log.Error().Err(err).Msgf("Error decoding stream event: %v", line)
			continue
		}
		switch ev.Type {
		case "model_download":
			if onModelDownload != nil {
				onModelDownload(ev.Model)
			}
		case "progress":
			if onProgress != nil {
				onProgress(ev.Progress)
			}
		case "result":
			finalResult = ev.Result
		case "error":
			return nil, errors.New(ev.Error)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("Error reading stream")
		return nil, err
	}
	if finalResult == nil {
		return nil, errors.New("no result received from transcription stream")
	}
	return finalResult, nil
}
func CheckTranscriptionServiceHealth() (ok bool, message string) {
	url := "http://" + os.Getenv("ASR_ENDPOINT") + "/healthcheck"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false, err.Error()
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body) // Just in case, but don't fail on body read

	if resp.StatusCode == http.StatusOK {
		return true, string(body)
	}
	return false, string(body)
}
