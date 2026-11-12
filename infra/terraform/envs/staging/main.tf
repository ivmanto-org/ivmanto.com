# Staging for ivmanto.com: one password-protected Cloud Run service that serves
# the generated static site (staging/server), so design and later content can be
# shown to invited groups before anything goes to production.
#
# Ownership split, as in envs/prod:
#   Terraform                  - APIs, registry, identity, secret shell, the service.
#   scripts/deploy-staging.sh  - the image (built from dev work, pushed, deployed).
#
# Access: the nestbalm.com org enforces Domain-Restricted Sharing, so allUsers
# cannot be granted run.invoker. The invoker IAM check is disabled instead and
# the container itself demands HTTP Basic auth (one user per group of viewers,
# from the staging-users secret) and refuses to start without credentials.
#
# Ingress: the org allows only internal-and-cloud-load-balancing. Nick set a
# project-level override of constraints/run.allowedIngress to "all" on
# ivmanto-com-sta (2026-11-12, by hand in Cloud Shell, not in Terraform). If the
# org policies are managed as code elsewhere, the override belongs there too.
#
# No backend here yet: /api/* answers 404 until the content phase adds it.

locals {
  apis = [
    "artifactregistry.googleapis.com",
    "iam.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com",
  ]

  # Public image that serves HTTP on $PORT; replaced by the first deploy.
  placeholder_image = "us-docker.pkg.dev/cloudrun/container/hello"

  labels = {
    managed-by = "terraform"
    workload   = "ivmanto-com-staging"
  }
}

resource "google_project_service" "this" {
  for_each = toset(local.apis)

  service            = each.value
  disable_on_destroy = false
}

resource "google_artifact_registry_repository" "staging" {
  location      = var.region
  repository_id = "staging"
  format        = "DOCKER"
  description   = "Images of the password-protected staging site"
  labels        = local.labels

  cleanup_policies {
    id     = "keep-last-10"
    action = "KEEP"
    most_recent_versions {
      keep_count = 10
    }
  }
  cleanup_policies {
    id     = "delete-older"
    action = "DELETE"
    condition {
      older_than = "604800s" # 7 days
    }
  }

  depends_on = [google_project_service.this]
}

resource "google_service_account" "site" {
  account_id   = "staging-site-sa"
  display_name = "Staging site runtime (reads the viewer credentials, nothing else)"

  depends_on = [google_project_service.this]
}

# Secret shell only: the user:password lines are added out-of-band
# (gcloud secrets versions add) and never live in Terraform or the repo.
resource "google_secret_manager_secret" "users" {
  secret_id = "staging-users"
  labels    = local.labels

  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }

  depends_on = [google_project_service.this]
}

resource "google_secret_manager_secret_iam_member" "site_reads_users" {
  secret_id = google_secret_manager_secret.users.id
  role      = "roles/secretmanager.secretAccessor"
  member    = google_service_account.site.member
}

resource "google_cloud_run_v2_service" "site" {
  name                 = "ivmanto-staging"
  location             = var.region
  ingress              = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = true # the container enforces Basic auth, see header
  deletion_protection  = false
  labels               = local.labels

  template {
    service_account = google_service_account.site.email

    scaling {
      min_instance_count = 0
      max_instance_count = 2
    }

    containers {
      image = local.placeholder_image

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "256Mi"
        }
        cpu_idle = true
      }

      env {
        name = "STAGING_USERS"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.users.secret_id
            version = "latest"
          }
        }
      }
    }
  }

  lifecycle {
    # The deploy script owns the image and the labels gcloud adds on deploy.
    ignore_changes = [
      client,
      client_version,
      template[0].containers[0].image,
      template[0].labels,
      template[0].revision,
    ]
  }

  depends_on = [google_secret_manager_secret_iam_member.site_reads_users]
}
