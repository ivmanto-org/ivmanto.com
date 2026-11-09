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

Secret **values** are never stored in Terraform; it only creates the Secret Manager
secrets. Add versions with `gcloud secrets versions add <name> --data-file=-`.

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
