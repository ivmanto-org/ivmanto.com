# Docker repository used by cloudbuild.yaml and cloudbuild-frontend-only.yaml
# (europe-west3-docker.pkg.dev/<project>/ivmanto-com-repo/<service>:<tag>).

locals {
  images_to_keep        = 10
  image_max_age_seconds = 30 * 24 * 60 * 60 # 30 days
}

resource "google_artifact_registry_repository" "images" {
  location      = var.region
  repository_id = "ivmanto-com-repo"
  format        = "DOCKER"
  description   = "ivmanto.com frontend and backend images (pushed by Cloud Build)"

  cleanup_policy_dry_run = false

  # Policies are evaluated together: KEEP wins over DELETE, so the newest
  # images (including the currently deployed ones) are never removed.
  cleanup_policies {
    id     = "keep-recent"
    action = "KEEP"
    most_recent_versions {
      keep_count = local.images_to_keep
    }
  }

  cleanup_policies {
    id     = "delete-old"
    action = "DELETE"
    condition {
      older_than = "${local.image_max_age_seconds}s"
    }
  }

  depends_on = [google_project_service.this]
}
