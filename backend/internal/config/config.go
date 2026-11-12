package config

import (
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
)

// Config holds all configuration for the application.
type Config struct {
	Service   ServiceConfig
	Email     EmailConfig
	GCal      GCalConfig
	GCP       GCPConfig
	Ideas     IdeasConfig
	Analytics AnalyticsConfig
	Blog      BlogConfig
	Advisor   AdvisorConfig
}

// ServiceConfig holds configuration for the HTTP service.
type ServiceConfig struct {
	Port string
}

// EmailConfig holds configuration for the SMTP email service.
type EmailConfig struct {
	SmtpHost      string
	SmtpPort      string
	SendFrom      string
	SendFromAlias string
	SmtpPass      string // Loaded from Secret Manager
	LeadTo        string // Recipient of advisor leads (LEADS_TO); default SendFrom
}

// GCalConfig holds configuration for the Google Calendar service.
type GCalConfig struct {
	CalendarID           string
	AvailableSlotSummary string
	ServiceAccountEmail  string // Runtime SA email, used as the impersonation source for DWD.
	ImpersonateUser      string // Workspace user to impersonate via Domain-Wide Delegation.
}

// GCPConfig holds project-level Google Cloud configuration.
type GCPConfig struct {
	ProjectID string
	Location  string
}

// Add a new struct for Ideas configuration
type IdeasConfig struct {
	GenerateIdeasPromptTemplate string `env:"GENERATE_IDEAS_PROMPT_TEMPLATE"`
}

// BlogConfig holds configuration for the GCS-backed blog system.
type BlogConfig struct {
	GCSBucket                 string
	ArticlesStatusToken       string // Optional shared secret for the operator-only status endpoint (Bearer header)
	PubSubPushSAEmail         string // Optional expected OIDC signer email for Pub/Sub pushes (empty = skip in local dev)
	PubSubPushAudience        string // Expected OIDC audience for Pub/Sub pushes (empty = prod default)
	FrontendRebuildWebhookURL string // Optional Cloud Build webhook URL to trigger frontend rebuilds
}

// AnalyticsConfig holds configuration for Google Analytics.
type AnalyticsConfig struct {
	ApiSecret     string `env:"GA_API_SECRET,required"`
	MeasurementID string `env:"GA_MEASUREMENT_ID,required"`
}

// AdvisorConfig holds configuration for the AI service advisor (#112).
// All fields are optional: an empty MistralAPIKey leaves the advisor
// unconfigured (POST /api/advisor answers 503) rather than failing boot —
// the advisor is a funnel enhancement, not a startup dependency.
type AdvisorConfig struct {
	MistralAPIKey string
	MistralModel  string // empty = advisor.DefaultMistralModel
	CataloguePath string // empty = the catalogue compiled into the binary
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	var missingVars []string

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		missingVars = append(missingVars, "SMTP_HOST")
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		missingVars = append(missingVars, "SMTP_PORT")
	}
	sendFrom := os.Getenv("SEND_FROM")
	if sendFrom == "" {
		missingVars = append(missingVars, "SEND_FROM")
	}
	sendFromAlias := os.Getenv("SEND_FROM_ALIAS")
	smtpPass := os.Getenv("SMTP_PASS")
	if smtpPass == "" {
		missingVars = append(missingVars, "SMTP_PASS")
	}
	// LEADS_TO is optional: when unset, advisor leads go to SEND_FROM.
	leadTo := os.Getenv("LEADS_TO")

	calendarID := os.Getenv("CALENDAR_ID")
	if calendarID == "" {
		missingVars = append(missingVars, "CALENDAR_ID")
	}
	availableSlotSummary := os.Getenv("GCAL_AVAILABLE_SLOT_SUMMARY")
	if availableSlotSummary == "" {
		missingVars = append(missingVars, "GCAL_AVAILABLE_SLOT_SUMMARY")
	}
	gcalSAEmail := os.Getenv("GCAL_SA_EMAIL")
	if gcalSAEmail == "" {
		missingVars = append(missingVars, "GCAL_SA_EMAIL")
	}
	gcalImpersonateUser := os.Getenv("GCAL_IMPERSONATE_USER")
	if gcalImpersonateUser == "" {
		missingVars = append(missingVars, "GCAL_IMPERSONATE_USER")
	}

	projectID, err := metadata.ProjectID()
	if err != nil {
		projectID = os.Getenv("GCP_PROJECT_ID")
	}
	if projectID == "" {
		missingVars = append(missingVars, "GCP_PROJECT_ID")
	}

	location := os.Getenv("GCP_LOCATION")
	if location == "" {
		missingVars = append(missingVars, "GCP_LOCATION")
	}

	generateIdeasPromptTemplate := os.Getenv("GENERATE_IDEAS_PROMPT_TEMPLATE")

	// Load Blog config
	gcsBlogBucket := os.Getenv("GCS_BLOG_BUCKET")
	if gcsBlogBucket == "" {
		missingVars = append(missingVars, "GCS_BLOG_BUCKET")
	}

	// Load optional blog push auth (empty = skip validation in local dev).
	// ARTICLES_STATUS_TOKEN: shared secret for the operator-only status
	// endpoint (Bearer header). PUBSUB_PUSH_SA_EMAIL: expected OIDC signer
	// for Pub/Sub pushes. PUBSUB_PUSH_AUDIENCE: OIDC audience (defaults to
	// the prod push endpoint URL in the blog handler).
	articlesStatusToken := os.Getenv("ARTICLES_STATUS_TOKEN")
	pubsubPushSAEmail := os.Getenv("PUBSUB_PUSH_SA_EMAIL")
	pubsubPushAudience := os.Getenv("PUBSUB_PUSH_AUDIENCE")

	// Load Analytics config
	gaApiSecret := os.Getenv("GA_API_SECRET")
	if gaApiSecret == "" {
		missingVars = append(missingVars, "GA_API_SECRET")
	}
	gaMeasurementID := os.Getenv("GA_MEASUREMENT_ID")
	if gaMeasurementID == "" {
		missingVars = append(missingVars, "GA_MEASUREMENT_ID")
	}

	// Load advisor config (all optional; see AdvisorConfig).
	mistralAPIKey := os.Getenv("MISTRAL_API_KEY")
	mistralModel := os.Getenv("MISTRAL_MODEL")
	cataloguePath := os.Getenv("CATALOGUE_PATH")

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	return &Config{
		Service: ServiceConfig{Port: port},
		Email:   EmailConfig{SmtpHost: smtpHost, SmtpPort: smtpPort, SendFrom: sendFrom, SendFromAlias: sendFromAlias, SmtpPass: smtpPass, LeadTo: leadTo},
		GCal: GCalConfig{
			CalendarID:           calendarID,
			AvailableSlotSummary: availableSlotSummary,
			ServiceAccountEmail:  gcalSAEmail,
			ImpersonateUser:      gcalImpersonateUser,
		},
		GCP:   GCPConfig{ProjectID: projectID, Location: location},
		Ideas: IdeasConfig{GenerateIdeasPromptTemplate: generateIdeasPromptTemplate},
		Analytics: AnalyticsConfig{
			ApiSecret:     gaApiSecret,
			MeasurementID: gaMeasurementID,
		},
		Blog: BlogConfig{
			GCSBucket:                 gcsBlogBucket,
			ArticlesStatusToken:       articlesStatusToken,
			PubSubPushSAEmail:         pubsubPushSAEmail,
			PubSubPushAudience:        pubsubPushAudience,
			FrontendRebuildWebhookURL: os.Getenv("FRONTEND_REBUILD_WEBHOOK_URL"),
		},
		Advisor: AdvisorConfig{
			MistralAPIKey: mistralAPIKey,
			MistralModel:  mistralModel,
			CataloguePath: cataloguePath,
		},
	}, nil
}
