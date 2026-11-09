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

variable "rate_limit_requests_per_minute" {
  description = "Cloud Armor per-IP request budget per minute before 429."
  type        = number
  default     = 300
}

variable "lb_service_users" {
  description = "Identities that apply Nestbalm's LB Terraform (nestbalm-app-prod) and reference ivmanto's backend services."
  type        = list(string)
  default     = ["user:nick-ops@nestbalm.com"]
}
