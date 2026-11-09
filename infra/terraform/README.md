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
gcloud auth application-default login   # an account with Owner on the project

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
