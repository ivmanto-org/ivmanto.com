# IV-INFRA-05: Cloud Run services (frontend + backend).
#
# Ownership split with Cloud Build (cloudbuild*.yaml stay unchanged):
#   Terraform  - the service itself, identity, ingress, scaling, resources.
#   Cloud Build - the image, and for the backend the env vars and secret refs
#                 (`gcloud run deploy --set-env-vars/--set-secrets` replaces them
#                 on every deploy, with values from the trigger substitutions, #105).
# So Terraform ignores those fields after creation instead of fighting deploys.
#
# Public reachability: the nestbalm.com org enforces Domain-Restricted Sharing,
# so `allUsers` cannot be granted run.invoker (Cloud Build's
# `--allow-unauthenticated` only logs a warning). Instead the invoker IAM check
# is disabled and ingress is limited to internal + Cloud Load Balancing: the
# *.run.app URL is not publicly routable, traffic arrives only via the shared
# global LB (Cloud Armor in front). Same approach as nestbalm-app-prod.

locals {
  # Public image that serves HTTP on $PORT; replaced by the first Cloud Build deploy.
  placeholder_image = "us-docker.pkg.dev/cloudrun/container/hello"

  cloud_run_services = {
    frontend = {
      name            = "ivmanto-com" # _FRONTEND_SERVICE_NAME in cloudbuild*.yaml
      service_account = google_service_account.frontend.email
      memory          = "512Mi"
    }
    backend = {
      name            = "ivmanto-backend-service" # _BACKEND_SERVICE_NAME
      service_account = google_service_account.backend.email
      memory          = "512Mi"
    }
  }

  cloud_run_max_instances = 3
  cloud_run_cpu           = "1"
}

resource "google_cloud_run_v2_service" "this" {
  for_each = local.cloud_run_services

  name                 = each.value.name
  location             = var.region
  ingress              = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"
  invoker_iam_disabled = true
  deletion_protection  = true

  template {
    service_account = each.value.service_account

    scaling {
      min_instance_count = 0
      max_instance_count = local.cloud_run_max_instances
    }

    containers {
      image = local.placeholder_image

      ports {
        container_port = 8080
      }

      resources {
        cpu_idle = true
        limits = {
          cpu    = local.cloud_run_cpu
          memory = each.value.memory
        }
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

  lifecycle {
    ignore_changes = [
      client,
      client_version,
      template[0].containers[0].image,
      template[0].containers[0].env,
      template[0].labels,
      template[0].revision,
    ]
  }

  depends_on = [google_project_service.this]
}
