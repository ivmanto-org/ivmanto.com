# Manual grants & out-of-band steps — 2026-11-10 (overnight)

PR #125 records, as code, the manual/out-of-band GCP and GitHub operations
performed on **2026-11-10** around IV-INFRA-08 (#105) and the Cloud Build /
GitHub connection. **No `terraform apply` was run as part of this PR** (and
none runs tonight per the overnight rule). Everything below is a *record* so
future applies reconcile against it instead of treating these as drift.

## What is recorded as Terraform

### 1. `ivmanto-github-oauth` Secret Manager secret shell (adopted)

Nick created the secret `ivmanto-github-oauth` in `ivmanto-com-prod-514215`
out-of-band to hold the dasiyes classic **PAT** (scopes `repo`, `read:org`,
`read:user`) used to authorise the Cloud Build 2nd-gen GitHub connection
(commit on the redma thread `ivmanto-relaunch`, 20:33 UTC).

- Added to `secrets.tf` as an **import-only secret shell** (no `*_secret_version`
  resource — the value stays out of Terraform, matching the repo's
  always-on rule). `prevent_destroy` guards it.
- The PAT value itself is **never** in this repo or this PR.

## Manual steps performed tonight (operator / Nick), logged for the record

These are the operations that happened out-of-band on 2026-11-10. The
chronology is from the `ivmanto-relaunch` redma thread (Clauco) and the
merged PR bodies; **any IAM binding listed as a grant is UNVERIFIED** (see
the caveat below) and must be confirmed from the live project before any
future apply touches related resources.

1. **`#105` stage 1 apply** (~19:1x UTC): the Cloud Build SA, its roles and
   the Cloud Build agent's Secret Manager rights (`cloud_build.tf`) were
   applied.
2. **Cloud Build 2nd-gen GitHub connection created + authorised** by Nick:
   `gcloud builds connections create github ivmanto-github --region=europe-west3
   --project=ivmanto-com-prod-514215`, then the Cloud Build GitHub App was
   installed on `ivmanto-org/ivmanto.com` (**app installation 180646944**).
   The connection lives **outside Terraform** (referenced by name via the
   `github_connection_name` variable, default `ivmanto-github`) — it is not
   a managed resource, so nothing here will create/destroy it.
3. **`#105` stage 2 apply** (~20:33 UTC): repository link + both build
   triggers (`google_cloudbuildv2_repository.site`,
   `google_cloudbuild_trigger.main_deploy` / `.frontend_rebuild`) which live
   in `cloud_build.tf`.
4. **Secret `ivmanto-github-oauth`** created with the dasiyes PAT (adopted
   as a secret shell in this PR — see above.)
5. **`#121` merged** → `enable_cloud_build_triggers` default flipped to
   `true` in `variables.tf`, so a plain apply no longer destroys the
   triggers.

## ⚠️ UNVERIFIED — the IAM grant question (needs Nick / audit log)

Clauco's record of tonight is scoped: his session made **no IAM grants**,
ran no applies, and live-checked that the **blog-push SA IAM policy is empty**
(no SA-scoped `serviceAccountUser`/`TokenCreator` binding). Whether any
**project-scoped** IAM grants were applied tonight (e.g. the `#115` Pub/Sub
push-signer `roles/iam.serviceAccountTokenCreator` for
`service-<project-number>@gcp-sa-pubsub.iam.gserviceaccount.com`, or the
Cloud Armor / LB service-user grants) is **NOT confirmed from this seat**,
which has no GCP credentials to read the project IAM policy (`iam policy
list`) or the audit logs.

**Do not treat the steps above as a complete IAM inventory.** Before any
future `terraform apply` against prod, resolve this by:
- `gcloud projects get-iam-policy ivmanto-com-prod-514215 --flatten="bindings[].members" --filter="bindings.role:..."`, and
- `gcloud logging read "protoPayload.methodName=SetIamPolicy AND timestamp>='2026-11-10T00:00:00Z'" --project=ivmanto-com-prod-514215 --freshness=2d`

and record any grants in the appropriate Terraform files. This item is
tracked for Nick on issue #112.
