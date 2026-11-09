# Bootstrap: creates the GCS bucket that holds Terraform state for envs/*.
#
# Chicken-and-egg: this config starts with local state, creates the bucket,
# then its own state is migrated into the bucket (prefix "bootstrap"):
#   terraform init && terraform apply
#   (backend block below was added afterwards, then: terraform init -migrate-state)

terraform {
  required_version = ">= 1.6.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.0"
    }
  }

  backend "gcs" {
    bucket = "ivmanto-com-prod-514215-tfstate"
    prefix = "bootstrap"
  }
}

variable "project_id" {
  type    = string
  default = "ivmanto-com-prod-514215"
}

variable "state_location" {
  type    = string
  default = "europe-west3"
}

locals {
  noncurrent_versions_to_keep = 20
}

provider "google" {
  project = var.project_id
}

resource "google_project_service" "storage" {
  service            = "storage.googleapis.com"
  disable_on_destroy = false
}

resource "google_storage_bucket" "tfstate" {
  name                        = "${var.project_id}-tfstate"
  location                    = var.state_location
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"
  force_destroy               = false

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      num_newer_versions = local.noncurrent_versions_to_keep
      with_state         = "ARCHIVED"
    }
    action {
      type = "Delete"
    }
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_project_service.storage]
}

output "state_bucket" {
  value = google_storage_bucket.tfstate.name
}
