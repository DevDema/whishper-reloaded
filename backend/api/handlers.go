package api

import (
	"fmt"
	"os"
	"time"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"codeberg.org/pluja/whishper/models"
)

func (s *Server) handleGetAllTranscriptions(c *fiber.Ctx) error {
	transcriptions := s.Db.GetAllTranscriptions()

	// Convert the transcriptions to JSON.
	json, err := json.Marshal(transcriptions)
	if err != nil {
		// 503 On vacation!
		return fiber.NewError(fiber.StatusServiceUnavailable, "On vacation!")
	}

	// Write the JSON to the response body.
	c.Set("Content-Type", "application/json")
	c.Write(json)
	return nil
}

func (s *Server) handleGetTranscriptionById(c *fiber.Ctx) error {
	id := c.Params("id")
	t := s.Db.GetTranscription(id)
	if t == nil {
		log.Warn().Msgf("Transcription with id %v not found", id)
		return fiber.NewError(fiber.StatusNotFound, "Not found")
	}

	// Convert the transcription to JSON.
	json, err := json.Marshal(t)
	if err != nil {
		// 503 On vacation!
		return fiber.NewError(fiber.StatusServiceUnavailable, "On vacation!")
	}

	// Write the JSON to the response body.
	c.Set("Content-Type", "application/json")
	c.Write(json)
	return nil
}

// This function receives data from a form to create a new transcription.
// If the transcription is created successfully, it returns a 201 Created status code and
// broadcasts the new transcription to all ws clients.
func (s *Server) handlePostTranscription(c *fiber.Ctx) error {
	log.Debug().Msg("POST /api/transcriptions")
	var transcription models.Transcription

	// we get the filename from the from
	var filename string
	if c.FormValue("sourceUrl") == "" {
		// Get the form file from the request.
		file, err := c.FormFile("file")
		if err != nil {
			log.Error().Err(err).Msg("Error getting file field from the form")
			return fiber.NewError(fiber.StatusBadRequest, "Bad request")
		}
		timeid := time.Now().Format("2006_01_02-150405000")
		filename = timeid + models.FileNameSeparator + file.Filename
		// if it's empty and there is no sourceurl we set a timestamp-based filename
		if filename == timeid+models.FileNameSeparator {
			filename = timeid + models.FileNameSeparator + time.Now().Format("2006_01_02-150405")
		}

		// Save the file to the uploads directory.
		err = c.SaveFile(file, fmt.Sprintf("%v/%v", os.Getenv("UPLOAD_DIR"), filename))
		if err != nil {
			log.Error().Err(err).Msgf("Error saving the form file to disk into %v", os.Getenv("UPLOAD_DIR"))
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	// Parse the body into the transcription struct.
	transcription.Language = c.FormValue("language")
	transcription.ModelSize = c.FormValue("modelSize")
	transcription.FileName = filename
	transcription.Status = models.TranscriptionStatusPending
	transcription.Task = "transcribe"
	transcription.SourceUrl = c.FormValue("sourceUrl")
	transcription.Device = c.FormValue("device")
	if transcription.Device != "cpu" && transcription.Device != "cuda" {
		log.Warn().Msgf("Device %v not supported, using cpu", transcription.Device)
		transcription.Device = "cpu"
	}
	// Parse new params
	if c.FormValue("beam_size") != "" {
		var beamSize int
		_, err := fmt.Sscanf(c.FormValue("beam_size"), "%d", &beamSize)
		if err == nil {
			transcription.BeamSize = beamSize
		}
	}
	transcription.InitialPrompt = c.FormValue("initial_prompt")
	hotwordsStr := c.FormValue("hotwords")
	if hotwordsStr != "" {
		// Split by comma and trim spaces
		var hotwords []string
		for _, hw := range SplitAndTrim(hotwordsStr, ",") {
			if hw != "" {
				hotwords = append(hotwords, hw)
			}
		}
		transcription.Hotwords = hotwords
	}

	log.Debug().Msgf("Transcription: %+v", transcription)
	// Save transcription to database
	res, err := s.Db.NewTranscription(&transcription)
	if err != nil {
		log.Error().Err(err).Msg("Error saving transcription to database")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Broadcast transcription to websocket clients
	s.BroadcastTranscription(res)
	s.NewTranscriptionCh <- true
	
	// Convert the transcription to JSON.
	json, err := json.Marshal(res)
	if err != nil {
		// 503 On vacation!
		return fiber.NewError(fiber.StatusServiceUnavailable, "On vacation!")
	}

	// Write the JSON to the response body.
	c.Set("Content-Type", "application/json")
	c.Write(json)
	return nil
}

// SplitAndTrim splits a string by sep and trims spaces from each part
func SplitAndTrim(s, sep string) []string {
   var out []string
   for _, part := range strings.Split(s, sep) {
	   trimmed := strings.TrimSpace(part)
	   out = append(out, trimmed)
   }
   return out
}

func (s *Server) handleDeleteTranscription(c *fiber.Ctx) error {
	// First get the transcription from the database
	id := c.Params("id")
	t := s.Db.GetTranscription(id)
	if t == nil {
		log.Warn().Msgf("Transcription with id %v not found", id)
		return fiber.NewError(fiber.StatusNotFound, "Not found")
	}

	// Then delete the file from disk
	err := os.Remove(fmt.Sprintf("%v/%v", os.Getenv("UPLOAD_DIR"), t.FileName))
	if err != nil {
		log.Error().Err(err).Msgf("Error deleting file %v", t.FileName)
	}

	// Finally delete the transcription from the database
	err = s.Db.DeleteTranscription(id)
	if err != nil {
		log.Error().Err(err).Msgf("Error deleting transcription %v", id)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Return status deleted
	c.Status(fiber.StatusOK)
	return nil
}

func (s *Server) handlePatchTranscription(c *fiber.Ctx) error {
	var transcription models.Transcription
	// Parse the body into the transcription struct.
	err := json.Unmarshal(c.Body(), &transcription)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing JSON body")
		return fiber.NewError(fiber.StatusBadRequest, "Bad request")
	}

	// Update the transcription in the database
	ut, err := s.Db.UpdateTranscription(&transcription)
	if err != nil {
		log.Error().Err(err).Msgf("Error updating transcription")
		if err.Error() == "no documents were modified" {
			return fiber.NewError(fiber.StatusNotModified, "Not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Write the JSON to the response body.
	s.BroadcastTranscription(ut)

	// Return status ok
	json, err := json.Marshal(&ut)
	if err != nil {
		// 503 On vacation!
		return fiber.NewError(fiber.StatusInternalServerError, "Error parsing json!")
	}

	c.Status(fiber.StatusOK)
	c.Write(json)
	return nil
}

func (s *Server) handleRenameFile(c *fiber.Ctx) error {
	id := c.Params("id")
	newFileName := c.FormValue("newFileName")
	
	if newFileName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "New file name is required")
	}

	// Get the transcription from db
	transcription := s.Db.GetTranscription(id)
	if transcription == nil {
		return fiber.NewError(fiber.StatusNotFound, "Transcription not found")
	}

	// Split current filename to get timeid part
	parts := strings.Split(transcription.FileName, models.FileNameSeparator)
	if len(parts) < 2 {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid filename format")
	}
	timeid := parts[0]

	// Create new filename
	newFullFileName := timeid + models.FileNameSeparator + newFileName
	oldPath := fmt.Sprintf("%v/%v", os.Getenv("UPLOAD_DIR"), transcription.FileName)
	newPath := fmt.Sprintf("%v/%v", os.Getenv("UPLOAD_DIR"), newFullFileName)

	// Rename the file on disk
	err := os.Rename(oldPath, newPath)
	if err != nil {
		log.Error().Err(err).Msgf("Error renaming file from %v to %v", oldPath, newPath)
		return fiber.NewError(fiber.StatusInternalServerError, "Error renaming file")
	}

	// Update filename in database
	transcription.FileName = newFullFileName
	updatedTranscription, err := s.Db.UpdateTranscription(transcription)
	if err != nil {
		// Try to revert the file rename as db update failed
		os.Rename(newPath, oldPath)
		log.Error().Err(err).Msg("Error updating filename in database")
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating database")
	}

	// Broadcast the change to websocket clients
	s.BroadcastTranscription(updatedTranscription)

	// Return the updated transcription
	return c.JSON(updatedTranscription)
}

func (s *Server) handleTranslate(c *fiber.Ctx) error {
	id := c.Params("id")
	targetLang := c.Params("target")

	transcription := s.Db.GetTranscription(id)

	// Set status as translating
	transcription.Status = models.TrannscriptionStatusTranslating
	s.Db.UpdateTranscription(transcription)
	s.BroadcastTranscription(transcription)

	err := transcription.Translate(targetLang)
	if err != nil {
		log.Debug().Err(err).Msg("Error with translation")
		return err
	}

	// Set as done
	transcription.Status = models.TranscriptionStatusDone
	s.Db.UpdateTranscription(transcription)
	s.BroadcastTranscription(transcription)
	return nil
}

func (s *Server) handleUploadJSON(c *fiber.Ctx) error {
	var request struct {
		TranscriptionId string      `json:"transcriptionId"`
		Result         interface{} `json:"result"`
	}

	// Parse the JSON body
	err := json.Unmarshal(c.Body(), &request)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing JSON body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON format")
	}

	// Validate required fields
	if request.TranscriptionId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "transcriptionId is required")
	}

	if request.Result == nil {
		return fiber.NewError(fiber.StatusBadRequest, "result is required")
	}

	// Get the transcription from the database
	transcription := s.Db.GetTranscription(request.TranscriptionId)
	if transcription == nil {
		return fiber.NewError(fiber.StatusNotFound, "Transcription not found")
	}

	// Validate the JSON structure
	resultJSON, err := json.Marshal(request.Result)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid result format")
	}

	// Try to unmarshal into WhisperResult to validate structure
	var whisperResult models.WhisperResult
	err = json.Unmarshal(resultJSON, &whisperResult)
	if err != nil {
		log.Error().Err(err).Msg("Error validating JSON structure")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid transcription result format")
	}

	// Basic validation - ensure required fields exist
	if whisperResult.Language == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required field: language")
	}

	if whisperResult.Text == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required field: text")
	}

	if len(whisperResult.Segments) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required field: segments")
	}

	// Update the transcription with the new result
	transcription.Result = whisperResult
	updatedTranscription, err := s.Db.UpdateTranscription(transcription)
	if err != nil {
		if err.Error() == "no documents were modified" {
			return c.Status(fiber.StatusNotModified).JSON(fiber.Map{
				"message": "No changes were made",
			})
		}
		log.Error().Err(err).Msg("Error updating transcription in database")
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating transcription")
	}

	// Broadcast the updated transcription to websocket clients
	s.BroadcastTranscription(updatedTranscription)

	// Return the updated transcription
	return c.JSON(updatedTranscription)
}
