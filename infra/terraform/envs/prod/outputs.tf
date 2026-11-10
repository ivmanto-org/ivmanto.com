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

output "cloud_run_services" {
  description = "Cloud Run service names, keyed by role; the serverless NEGs (#103) point at these."
  value       = { for k, s in google_cloud_run_v2_service.this : k => s.name }
}

output "lb_backend_services" {
  description = "Backend service self-links for Nestbalm's URL map host rules (#104)."
  value       = { for k, v in google_compute_backend_service.this : k => v.self_link }
}

output "cloud_build_sa_email" {
  value = google_service_account.cloud_build.email
}

output "frontend_rebuild_trigger_id" {
  description = "Webhook trigger id for the frontend-rebuild URL (#105); null until triggers are enabled."
  value       = one(google_cloudbuild_trigger.frontend_rebuild[*].trigger_id)
}
