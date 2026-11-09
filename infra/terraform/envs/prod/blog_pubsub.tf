# IV-INFRA-04: Blog GCS bucket + Pub/Sub refresh pipeline.
#
# Flow: markdown published to the bucket -> GCS notification -> Pub/Sub topic
# `blog-refresh` -> push subscription -> backend endpoint
# POST /api/_internal/pubsub/blog-refresh (backend/internal/blog/handler.go)
# -> debounce -> re-scan bucket -> rebuild cache -> fire frontend rebuild webhook.
#
# Routing/authentication notes:
#   - The backend Cloud Run service has ingress INTERNAL_AND_CLOUD_LOAD_BALANCING
#     with invoker IAM disabled (DRS workaround, see cloud_run.tf), so Pub/Sub
#     cannot push to the *.run.app URL (admitted only under strict `Internal`
#     ingress). The subscription therefore pushes through the shared LB's
#     custom-domain URL, which only resolves after the DNS cut-over (#106).
#     Until then pushes will fail and retry — expected, harmless.
#   - Authentication is the backend's own shared-token check (?token=, value
#     from Secret Manager `pubsub-push-token`). That value NEVER enters
#     Terraform or state: after DNS cut-over the operator appends it to the
#     push endpoint out-of-band (see PR body) and this resource ignores
#     push-endpoint drift so Terraform never fights or leak it.

locals {
  # Must match the _GCS_BLOG_BUCKET default in cloudbuild.yaml (do-not-touch).
  blog_bucket_name = "ivmanto_com_blog_articles"
  blog_topic_name  = "blog-refresh"
}

resource "google_storage_bucket" "blog" {
  name                        = local.blog_bucket_name
  location                    = var.region
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  labels = {
    managed-by = "terraform"
    workload   = "ivmanto-com"
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.this]
}

# The backend reads published article markdown from this bucket.
resource "google_storage_bucket_iam_member" "backend_reads_blog" {
  bucket = google_storage_bucket.blog.name
  role   = "roles/storage.objectViewer"
  member = google_service_account.backend.member
}

# --- GCS -> Pub/Sub notification ---

data "google_storage_project_service_account" "gcs" {
  project = var.project_id
}

resource "google_pubsub_topic" "blog_refresh" {
  name                       = local.blog_topic_name
  message_retention_duration = "604800s" # 7 days

  depends_on = [google_project_service.this]
}

# The GCS service account is the identity that publishes notification messages.
resource "google_pubsub_topic_iam_member" "gcs_notifies" {
  topic  = google_pubsub_topic.blog_refresh.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:${data.google_storage_project_service_account.gcs.email_address}"
}

resource "google_storage_notification" "blog" {
  bucket         = google_storage_bucket.blog.name
  payload_format = "JSON_API_V1"
  topic          = google_pubsub_topic.blog_refresh.id
  event_types    = ["OBJECT_FINALIZE", "OBJECT_DELETE"]

  depends_on = [google_pubsub_topic_iam_member.gcs_notifies]
}

# --- Push subscription to the backend ---

variable "blog_push_endpoint" {
  description = "HTTPS URL Pub/Sub pushes blog events to. Must go through the shared LB (public ivmanto.com URL); the *.run.app URL is not admitted under internal-and-cloud-load-balancing ingress."
  type        = string
  default     = "https://ivmanto.com/api/_internal/pubsub/blog-refresh"
}

resource "google_pubsub_subscription" "blog_refresh_push" {
  name  = "blog-refresh-push"
  topic = google_pubsub_topic.blog_refresh.id

  ack_deadline_seconds = 20

  # Never expire; a topic with no live subscription silently drops events.
  expiration_policy {
    ttl = ""
  }

  push_config {
    push_endpoint = var.blog_push_endpoint

    # No OIDC token: Cloud Run invoker IAM is disabled on this org, so the
    # backend validates the shared token in the query string itself.
  }

  # The operator appends the ?token=<pubsub-push-token> query param out-of-band
  # once DNS is live (#106). Values from Secret Manager never enter state.
  lifecycle {
    ignore_changes = [push_config[0].push_endpoint]
  }

  depends_on = [google_project_service.this]
}
