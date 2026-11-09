output "artifact_registry_repo" {
  description = "Image path prefix used by Cloud Build."
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.images.repository_id}"
}

output "frontend_sa_email" {
  value = google_service_account.frontend.email
}

output "backend_sa_email" {
  value = google_service_account.backend.email
}

output "gcal_sa_email" {
  description = "Value for GCAL_SA_EMAIL."
  value       = google_service_account.gcal.email
}

output "gcal_sa_client_id" {
  description = "OAuth client ID to authorise for DWD in the Workspace Admin Console (#107)."
  value       = google_service_account.gcal.unique_id
}
