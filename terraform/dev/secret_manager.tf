# Secret Manager
resource "google_secret_manager_secret" "api_token" {
  secret_id = "api_token"
  replication { automatic = true }
}

data "google_compute_default_service_account" "cdsa" {}
resource "google_secret_manager_secret_iam_member" "cdsa_api_t" {
  secret_id = google_secret_manager_secret.api_token.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_compute_default_service_account.cdsa.email}"
}
