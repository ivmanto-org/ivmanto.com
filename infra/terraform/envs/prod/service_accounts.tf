# Runtime identities. No service account keys are ever created.
#
#   frontend  - Nuxt SSG container; needs no Google API access.
#   backend   - Go API (ivmanto-backend-sa, name kept from cloudbuild.yaml).
#   gcal      - Domain-Wide Delegation identity for Google Calendar. The backend
#               impersonates it (GCAL_SA_EMAIL), and it impersonates the Workspace
#               user (GCAL_IMPERSONATE_USER). Its OAuth client ID must be authorised
#               in the ivmanto.com Workspace Admin Console (IV-INFRA-10, #107).

resource "google_service_account" "frontend" {
  account_id   = "ivmanto-frontend-sa"
  display_name = "ivmanto.com frontend (Cloud Run)"

  depends_on = [google_project_service.this]
}

resource "google_service_account" "backend" {
  account_id   = "ivmanto-backend-sa"
  display_name = "ivmanto.com backend (Cloud Run)"

  depends_on = [google_project_service.this]
}

resource "google_service_account" "gcal" {
  account_id   = "ivmanto-gcal-dwd-sa"
  display_name = "ivmanto.com Google Calendar DWD"
  description  = "Impersonated by the backend; authorised for DWD in the ivmanto.com Workspace."

  depends_on = [google_project_service.this]
}

# blog_push - OIDC signer identity for the blog-refresh Pub/Sub push
# subscription (blog_gcs.tf, #115). Pub/Sub attaches this SA's OIDC token to
# each push request (Authorization: Bearer header); the backend verifies the
# token's signature, audience and this SA's email. No keys are ever created.
resource "google_service_account" "blog_push" {
  account_id   = "ivmanto-blog-push-sa"
  display_name = "ivmanto.com blog Pub/Sub push (OIDC signer)"
  description  = "Signs the OIDC bearer token the blog-refresh Pub/Sub push subscription sends to the backend (#115)."

  depends_on = [google_project_service.this]
}

# Backend -> gcal SA: mint tokens / sign the DWD JWT (impersonate.CredentialsTokenSource).
resource "google_service_account_iam_member" "backend_impersonates_gcal" {
  service_account_id = google_service_account.gcal.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = google_service_account.backend.member
}

# Backend -> Vertex AI (ideas generation).
resource "google_project_iam_member" "backend_vertex_user" {
  project = var.project_id
  role    = "roles/aiplatform.user"
  member  = google_service_account.backend.member
}

# Backend -> secrets it reads at runtime (cloudbuild.yaml --set-secrets).
# The secrets themselves are created in secrets.tf (IV-INFRA-03, #100); they must
# exist before this is applied.
locals {
  backend_secret_ids = [
    "smtp-password",
    "ideas-prompt",
    "ga-api-secret",
    "pubsub-push-token",
    "frontend-rebuild-webhook-url",
  ]
}

resource "google_secret_manager_secret_iam_member" "backend_secret_accessor" {
  for_each = toset(local.backend_secret_ids)

  project   = var.project_id
  secret_id = each.value
  role      = "roles/secretmanager.secretAccessor"
  member    = google_service_account.backend.member
}
