package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/vertexai/genai"
	"github.com/joho/godotenv"
	"ivmanto.com/backend/internal/advisor"
	"ivmanto.com/backend/internal/analytics"
	"ivmanto.com/backend/internal/articles"
	"ivmanto.com/backend/internal/blog"
	"ivmanto.com/backend/internal/booking"
	"ivmanto.com/backend/internal/config"
	"ivmanto.com/backend/internal/contact"
	"ivmanto.com/backend/internal/email"
	"ivmanto.com/backend/internal/gcal"
	"ivmanto.com/backend/internal/ideas"
	"ivmanto.com/backend/internal/leads"
	"ivmanto.com/backend/internal/middleware"
)

func main() {
	// Load .env file for local development.
	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file not found, loading config from environment")
	}

	// 1. Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger) // Set as default for convenience

	// Get the port from the environment variable.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
		slog.Info("Defaulting to port", "number", port)
	}

	// 2. Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("could not load config", "error", err)
		os.Exit(1)
	}

	// 3. Initialize services
	emailService := email.NewSmtpService(&cfg.Email, logger)
	ctx := context.Background()

	// For GCal, we use Application Default Credentials (ADC).
	// On Cloud Run, this uses the attached service account's identity.
	// For local development, run `gcloud auth application-default login`.

	gcalSvc, err := gcal.NewService(ctx, cfg)
	if err != nil {
		slog.Error("Failed to create Google Calendar service", "error", err)
		os.Exit(1)
	}

	// Initialize Vertex AI client once.
	genaiClient, err := genai.NewClient(ctx, cfg.GCP.ProjectID, cfg.GCP.Location)
	if err != nil {
		slog.Error("Failed to create Vertex AI client", "error", err)
		os.Exit(1)
	}
	defer genaiClient.Close()

	// Initialize GCS client for blog storage.
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		slog.Error("Failed to create GCS storage client", "error", err)
		os.Exit(1)
	}
	defer storageClient.Close()

	// Initialize blog pipeline: GCS storage -> markdown parser -> in-memory cache.
	blogStorage := blog.NewStorage(storageClient, cfg.Blog.GCSBucket, logger)
	blogParser := blog.NewParser()
	blogCache, err := blog.NewCache(ctx, blogStorage, blogParser, logger)
	if err != nil {
		slog.Error("Failed to initialize blog cache", "error", err)
		os.Exit(1)
	}
	defer blogCache.Stop()

	// Initialize Analytics Tracker. This requires GA_API_SECRET and GA_MEASUREMENT_ID env vars.
	trackerSvc, err := analytics.NewTracker(cfg.Analytics.ApiSecret, cfg.Analytics.MeasurementID, logger)
	if err != nil {
		slog.Error("Failed to create analytics tracker", "error", err)
		os.Exit(1)
	}

	// 4. Initialize handlers, passing dependencies
	contactHandler := contact.NewHandler(logger, emailService)
	leadsHandler := leads.NewHandler(logger, emailService)
	bookingHandler := booking.NewHandler(logger, gcalSvc, emailService, trackerSvc)
	ideasHandler := ideas.NewHandler(logger, genaiClient, emailService, cfg.Ideas.GenerateIdeasPromptTemplate)
	articlesHandler := articles.NewHandler(logger)
	blogHandler := blog.NewHandler(logger, blogCache, cfg.Blog.ArticlesStatusToken, cfg.Blog.PubSubPushSAEmail, cfg.Blog.PubSubPushAudience, cfg.Blog.FrontendRebuildWebhookURL)

	// AI service advisor (#112): load the shared catalogue and build the
	// Mistral provider when MISTRAL_API_KEY is set. Everything is
	// fail-soft: any misconfiguration (missing key, missing/unparseable
	// catalogue) leaves the provider nil, so /api/advisor answers 503
	// rather than preventing the rest of the backend (booking, contact,
	// blog) from starting.
	cataloguePath := cfg.Advisor.CataloguePath
	if cataloguePath == "" {
		cataloguePath = advisor.DefaultCataloguePath
	}

	var advisorProvider advisor.Provider
	var advCatalogue *advisor.Catalogue
	if cfg.Advisor.MistralAPIKey != "" {
		var cerr error
		advCatalogue, cerr = advisor.LoadCatalogue(cataloguePath)
		if cerr != nil {
			slog.Error("MISTRAL_API_KEY set but service catalogue could not be loaded — advisor DISABLED (503)", "path", cataloguePath, "error", cerr)
		} else {
			advisorProvider, err = advisor.NewMistralProvider(advisor.MistralOptions{
				APIKey:    cfg.Advisor.MistralAPIKey,
				Model:     cfg.Advisor.MistralModel,
				Catalogue: advCatalogue,
			})
			if err != nil {
				slog.Error("could not build advisor provider — advisor DISABLED (503)", "error", err)
				advisorProvider = nil
				advCatalogue = nil
			} else {
				model := cfg.Advisor.MistralModel
				if model == "" {
					model = advisor.DefaultMistralModel
				}
				slog.Info("advisor enabled", "model", model, "catalogue", cataloguePath)
			}
		}
	} else {
		slog.Warn("MISTRAL_API_KEY unset — /api/advisor will answer 503 (not configured)")
	}
	advisorHandler := advisor.NewHandler(logger, advisorProvider, advCatalogue)

	// Security note (#115 review, LOW-1): the /api/_internal/* endpoints sit
	// behind Cloud Run with invoker_iam_disabled and the Load Balancer logging
	// enabled — the backend is the ONLY gate. Warn loudly at startup if either
	// auth knob is empty in a non-local run so a misconfigured deploy is caught
	// at boot, not discovered as a public endpoint later.
	if cfg.Blog.PubSubPushSAEmail == "" || cfg.Blog.ArticlesStatusToken == "" {
		logger.Warn("blog push/status auth DISABLED (empty PUBSUB_PUSH_SA_EMAIL or ARTICLES_STATUS_TOKEN) — endpoints are unauthenticated; set both in production",
			"pubsub_push_sa_email", cfg.Blog.PubSubPushSAEmail,
			"articles_status_token_set", cfg.Blog.ArticlesStatusToken != "")
	}

	// 5. Register routes
	mux := http.NewServeMux()
	contactHandler.RegisterRoutes(mux)
	bookingHandler.RegisterRoutes(mux)
	leadsHandler.RegisterRoutes(mux)
	ideasHandler.RegisterRoutes(mux)
	articlesHandler.RegisterRoutes(mux)
	blogHandler.RegisterRoutes(mux)
	advisorHandler.RegisterRoutes(mux)

	// 6. Apply middleware
	var finalHandler http.Handler = mux
	finalHandler = middleware.Cors(finalHandler)
	finalHandler = middleware.RequestLogger(logger, finalHandler)

	// 7. Start server
	slog.Info("Starting server", "port", port)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: finalHandler,
	}
	if err := server.ListenAndServe(); err != nil {
		slog.Error("could not start server", "error", err)
		os.Exit(1)
	}
}
