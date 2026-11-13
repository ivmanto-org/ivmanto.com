locals {
  required_apis = [
    "aiplatform.googleapis.com",       # ideas generation (Vertex AI / Gemini)
    "apikeys.googleapis.com",          # API key on the frontend-rebuild webhook URL (#105)
    "artifactregistry.googleapis.com", # container images
    "calendar-json.googleapis.com",    # booking (Google Calendar via DWD)
    "cloudbuild.googleapis.com",       # CI/CD (cloudbuild*.yaml)
    "cloudresourcemanager.googleapis.com",
    "cloudscheduler.googleapis.com", # daily booking-retention run (booking_retention.tf)
    "compute.googleapis.com",        # serverless NEGs, backend services, Cloud Armor
    "iam.googleapis.com",
    "iamcredentials.googleapis.com", # DWD token impersonation
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "pubsub.googleapis.com", # blog publish notifications
    "run.googleapis.com",    # frontend + backend services
    "secretmanager.googleapis.com",
    "serviceusage.googleapis.com",
    "storage.googleapis.com", # blog bucket
  ]
}

resource "google_project_service" "this" {
  for_each = toset(local.required_apis)

  service            = each.value
  disable_on_destroy = false
}
