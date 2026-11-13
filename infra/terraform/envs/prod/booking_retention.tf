# Booking-data retention (Nick, 2026-11-13; PR #140): name, work email and notes
# from /booking are deleted 10 days after the call. The Google Calendar event is
# the booking record, so the job wipes and deletes expired events there.
#
#   Cloud Scheduler (daily) -> Cloud Run Job `ivmanto-booking-cleanup`
#
# Ownership split, same as the services in cloud_run.tf:
#   Terraform   - the job, its identity, the schedule and who may start it.
#   Cloud Build - the image, command, env vars and secret refs (`gcloud run jobs
#                 deploy` in cloudbuild.yaml), so the job always runs the same
#                 backend image as the service.
# Until the first Cloud Build deploy the job holds Google's placeholder job
# image, which prints a line and exits 0: a scheduled run before that deploy
# deletes nothing and fails nothing.

locals {
  booking_cleanup_job_name = "ivmanto-booking-cleanup" # _BOOKING_CLEANUP_JOB_NAME in cloudbuild.yaml
  placeholder_job_image    = "us-docker.pkg.dev/cloudrun/container/job"

  # Daily, at night for the people whose data it is. Retention is counted from
  # the end of the call, so the hour only shifts the deletion by under a day.
  booking_cleanup_schedule = "15 3 * * *"
  booking_cleanup_timezone = "Europe/Berlin"
}

resource "google_cloud_run_v2_job" "booking_cleanup" {
  name                = local.booking_cleanup_job_name
  location            = var.region
  deletion_protection = true

  template {
    task_count = 1

    template {
      # Same identity as the backend: it already holds the token-creator grant
      # on the Calendar DWD account and nothing else is needed.
      service_account = google_service_account.backend.email
      max_retries     = 1
      timeout         = "600s"

      containers {
        image = local.placeholder_job_image

        resources {
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [
      client,
      client_version,
      template[0].labels,
      template[0].template[0].containers[0].image,
      template[0].template[0].containers[0].command,
      template[0].template[0].containers[0].args,
      template[0].template[0].containers[0].env,
    ]
  }

  depends_on = [google_project_service.this]
}

# Identity Cloud Scheduler uses to start the job. It can start this one job and
# nothing else. No keys are created.
resource "google_service_account" "booking_cleanup_scheduler" {
  account_id   = "ivmanto-booking-cleanup-sched"
  display_name = "ivmanto.com booking cleanup (Cloud Scheduler)"
  description  = "Starts the ivmanto-booking-cleanup Cloud Run Job once a day."

  depends_on = [google_project_service.this]
}

resource "google_cloud_run_v2_job_iam_member" "scheduler_runs_booking_cleanup" {
  name     = google_cloud_run_v2_job.booking_cleanup.name
  location = google_cloud_run_v2_job.booking_cleanup.location
  role     = "roles/run.invoker"
  member   = google_service_account.booking_cleanup_scheduler.member
}

resource "google_cloud_scheduler_job" "booking_cleanup" {
  name        = "ivmanto-booking-cleanup-daily"
  description = "Daily booking-data retention run (delete 10 days after the call)."
  region      = var.region
  schedule    = local.booking_cleanup_schedule
  time_zone   = local.booking_cleanup_timezone

  # The job itself decides success; a failed start is retried a few times.
  attempt_deadline = "180s"

  retry_config {
    retry_count          = 3
    min_backoff_duration = "60s"
    max_backoff_duration = "600s"
  }

  http_target {
    http_method = "POST"
    uri         = "https://run.googleapis.com/v2/projects/${var.project_id}/locations/${var.region}/jobs/${google_cloud_run_v2_job.booking_cleanup.name}:run"

    oauth_token {
      service_account_email = google_service_account.booking_cleanup_scheduler.email
      scope                 = "https://www.googleapis.com/auth/cloud-platform"
    }
  }

  depends_on = [
    google_project_service.this,
    google_cloud_run_v2_job_iam_member.scheduler_runs_booking_cleanup,
  ]
}
