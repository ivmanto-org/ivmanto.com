#!/usr/bin/env bash
# Builds the static site from the current checkout and deploys it to the
# password-protected staging service (infra/terraform/envs/staging).
# Usage: scripts/deploy-staging.sh            build, push and deploy
#        scripts/deploy-staging.sh --no-deploy  build and push only
set -euo pipefail

PROJECT_ID="${STAGING_PROJECT_ID:-ivmanto-com-sta}"
REGION="${STAGING_REGION:-europe-west3}"
SERVICE="ivmanto-staging"
REPOSITORY="${REGION}-docker.pkg.dev/${PROJECT_ID}/staging"

repo_root="$(git rev-parse --show-toplevel)"
cd "${repo_root}"

if [[ -n "$(git status --porcelain --untracked-files=no)" ]]; then
  echo "Uncommitted changes: commit or stash them so staging matches a commit." >&2
  exit 1
fi

commit="$(git rev-parse --short HEAD)"
image="${REPOSITORY}/site:${commit}"
context="$(mktemp -d)"
trap 'rm -rf "${context}"' EXIT

echo "Generating the static site at ${commit}"
node scripts/check-dashes.mjs
npx nuxi generate

cp staging/server/Dockerfile staging/server/go.mod staging/server/main.go "${context}/"
cp -r .output/public "${context}/site"

echo "Building and pushing ${image}"
docker build --quiet --tag "${image}" "${context}"
docker push --quiet "${image}"

if [[ "${1:-}" == "--no-deploy" ]]; then
  echo "Pushed ${image} (not deployed)"
  exit 0
fi

# The service account is pinned so a recreated service never falls back to the
# compute default (which cannot read the credentials secret).
gcloud run deploy "${SERVICE}" --project="${PROJECT_ID}" --region="${REGION}" --image="${image}" \
  --service-account="staging-site-sa@${PROJECT_ID}.iam.gserviceaccount.com" --quiet
gcloud run services describe "${SERVICE}" --project="${PROJECT_ID}" --region="${REGION}" --format='value(status.url)'
