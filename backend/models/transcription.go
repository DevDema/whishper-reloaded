package models

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	ltr "github.com/snakesel/libretranslate"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transcription struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Status                  int                `bson:"status" json:"status"`
	Language                string             `bson:"language" json:"language"`
	ModelSize               string             `bson:"modelSize" json:"modelSize"`
	Task                    string             `bson:"task" json:"task"`
	Device                  string             `bson:"device" json:"device"`
	FileName                string             `bson:"fileName" json:"fileName"`
	SourceUrl               string             `bson:"sourceUrl" json:"sourceUrl"`
	BeamSize                int                `bson:"beam_size" json:"beam_size"`
	InitialPrompt           string             `bson:"initial_prompt" json:"initial_prompt"`
	Hotwords                []string           `bson:"hotwords" json:"hotwords"`
	VadFilter               bool               `bson:"vad_filter" json:"vad_filter"`
	VadThreshold            *float64           `bson:"vad_threshold,omitempty" json:"vad_threshold,omitempty"`
	VadMinSpeechDurationMS  *int               `bson:"vad_min_speech_duration_ms,omitempty" json:"vad_min_speech_duration_ms,omitempty"`
	VadMinSilenceDurationMS *int               `bson:"vad_min_silence_duration_ms,omitempty" json:"vad_min_silence_duration_ms,omitempty"`
	Result                  WhisperResult      `bson:"result" json:"result"`
	Translations            []Translation      `bson:"translations" json:"translations"`
	WordsCount              int                `bson:"words_count,omitempty" json:"words_count,omitempty"`
	Progress                float64            `bson:"progress,omitempty" json:"progress,omitempty"`
	DownloadingModel        bool               `bson:"downloading_model,omitempty" json:"downloadingModel,omitempty"`
}

type TranscriptionListItem struct {
	ID                      string                `json:"id"`
	Status                  int                   `json:"status"`
	Language                string                `json:"language"`
	ModelSize               string                `json:"modelSize"`
	Task                    string                `json:"task"`
	Device                  string                `json:"device"`
	FileName                string                `json:"fileName"`
	SourceUrl               string                `json:"sourceUrl"`
	BeamSize                int                   `json:"beam_size"`
	InitialPrompt           string                `json:"initial_prompt"`
	Hotwords                []string              `json:"hotwords"`
	VadFilter               bool                  `json:"vad_filter"`
	VadThreshold            *float64              `json:"vad_threshold,omitempty"`
	VadMinSpeechDurationMS  *int                  `json:"vad_min_speech_duration_ms,omitempty"`
	VadMinSilenceDurationMS *int                  `json:"vad_min_silence_duration_ms,omitempty"`
	Duration                float64               `json:"duration"`
	WordsCount              int                   `json:"words_count"`
	Progress                float64               `json:"progress,omitempty"`
	DownloadingModel        bool                  `json:"downloadingModel,omitempty"`
	Translations            []TranslationListItem `json:"translations"`
}

type TranslationListItem struct {
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage"`
	Status         int    `json:"translationStatus"`
}

func (t *Transcription) Translate(target string) error {
	for _, translation := range t.Translations {
		if translation.TargetLanguage == target {
			log.Debug().Msgf("Translation for %v already exists!", target)
			return fmt.Errorf("translation for %v already exists", target)
		}
	}

	translate := ltr.New(ltr.Config{
		Url: fmt.Sprintf("http://%v", os.Getenv("TRANSLATION_ENDPOINT")),
	})

	var translation Translation
	translation.SourceLanguage = t.Language
	translation.TargetLanguage = target

	trtext, err := translate.Translate(t.Result.Text, translation.SourceLanguage, translation.TargetLanguage)
	if err != nil {
		log.Debug().Err(err).Msgf("Error translating text...")
		return err
	}
	translatedText := trtext

	translatedSegments := make([]Segment, len(t.Result.Segments))
	copy(translatedSegments, t.Result.Segments)
	for i, seg := range t.Result.Segments {
		trtext, err := translate.Translate(seg.Text, translation.SourceLanguage, translation.TargetLanguage)
		if err != nil {
			log.Debug().Err(err).Msgf("Error translating segment text...")
			return err
		}
		translatedSegments[i].Text = trtext
		// Word-level data is lost, since we can't make sure that words will be in the same order and number as the final translation.
		// For example, if we translate "The big home" to Spanish, we could get "La casa grande", thus words changed order.
		translatedSegments[i].Words = []Word{}
	}

	translation.Result.Text = translatedText
	translation.Result.Segments = translatedSegments
	translation.Result.Language = translation.TargetLanguage
	t.Translations = append(t.Translations, translation)
	return nil
}
