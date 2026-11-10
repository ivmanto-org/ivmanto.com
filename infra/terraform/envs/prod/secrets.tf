# IV-INFRA-03: Secret Manager secret shells (import-only).
#
# All six secrets already exist in ivmanto-com-prod-514215 (values migrated
# out-of-band by Nick, each with version 1 enabled and user-managed
# replication in europe-west3). Terraform adopts them via import blocks —
# it must NOT create, replace or destroy them:
#   - replication is pinned to user-managed europe-west3; any mismatch
#     forces replacement, which would DESTROY the stored values;
#   - lifecycle.prevent_destroy is a second line of defence;
#   - no google_secret_manager_secret_version resources: secret values are
#     added by Nick out-of-band
#     (gcloud secrets versions add <name> --data-file=-), never in Terraform.
#
# Accessor IAM bindings do NOT live here: they ride in #99 (IV-INFRA-02 /
# service_accounts.tf) because they reference google_service_account.backend.
#
# Acceptance: terraform plan = 6 to import, 0 to add, 0 to change, 0 to destroy.

locals {
  secrets = {
    smtp-password                = "SMTP password for the backend service"
    ideas-prompt                 = "LLM prompt used by the ideas-generation endpoint"
    ga-api-secret                = "Google Analytics API secret (Measurement Protocol)"
    pubsub-push-token            = "Shared bearer token for the operator-only articles status endpoint (Authorization header, constant-time compare, #115)"
    frontend-rebuild-webhook-url = "Cloud Build webhook trigger URL used by the backend after a blog publish (new version once #105 creates the trigger)"
    webhook-secret               = "Cloud Build webhook trigger secret (GitHub webhook)"
  }
}

resource "google_secret_manager_secret" "this" {
  for_each = local.secrets

  project = var.project_id

  secret_id = each.key

  labels = {
    managed-by = "terraform"
    workload   = "ivmanto-com"
  }

  replication {
    user_managed {
      replicas {
        location = "europe-west3"
      }
    }
  }

  lifecycle {
    prevent_destroy = true
  }
}

import {
  to = google_secret_manager_secret.this["smtp-password"]
  id = "projects/ivmanto-com-prod-514215/secrets/smtp-password"
}

import {
  to = google_secret_manager_secret.this["ideas-prompt"]
  id = "projects/ivmanto-com-prod-514215/secrets/ideas-prompt"
}

import {
  to = google_secret_manager_secret.this["ga-api-secret"]
  id = "projects/ivmanto-com-prod-514215/secrets/ga-api-secret"
}

import {
  to = google_secret_manager_secret.this["pubsub-push-token"]
  id = "projects/ivmanto-com-prod-514215/secrets/pubsub-push-token"
}

import {
  to = google_secret_manager_secret.this["frontend-rebuild-webhook-url"]
  id = "projects/ivmanto-com-prod-514215/secrets/frontend-rebuild-webhook-url"
}

import {
  to = google_secret_manager_secret.this["webhook-secret"]
  id = "projects/ivmanto-com-prod-514215/secrets/webhook-secret"
}
