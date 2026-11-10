# ivmanto.com infrastructure (Terraform)

Target: GCP project `ivmanto-com-prod-514215` (nestbalm.com org, folder `Prod`).
Public traffic arrives through the **shared Nestbalm global external Application Load
Balancer** (`nestbalm-app-prod`); its URL map and certificate live in the Nestbalm repo
(`Nestbalm/nestbalm-saas`, `infra/terraform/modules/app_edge`). This repo owns everything
inside the ivmanto project. Deploys stay on Cloud Build (`cloudbuild*.yaml`).

Plan and task list: `tasks/2026-11-09-gcp-infra-terraform.md`.

## Layout

| Path | Purpose | State |
|---|---|---|
| `bootstrap/` | Terraform state bucket `ivmanto-com-prod-514215-tfstate` | same bucket, prefix `bootstrap` |
| `envs/prod/` | Everything else for production | same bucket, prefix `envs/prod` |

## Usage

```bash
gcloud auth application-default login   # operator account, see "Operator permissions"

# one-time (already done 2026-11-09): with the backend block in bootstrap/main.tf
# commented out, run init + apply, then restore the block and migrate local state
cd infra/terraform/bootstrap
terraform init && terraform apply
terraform init -migrate-state

# day-to-day
cd infra/terraform/envs/prod
terraform init
terraform plan
terraform apply
```

Secret **values** are never stored in Terraform; it adopts the existing Secret
Manager secrets via `import` blocks (replication pinned to `europe-west3`,
`prevent_destroy` on). Add versions with
`gcloud secrets versions add <name> --data-file=-`.

## Operator permissions

Owner works but is broader than needed. A least-privilege operator for `envs/prod` needs:

| Role | For |
|---|---|
| `roles/serviceusage.serviceUsageAdmin` | API enablement |
| `roles/storage.admin` | state bucket, blog bucket |
| `roles/artifactregistry.admin` | image repository |
| `roles/iam.serviceAccountAdmin` | service accounts and their IAM |
| `roles/resourcemanager.projectIamAdmin` | project-level role grants |
| `roles/secretmanager.admin` | secrets and accessor bindings |
| `roles/run.admin` | Cloud Run services |
| `roles/compute.loadBalancerAdmin`, `roles/compute.securityAdmin` | NEGs, backend services, Cloud Armor |
| `roles/pubsub.admin` | blog publish notifications |
| `roles/cloudbuild.builds.editor`, `roles/cloudbuild.connectionAdmin` | Cloud Build connection and triggers |

Narrowing the current Owner-based operator to this set is tracked as a follow-up.

## Cloud Build (#105)

CI/CD stays on Cloud Build and `cloudbuild*.yaml` is unchanged; `cloud_build.tf` recreates the triggers. Rollout:

1. Apply with the default `enable_cloud_build_triggers = false`. This creates the build SA `ivmanto-cloudbuild-sa` and its roles, and lets the Cloud Build service agent store the connection's OAuth token.
2. **Operator, once:** create and authorise the GitHub connection.
   ```bash
   gcloud builds connections create github ivmanto-github \
     --region=europe-west3 --project=ivmanto-com-prod-514215
   ```
   Open the printed link, authorise as the GitHub account that owns `ivmanto-org`, and install the Cloud Build GitHub App on `ivmanto-org/ivmanto.com`. Check the result:
   `gcloud builds connections describe ivmanto-github --region=europe-west3 --format='value(installationState.stage)'` should print `COMPLETE`.
3. Apply with `-var enable_cloud_build_triggers=true` (then set it in a tfvars or default). This creates the repository link and both triggers.
4. **Operator, once:** the frontend-rebuild webhook URL the backend calls after a blog publish. The value stays out of Terraform.
   ```bash
   KEY=$(gcloud services api-keys create --display-name="cloudbuild-webhook" \
     --api-target=service=cloudbuild.googleapis.com \
     --project=ivmanto-com-prod-514215 --format='value(response.keyString)')
   SECRET=$(gcloud secrets versions access latest --secret=webhook-secret --project=ivmanto-com-prod-514215)
   URL="https://cloudbuild.googleapis.com/v1/projects/ivmanto-com-prod-514215/locations/europe-west3/triggers/ivmanto-frontend-rebuild:webhook?key=${KEY}&secret=${SECRET}&trigger=ivmanto-frontend-rebuild&projectId=ivmanto-com-prod-514215"
   printf '%s' "$URL" | gcloud secrets versions add frontend-rebuild-webhook-url --data-file=- --project=ivmanto-com-prod-514215
   ```
   The backend reads `:latest` at deploy, so the next push-to-main build picks it up.

The main trigger overrides `_BACKEND_SA_NAME`, `_GCAL_SA_EMAIL` and `_GCP_PROJECT_ID`, because the yaml's defaults still name the old project.
