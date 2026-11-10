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
# Push auth (#115): because Cloud Run has invoker_iam_disabled, the environment
# does NOT validate the OIDC token a push subscription normally carries - the
# backend is the ONLY gate. The subscription's oidc_token block attaches a
# Google-signed OIDC JWT for google_service_account.blog_push
# (Authorization: Bearer header, audience = the push endpoint). The backend
# verifies signature, issuer, audience and expiry, and checks the token's
# email claim equals the expected signer SA (PUBSUB_PUSH_SA_EMAIL). No secret
# value exists in the URL or in Terraform, so LB request logging of this
# backend is enabled (lb_backends.tf). The token value is never logged or
# echoed by the backend.
#
# One-time operator step BEFORE the first apply that adds oidc_token: grant the
# identity creating/updating the subscription (the Terraform applier) the
# Service Account Token Creator role on the push SA:
#
#   gcloud iam service-accounts add-iam-policy-binding \
#     ivmanto-blog-push-sa@ivmanto-com-prod-514215.iam.gserviceaccount.com \
#     --member="user:<applier-identity>" \
#     --role="roles/iam.serviceAccountTokenCreator"
#
# The articles status endpoint (/api/_internal/articles/status) is operator-
# only and uses a separate shared bearer token (ARTICLES_STATUS_TOKEN, read
# from the pubsub-push-token Secret Manager secret; constant-time compare,
# Authorization header, never the URL).
#
# Acceptance: plan shows 0 import / +N add / 0 change / 0 destroy (all new).

locals {
  blog_bucket_name   = "ivmanto-com-prod-514215-blog" # _GCS_BLOG_BUCKET in cloudbuild.yaml
  blog_topic_name    = "ivmanto-blog-refresh"
  blog_subscription  = "ivmanto-blog-refresh-push"
  blog_push_endpoint = "https://ivmanto.com/api/_internal/pubsub/blog-refresh"
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
# The cache write is an OVERWRITE of metadata.json (gcs.go Writer), which
# GCS implements as delete+create -> objectCreator alone fails on the 2nd
# write (needs storage.objects.delete). objectUser bundles viewer+creator
# + delete, so one binding covers read, create and overwrite.
resource "google_storage_bucket_iam_member" "backend_blog_access" {
  bucket = google_storage_bucket.blog.name
  role   = "roles/storage.objectUser"
  member = google_service_account.backend.member
}

# --- Pub/Sub -----------------------------------------------------------
resource "google_pubsub_topic" "blog" {
  name    = local.blog_topic_name
  project = var.project_id

  depends_on = [google_project_service.this]
}

# GCS notifications are published by the project's GCS service agent
# (gs-project-accounts). Resolving it via the data source (instead of a
# hard-coded service-${project.number}@gs-project-accounts... address)
# makes Terraform verify/provision the service agent exists in a fresh
# project before the IAM grant.
data "google_storage_project_service_account" "gcs" {
  project = var.project_id
}

resource "google_pubsub_topic_iam_member" "blog_gcs_publisher" {
  project = var.project_id
  topic   = google_pubsub_topic.blog.name
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${data.google_storage_project_service_account.gcs.email_address}"

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
    # OIDC push auth (#115): Pub/Sub attaches the blog_push SA's Google-signed
    # JWT as an Authorization: Bearer header (audience = push endpoint). Cloud
    # Run does not validate it (invoker_iam_disabled); the backend does - see
    # the header comment. No secret in the URL or in Terraform.
    push_endpoint = local.blog_push_endpoint

    oidc_token {
      service_account_email = google_service_account.blog_push.email
      audience              = local.blog_push_endpoint
    }
  }

  depends_on = [google_project_service.this]
}
