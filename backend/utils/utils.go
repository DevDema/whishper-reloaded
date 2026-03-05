package utils

import (
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

	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		log.Debug().Err(err).Msg("Error downloading media")
		return "", err
	}

	filename := fmt.Sprintf("%v%v%v", t.ID.Hex(), models.FileNameSeparator, result.Info.Title)
	filename = SanitizeFilename(filename)

	defer downloadResult.Close()
	f, err := os.Create(filepath.Join(os.Getenv("UPLOAD_DIR"), filename))
	if err != nil {
		log.Debug().Err(err).Msg("Error creating file")
		return "", err
	}
	defer f.Close()
	io.Copy(f, downloadResult)

	return filename, nil
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
