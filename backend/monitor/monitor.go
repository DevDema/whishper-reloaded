package monitor

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"codeberg.org/pluja/whishper/api"
	"codeberg.org/pluja/whishper/models"
	"codeberg.org/pluja/whishper/utils"
)

func StartMonitor(s *api.Server) {
	log.Info().Msg("Starting monitor!")
	go func() {
		for {
			// Wait for new transcription to be added to the database
			// notification will be received through the NewTranscriptionCh channel
			<-s.NewTranscriptionCh
			pendingTranscriptions := s.Db.GetPendingTranscriptions()
			log.Debug().Msgf("Pending transcriptions: %v", len(pendingTranscriptions))
			for _, pt := range pendingTranscriptions {
				log.Debug().Msgf("Taking pending transcription %v", pt.ID)
				if pt.Status == models.TranscriptionStatusPending {
					err := transcribe(s, pt)
					if err != nil {
						log.Error().Err(err).Msg("Error transcribing")
						pt.Status = models.TranscriptionStatusError
						pt.DownloadingModel = false
						ut, err := s.Db.UpdateTranscription(pt)
						if err != nil {
							log.Error().Err(err).Msg("Error updating transcription")
						}
						s.BroadcastTranscription(ut)
						continue
					}
				}
			}
		}
	}()
}

func transcribe(s *api.Server, t *models.Transcription) error {
	// Update transcription status
	t.Status = models.TranscriptionStatusRunning
	log.Debug().Msgf("Updating transcription %v", t)
	_, err := s.Db.UpdateTranscription(t)
	if err != nil {
		log.Error().Err(err).Msg("Error updating transcription")
		return err
	}
	s.BroadcastTranscription(t)

	if t.SourceUrl != "" {
		// Download media
		fn, err := utils.DownloadMedia(t)
		if err != nil {
			log.Error().Err(err).Msg("Error downloading media")
			return err
		}
		t.FileName = fn
		s.BroadcastTranscription(t)
	}

	// Prepare multipart form data
	body, writer, err := prepareMultipartFormData(t)
	if err != nil {
		log.Error().Err(err).Msg("Error preparing multipart form data")
		return err
	}

	// Send transcription request to transcription service. We use the
	// streaming endpoint so we can report progress while it runs.
	var lastBroadcast float64
	res, err := utils.SendTranscriptionRequestStream(t, body, writer, func(progress float64) {
		// Once real progress arrives the model is loaded, so clear the
		// downloading flag (force a broadcast on this transition).
		downloadingCleared := t.DownloadingModel
		// Throttle: only persist/broadcast on meaningful changes (>= 1%).
		if progress-lastBroadcast < 0.01 && progress < 1.0 && !downloadingCleared {
			return
		}
		lastBroadcast = progress
		t.Progress = progress
		t.DownloadingModel = false
		if _, uerr := s.Db.UpdateTranscription(t); uerr != nil {
			log.Error().Err(uerr).Msg("Error updating transcription progress")
			return
		}
		s.BroadcastTranscription(t)
	}, func(model string) {
		// The transcription service is downloading the model weights.
		log.Info().Msgf("Downloading model %v for transcription %v", model, t.ID.Hex())
		t.DownloadingModel = true
		t.Progress = 0
		if _, uerr := s.Db.UpdateTranscription(t); uerr != nil {
			log.Error().Err(uerr).Msg("Error updating transcription download state")
			return
		}
		s.BroadcastTranscription(t)
	})
	if err != nil {
		log.Error().Err(err).Msg("Error sending transcription request")
		return err
	}

	t.Result = *res
	t.Translations = []models.Translation{}
	t.Status = models.TranscriptionStatusDone
	t.Progress = 1.0
	t.DownloadingModel = false
	_, err = s.Db.UpdateTranscription(t)
	if err != nil {
		log.Error().Err(err).Msg("Error updating transcription")
		return err
	}
	s.BroadcastTranscription(t)
	return nil
}

func prepareMultipartFormData(t *models.Transcription) (*bytes.Buffer, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", t.FileName)
	if err != nil {
		log.Error().Err(err).Msg("Error creating form file")
		return nil, nil, err
	}

	// Read file from disk
	filePath := filepath.Join(os.Getenv("UPLOAD_DIR"), t.FileName)
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Error opening file")
		return nil, nil, err
	}
	defer file.Close()

	_, err = io.Copy(part, file)
	if err != nil {
		log.Error().Err(err).Msg("Error copying file")
		return nil, nil, err
	}

	err = writer.Close()
	if err != nil {
		log.Error().Err(err).Msg("Error closing writer")
		return nil, nil, err
	}

	return body, writer, nil
}
