package api

import (
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"

	"codeberg.org/pluja/whishper/database"
	"codeberg.org/pluja/whishper/models"
	"codeberg.org/pluja/whishper/utils"
)

type Server struct {
	ListenAddr         string
	Router             *fiber.App
	Db                 database.Db
	NewTranscriptionCh chan bool
	clients            []*websocket.Conn
}

func NewServer(listenAddr string, db database.Db) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Router: fiber.New(fiber.Config{
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
			BodyLimit:    100000 * 1024 * 1024, // Increase body limit to 100000MB (100GB)
			ServerHeader: "Fiber",              // Optional, for easier debugging
		}),
		Db:                 db,
		clients:            make([]*websocket.Conn, 0),
		NewTranscriptionCh: make(chan bool, 100),
	}
}

func (s *Server) Run() {
	s.SetupWebsocket()
	s.SetupMiddleware()
	s.RegisterRoutes()
	s.Router.Listen(s.ListenAddr)
}

func (s *Server) SetupWebsocket() {
	s.Router.Get("/ws/transcriptions", websocket.New(func(c *websocket.Conn) {

		// Add this connection to the slice of clients
		s.clients = append(s.clients, c)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				// Check for normal close error (1000) or going away error (1001)
				log.Debug().Err(err).Msgf("Error received as a warning!")
				if err.Error() != "websocket: close 1000 (normal)" &&
					err.Error() != "websocket: close 1001 (going away)" {
					log.Debug().Err(err).Msgf("Error reading message")
				}
				// Remove the client from the slice if it has disconnected
				s.clients = removeWsClient(s.clients, c)
				return
			}
			s.handleWebsocketMessage(c, msg)
		}
	}))
}

func (s *Server) BroadcastTranscription(t *models.Transcription) {
	// Convert the transcription to JSON.
	json, err := json.Marshal(&t)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling transcription to JSON:")
		return
	}
	for _, client := range s.clients {
		if err := client.WriteMessage(websocket.TextMessage, json); err != nil {
			log.Error().Err(err).Msg("Error broadcasting message:")
		}
	}
}

func (s *Server) SetupMiddleware() {
	s.Router.Use(cors.New())
}

func (s *Server) RegisterRoutes() {
	// Static routes
	s.Router.Static("/api/video", os.Getenv("UPLOAD_DIR"))

	// Register HTTP route for getting initial state.
	s.Router.Get("/api/transcriptions", func(c *fiber.Ctx) error {
		err := s.handleGetAllTranscriptions(c)
		if err != nil {
			log.Error().Err(err).Msg("Error handling POST /api/transcriptions")
		}
		return err
	})

	// Register HTTP route for getting initial state.
	s.Router.Get("/api/transcriptions/:id", func(c *fiber.Ctx) error {
		log.Debug().Msgf("GET /api/transcriptions/%v", c.Params("id"))
		err := s.handleGetTranscriptionById(c)
		if err != nil {
			log.Error().Err(err).Msg("Error handling GET /api/transcriptions/:id")
		}
		return err
	})

	// Register HTTP route for getting initial state.
	s.Router.Get("/api/translate/:id/:target", func(c *fiber.Ctx) error {
		log.Debug().Msgf("GET /api/translate/%v/%v", c.Params("id"), c.Params("target"))
		err := s.handleTranslate(c)
		if err != nil {
			log.Error().Err(err).Msg("Error handling GET /api/translate/:id/:source")
		}
		return err
	})

	// Register HTTP route for renaming files
	s.Router.Post("/api/rename/:id", func(c *fiber.Ctx) error {
		log.Debug().Msgf("POST /api/rename/%v", c.Params("id"))
		err := s.handleRenameFile(c)
		if err != nil {
			log.Error().Err(err).Msg("Error handling POST /api/rename/:id")
		}
		return err
	})

	// Register HTTP route for receiving the form data and creating new transcription job.
	s.Router.Post("/api/transcriptions", func(c *fiber.Ctx) error {
		log.Debug().Msg("POST /api/transcriptions")
		err := s.handlePostTranscription(c)
		if err != nil {
			log.Error().Err(err).Msg("Error handling POST /api/transcriptions")
		}
		return err
	})

	s.Router.Patch("/api/transcriptions", func(c *fiber.Ctx) error {
		//log.Debug().Msgf("PATCH /api/transcriptions/%v", c.Params("id"))
		err := s.handlePatchTranscription(c)
		if err != nil {
			log.Error().Err(err).Msg("Error handling PATCH /api/transcriptions")
		}
		return err
	})

	// Register HTTP route for receiving the form data and creating new transcription job.
	s.Router.Delete("/api/transcriptions/:id", func(c *fiber.Ctx) error {
		log.Debug().Msgf("DELETE /api/transcriptions/%v", c.Params("id"))
		err := s.handleDeleteTranscription(c)
		if err != nil {
			log.Error().Err(err).Msg("Error handling DELETE /api/transcriptions")
		}
		return err
	})

	// Register HTTP route for uploading JSON to replace transcription result
	s.Router.Post("/api/upload", func(c *fiber.Ctx) error {
		log.Debug().Msg("POST /api/upload")
		err := s.handleUploadJSON(c)
		if err != nil {
			log.Error().Err(err).Msg("Error handling POST /api/upload")
		}
		return err
	})

	s.Router.Get("/api/status", func(c *fiber.Ctx) error {
		healthy, msg := utils.CheckTranscriptionServiceHealth()
		if healthy {
			return c.JSON(fiber.Map{
				"status": "ok",
				"service_message": msg,
			})
		}

		// If the health check failed, it may be because the transcription-api is busy
		// processing a running transcription and not responding. Check the DB for
		// running transcriptions and fall back to reporting a likely running state.
		running := s.Db.GetRunningTranscription()
		if running != nil && len(running) > 0 {
			log.Debug().Msgf("Transcription service healthcheck failed but %d running transcriptions found", len(running))
			return c.JSON(fiber.Map{
				"status": "ok",
				"service_message": "transcription service unreachable but there are running transcriptions",
			})
		}

		// No running transcriptions -> real outage
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "error",
			"error":  "transcription service unavailable",
			"service_message": msg,
		})
	})
}

// Helper function to remove a WebSocket connection from the slice
func removeWsClient(s []*websocket.Conn, r *websocket.Conn) []*websocket.Conn {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
