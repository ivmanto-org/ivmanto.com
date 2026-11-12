variable "project_id" {
  description = "Staging project. Separate from prod so nothing staged can touch production."
  type        = string
  default     = "ivmanto-com-sta"
}

variable "region" {
  description = "Region for all staging resources."
  type        = string
  default     = "europe-west3"
}
