# GCP infrastructure as Terraform (ivmanto-com-prod)

## Goal
Recreate the ivmanto.com runtime on Google Cloud with Terraform, in the dedicated project
`ivmanto-com-prod-514215` (nestbalm.com org, folder `Prod`, billing `ivmanto-billing-accnt`),
served through the existing Nestbalm global external Application Load Balancer.

## Decisions (agreed 2026-11-09)
- Project: `ivmanto-com-prod-514215`, created manually, Terraform takes it over from here.
- Load balancer: **shared** with Nestbalm (`ff-app-*` in `nestbalm-app-prod`, IP `136.68.185.92`),
  using cross-project backend service referencing (same org, no Shared VPC needed).
- Terraform state: GCS bucket in `ivmanto-com-prod-514215`.
- Region: `europe-west3` (as before).

## Split of responsibility

| Lives in | Resources |
|---|---|
| **this repo** `infra/terraform/` → `ivmanto-com-prod-514215` | APIs, Artifact Registry, service accounts, Secret Manager secrets (shells only), Cloud Run frontend + backend, serverless NEGs, backend services, Cloud Armor policy, blog GCS bucket, Pub/Sub (blog refresh), IAM grant `roles/compute.loadBalancerServiceUser` to Nestbalm's Terraform identity, CI/CD identity (WIF) |
| **Nestbalm repo** `infra/terraform/modules/app_edge` → `nestbalm-app-prod` | URL map host rules `ivmanto.com`, `www.ivmanto.com` → ivmanto backends (`/api/*` → backend, rest → frontend), managed cert for ivmanto domains added to the HTTPS proxy |
| **DNS** (ivmanto.com registrar/DNS) | `A ivmanto.com`, `A www.ivmanto.com` → `136.68.185.92` |

## Plan
- [x] **IV-INFRA-1** Bootstrap: state bucket (versioned, uniform access), enable APIs, `infra/terraform/envs/prod` skeleton with `backend "gcs"`, pinned provider versions
- [ ] **IV-INFRA-2** Artifact Registry repo, runtime service accounts (frontend, backend) with least-privilege IAM
- [ ] **IV-INFRA-3** Secret Manager secret shells: `smtp-password`, `ideas-prompt`, `ga-api-secret`, `pubsub-push-token`, `frontend-rebuild-webhook-url` (values added by Nick out-of-band, never in Terraform)
- [ ] **IV-INFRA-4** Blog GCS bucket + Pub/Sub notification → backend push subscription
- [ ] **IV-INFRA-5** Cloud Run services `ivmanto-com` (frontend) and `ivmanto-backend-service`, ingress `internal-and-cloud-load-balancing`, placeholder image until CI pushes the first build
- [ ] **IV-INFRA-6** Serverless NEGs + backend services + Cloud Armor policy; grant LB service user to Nestbalm's identity; output backend service self-links
- [ ] **IV-INFRA-7** Nestbalm repo PR: host rules + ivmanto cert in `app_edge` (separate PR in `Nestbalm/nestbalm-saas`, reviewed there)
- [ ] **IV-INFRA-8** CI/CD stays on **Cloud Build, `cloudbuild*.yaml` unchanged**. Terraform recreates: Cloud Build GitHub connection (2nd gen) to `ivmanto-org/ivmanto.com` (Nick authorises the GitHub App once), push-to-`main` trigger → `cloudbuild.yaml`, webhook trigger → `cloudbuild-frontend-only.yaml` (its URL stored in secret `frontend-rebuild-webhook-url`, used by the backend after a blog publish), Cloud Build SA with deploy rights
- [ ] **IV-INFRA-9** DNS cut-over, wait for cert ACTIVE, smoke test (`/`, `/api/health`, booking, blog)
- [ ] **IV-INFRA-10** Calendar DWD: new backend SA's client ID authorised in ivmanto.com Workspace admin (Nick)

## Needs Nick
- Secret values (SMTP, GA API secret, ideas prompt) for IV-INFRA-3
- Workspace admin DWD authorisation for IV-INFRA-10
- DNS change for IV-INFRA-9
- Review/merge of the Nestbalm LB PR (IV-INFRA-7)

## Risks
- Shared LB: a bad URL-map change affects nestbalm.com too. Mitigation: Nestbalm-side change is additive host rules only; `terraform plan` reviewed before apply.
- Managed cert only goes ACTIVE after DNS points at the LB, so expect a short window without HTTPS on ivmanto.com at cut-over (site is down today anyway).

## Review
_(to be filled in after implementation)_
