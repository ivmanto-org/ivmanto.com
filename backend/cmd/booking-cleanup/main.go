// Command booking-cleanup is the retention deletion for ivmanto.com
// booking data (Nick, 2026-11-13): name, work email and notes from
// /booking are kept 10 days and then deleted. The Google Calendar event
// is the booking record — there is no database copy — so deleting the
// event deletes the record. Cancellations need no cleanup run: the
// client PII is wiped in place at cancellation time (see
// gcal.CancelBooking).
//
// The 10 days count from the END of the call, never from the booking
// creation (a booking made today for a call in three weeks must survive
// the run; the booking window is 90 days).
//
// Intended to run daily (Cloud Run Job + Cloud Scheduler, or any host
// with ADC for the DWD impersonation chain). -dry-run lists what would
// be deleted without deleting.
//
// Logs carry event IDs only — never names, emails or notes — so the
// cleanup itself does not re-persist the data it deletes.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"ivmanto.com/backend/internal/config"
	"ivmanto.com/backend/internal/gcal"
)

// defaultRetentionDays is Nick's stated retention: 10 days, then deleted.
const defaultRetentionDays = 10

func main() {
	dryRun := flag.Bool("dry-run", false, "log which events would be deleted, delete nothing")
	flag.Parse()

	// Load .env file for local development, mirroring cmd/server.
	if err := godotenv.Load(); err != nil {
		slog.Info(".env file not found, loading config from environment")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	retentionDays := defaultRetentionDays
	if v := os.Getenv("BOOKING_RETENTION_DAYS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			slog.Error("invalid BOOKING_RETENTION_DAYS", "value", v)
			os.Exit(1)
		}
		retentionDays = n
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("could not load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	svc, err := gcal.NewService(ctx, cfg)
	if err != nil {
		slog.Error("could not create calendar service", "error", err)
		os.Exit(1)
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	slog.Info("starting booking retention cleanup",
		"retention_days", retentionDays,
		"cutoff", cutoff.UTC().Format(time.RFC3339),
		"dry_run", *dryRun)

	ids, err := svc.CleanupExpiredBookings(ctx, cutoff, *dryRun)
	if err != nil {
		slog.Error("booking retention cleanup finished with errors", "handled", len(ids), "error", err)
		os.Exit(1)
	}
	action := "deleted"
	if *dryRun {
		action = "would-delete"
	}
	slog.Info("booking retention cleanup complete", "action", action, "count", len(ids))
	if *dryRun {
		for _, id := range ids {
			fmt.Fprintf(os.Stdout, "dry-run: would delete event %s\n", id)
		}
	}
}
