output "staging_url" {
  description = "Address to share with invited viewers (together with their group's password)."
  value       = google_cloud_run_v2_service.site.uri
}

output "image_repository" {
  description = "Where scripts/deploy-staging.sh pushes the staging image."
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.staging.repository_id}"
}
