# Cloud Run
resource "google_cloud_run_service" "getcre-server" {
  name     = "getcre-server"
  location = "europe-west1"

  # let it not explode when external changes are made
  autogenerate_revision_name = true

  template {
    spec {
      containers {
        image = "europe-west1-docker.pkg.dev/cloud-run-go-boilerplate/cloud-run-go-boilerplate-containers/cloud-run-go-boilerplate:sha-1517d4e"
        env {
          name = "API_TOKEN"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.api_token.secret_id
              key  = "2"
            }
          }
        }
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
  lifecycle {
    ignore_changes = [template[0].spec[0].containers[0].image]
  }
}

### Publicly Open
data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}
resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = google_cloud_run_service.getcre-server.location
  project     = google_cloud_run_service.getcre-server.project
  service     = google_cloud_run_service.getcre-server.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

