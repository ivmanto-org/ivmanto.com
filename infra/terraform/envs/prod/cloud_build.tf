# IV-INFRA-08 (#105): CI/CD stays on Cloud Build with cloudbuild*.yaml
# unchanged. This recreates what the old project had:
#   - push to main         -> cloudbuild.yaml (frontend + backend)
#   - webhook (blog publish) -> cloudbuild-frontend-only.yaml
# Both run as a dedicated build SA. Triggers use a Cloud Build 2nd-gen
# GitHub connection, which a human authorises once (GitHub App OAuth), so it
# is created outside Terraform and referenced by name.
#
# Rollout (two applies):
#   1. Apply with enable_cloud_build_triggers = false (default): the build
#      SA, its roles, and the Cloud Build agent's rights.
#   2. Nick creates and authorises the connection (README "Cloud Build"):
#        gcloud builds connections create github ivmanto-github \
#          --region=europe-west3 --project=ivmanto-com-prod-514215
#   2a. Nick opens the printed link and installs the Cloud Build GitHub App
#       on ivmanto-org/ivmanto.com. This happened on 2026-11-10 (app
#       installation 180646944) — see
#       infra/terraform/manual-grants-2026-11-10.md (#125).
#   3. Apply with enable_cloud_build_triggers = true: repository link and
#      both triggers.
#   4. Webhook URL: create an API key restricted to Cloud Build and add a new
#      version of secret frontend-rebuild-webhook-url (README). The value
#      never enters Terraform.

locals {
  cloud_build_agent = "serviceAccount:service-${data.google_project.this.number}@gcp-sa-cloudbuild.iam.gserviceaccount.com"
  github_connection = "projects/${var.project_id}/locations/${var.region}/connections/${var.github_connection_name}"

  # Overrides for cloudbuild.yaml's defaults, which still name the old
  # project. The file itself stays unchanged.
  main_build_substitutions = {
    _BACKEND_SA_NAME = google_service_account.backend.email
    _GCAL_SA_EMAIL   = google_service_account.gcal.email
    _GCP_PROJECT_ID  = var.project_id
  }
}

data "google_project" "this" {}

resource "google_service_account" "cloud_build" {
  account_id   = "ivmanto-cloudbuild-sa"
  display_name = "ivmanto.com Cloud Build"
  description  = "Runs the ivmanto.com build triggers: builds, pushes and deploys both Cloud Run services."
}

resource "google_project_iam_member" "cloud_build" {
  for_each = toset([
    "roles/run.developer",     # gcloud run deploy
    "roles/logging.logWriter", # CLOUD_LOGGING_ONLY build logs
  ])

  project = var.project_id
  role    = each.value
  member  = google_service_account.cloud_build.member
}

resource "google_artifact_registry_repository_iam_member" "cloud_build_writer" {
  location   = google_artifact_registry_repository.images.location
  repository = google_artifact_registry_repository.images.name
  role       = "roles/artifactregistry.writer"
  member     = google_service_account.cloud_build.member
}

# Deploying a revision that runs as these SAs needs actAs on them.
resource "google_service_account_iam_member" "cloud_build_acts_as" {
  for_each = {
    frontend = google_service_account.frontend.name
    backend  = google_service_account.backend.name
  }

  service_account_id = each.value
  role               = "roles/iam.serviceAccountUser"
  member             = google_service_account.cloud_build.member
}

# The Cloud Build service agent stores the GitHub OAuth token of the
# connection in Secret Manager (Google's documented requirement for 2nd-gen
# connections) and reads the webhook trigger's secret.
resource "google_project_iam_member" "cloud_build_agent_secrets" {
  project = var.project_id
  role    = "roles/secretmanager.admin"
  member  = local.cloud_build_agent
}

resource "google_cloudbuildv2_repository" "site" {
  count = var.enable_cloud_build_triggers ? 1 : 0

  location          = var.region
  name              = "ivmanto.com"
  parent_connection = local.github_connection
  remote_uri        = "https://github.com/ivmanto-org/ivmanto.com.git"
}

resource "google_cloudbuild_trigger" "main_deploy" {
  count = var.enable_cloud_build_triggers ? 1 : 0

  location        = var.region
  name            = "ivmanto-main-deploy"
  description     = "Push to main: build and deploy frontend + backend (cloudbuild.yaml)."
  service_account = google_service_account.cloud_build.id
  filename        = "cloudbuild.yaml"
  substitutions   = local.main_build_substitutions

  repository_event_config {
    repository = google_cloudbuildv2_repository.site[0].id
    push {
      branch = "^main$"
    }
  }
}

data "google_secret_manager_secret_version" "webhook" {
  secret = google_secret_manager_secret.this["webhook-secret"].id
}

resource "google_cloudbuild_trigger" "frontend_rebuild" {
  count = var.enable_cloud_build_triggers ? 1 : 0

  location        = var.region
  name            = "ivmanto-frontend-rebuild"
  description     = "Webhook from the backend after a blog publish: rebuild the SSG frontend (cloudbuild-frontend-only.yaml)."
  service_account = google_service_account.cloud_build.id

  webhook_config {
    secret = data.google_secret_manager_secret_version.webhook.name
  }

  source_to_build {
    repository = google_cloudbuildv2_repository.site[0].id
    ref        = "refs/heads/main"
    repo_type  = "GITHUB"
  }

  git_file_source {
    path       = "cloudbuild-frontend-only.yaml"
    repository = google_cloudbuildv2_repository.site[0].id
    revision   = "refs/heads/main"
    repo_type  = "GITHUB"
  }
}
