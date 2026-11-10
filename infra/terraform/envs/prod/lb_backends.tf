# IV-INFRA-06: load-balancer backends for Nestbalm's shared global external
# Application LB (cross-project backend service referencing). The NEGs,
# backend services and Cloud Armor policy live here, next to Cloud Run; the
# URL map, certificates and forwarding rule stay in nestbalm-app-prod (#104).

locals {
  # Request logs record the full URL. The Pub/Sub push and the status
  # endpoint no longer carry their token in the URL (OIDC Bearer header
  # and header-shared-secret respectively, #115), so the API backend can
  # be logged without leaking auth material.
  lb_backends = {
    frontend = { log_requests = true }
    backend  = { log_requests = true }
  }
}

resource "google_compute_region_network_endpoint_group" "this" {
  for_each = local.lb_backends

  name                  = "${local.cloud_run_services[each.key].name}-neg"
  region                = var.region
  network_endpoint_type = "SERVERLESS"

  cloud_run {
    service = google_cloud_run_v2_service.this[each.key].name
  }
}

resource "google_compute_security_policy" "edge" {
  name        = "ivmanto-edge-armor"
  description = "Per-IP rate limit for ivmanto.com behind Nestbalm's shared LB."

  rule {
    action      = "throttle"
    priority    = 1000
    description = "Per-IP edge rate limit"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    rate_limit_options {
      conform_action = "allow"
      exceed_action  = "deny(429)"
      enforce_on_key = "IP"
      rate_limit_threshold {
        count        = var.rate_limit_requests_per_minute
        interval_sec = 60
      }
    }
  }

  rule {
    action      = "allow"
    priority    = 2147483647
    description = "Default allow"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
  }
}

resource "google_compute_backend_service" "this" {
  for_each = local.lb_backends

  name                  = "${local.cloud_run_services[each.key].name}-backend"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  protocol              = "HTTPS"
  security_policy       = google_compute_security_policy.edge.id

  backend {
    group = google_compute_region_network_endpoint_group.this[each.key].id
  }

  log_config {
    enable      = each.value.log_requests
    sample_rate = each.value.log_requests ? 1.0 : 0
  }
}

# Lets Nestbalm's URL map (another project) reference these backend services.
resource "google_project_iam_member" "lb_service_user" {
  for_each = toset(var.lb_service_users)

  project = var.project_id
  role    = "roles/compute.loadBalancerServiceUser"
  member  = each.value
}
