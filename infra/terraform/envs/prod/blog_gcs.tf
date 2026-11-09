# IV-INFRA-04: blog GCS bucket + Pub/Sub refresh pipeline.
#
# Pipeline: blog authors drop .md articles (and images) into the blog GCS
# bucket -> GCS sends an object notification to the Pub/Sub topic -> the push
# subscription POSTs to the backend's blog-refresh endpoint
# (POST /api/_internal/pubsub/blog-refresh, backend/internal/blog/handler.go),
# which debounces a cache refresh and, for new/updated articles, triggers a
# frontend-only Cloud Build rebuild.
#
# Public routing: the push endpoint is served through the shared Nestbalm
# global external LB (host rule /api/* -> backend service, IV-INFRA-7),
# NOT the *.run.app URL - the Cloud Run services (IV-INFRA-05) run with
# invoker_iam_disabled and ingress internal-and-cloud-load-balancing, so the
# only admitted external traffic is LB traffic. (Whether Pub/Sub push to the
# LB URL itself passes is verified empirically on first apply; the push_config
# below is the mechanism regardless.)
#
# Push auth: because Cloud Run has invoker_iam_disabled, the OIDC token a push
# subscription normally carries is NOT validated by the environment. The
# backend is therefore the ONLY gate: it requires a shared token in the query
# string (?token=<PUBSUB_PUSH_TOKEN>, constant-time compare, no token echo).
# Following the IV-INFRA-03 secrets convention, the token VALUE never appears
# in Terraform: Terraform creates the subscription with the base endpoint and
# ignores the push endpoint, and the operator appends the token out-of-band:
#
#   gcloud pubsub subscriptions update ivmanto-blog-refresh-push \
#     --push-endpoint="https://ivmanto.com/api/_internal/pubsub/blog-refresh?token=<val>"
#
# Until that patch lands the endpoint carries no token and the backend returns
# 403 (fail-closed) - no notifications are dropped silently during the window.
#
# Acceptance: plan shows 0 import / +N add / 0 change / 0 destroy (all new).

locals {
  blog_bucket_name   = "ivmanto_com_blog_articles" # _GCS_BLOG_BUCKET in cloudbuild.yaml
  blog_topic_name    = "ivmanto-blog-refresh"
  blog_subscription  = "ivmanto-blog-refresh-push"
  blog_push_endpoint = "https://ivmanto.com/api/_internal/pubsub/blog-refresh"
}

data "google_project" "this" {
  project_id = var.project_id
}

# --- Bucket ------------------------------------------------------------
resource "google_storage_bucket" "blog" {
  name                        = local.blog_bucket_name
  project                     = var.project_id
  location                    = var.region
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.this]
}

# Backend reads the .md articles and writes the metadata.json cache
# (backend/internal/blog/gcs.go: ListMarkdownFiles/ReadFile/WriteMetadataCache).
resource "google_storage_bucket_iam_member" "backend_blog_reader" {
  bucket = google_storage_bucket.blog.name
  role   = "roles/storage.objectViewer"
  member = google_service_account.backend.member
}

resource "google_storage_bucket_iam_member" "backend_blog_cache_writer" {
  bucket = google_storage_bucket.blog.name
  role   = "roles/storage.objectCreator"
  member = google_service_account.backend.member
}

# --- Pub/Sub -----------------------------------------------------------
resource "google_pubsub_topic" "blog" {
  name    = local.blog_topic_name
  project = var.project_id

  depends_on = [google_project_service.this]
}

# GCS notifications are published by the gs-project-accounts service account,
# so it needs pubsub.publisher on the topic.
resource "google_pubsub_topic_iam_member" "blog_gcs_publisher" {
  project = var.project_id
  topic   = google_pubsub_topic.blog.name
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:service-${data.google_project.this.number}@gs-project-accounts.iam.gserviceaccount.com"

  depends_on = [google_project_service.this]
}

resource "google_storage_notification" "blog" {
  bucket         = google_storage_bucket.blog.name
  topic          = google_pubsub_topic.blog.id
  event_types    = ["OBJECT_FINALIZE", "OBJECT_DELETE"]
  payload_format = "JSON_API_V1" # backend reads message attributes eventType/objectId
}

resource "google_pubsub_subscription" "blog_push" {
  name    = local.blog_subscription
  project = var.project_id
  topic   = google_pubsub_topic.blog.id

  ack_deadline_seconds       = 60
  message_retention_duration = "86400s"

  push_config {
    # OIDC token intentionally not configured: invoker_iam_disabled means it is
    # not validated; the ?token= query guard is the auth (see header comment).
    push_endpoint = local.blog_push_endpoint
  }

  # Operator appends ?token=<secret> out-of-band (never in Terraform).
  lifecycle {
    ignore_changes = [
      push_config[0].push_endpoint,
    ]
  }

  depends_on = [google_project_service.this]
}
