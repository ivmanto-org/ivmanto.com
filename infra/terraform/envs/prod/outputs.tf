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

output "blog_bucket" {
  description = "Blog markdown bucket; must match _GCS_BLOG_BUCKET in cloudbuild.yaml."
  value       = google_storage_bucket.blog.name
}

output "blog_pubsub" {
  description = "Blog refresh pipeline coordinates (topic, push subscription, endpoint)."
  value = {
    topic        = google_pubsub_topic.blog_refresh.name
    subscription = google_pubsub_subscription.blog_refresh_push.name
    endpoint     = var.blog_push_endpoint
  }
}
