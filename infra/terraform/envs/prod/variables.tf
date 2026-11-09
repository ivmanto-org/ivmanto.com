variable "project_id" {
  description = "GCP project hosting ivmanto.com (nestbalm.com org, folder Prod)."
  type        = string
  default     = "ivmanto-com-prod-514215"
}

variable "region" {
  description = "Region for Cloud Run, Artifact Registry and regional resources."
  type        = string
  default     = "europe-west3"
}
