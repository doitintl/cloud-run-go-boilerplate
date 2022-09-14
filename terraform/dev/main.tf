# Workload Identity Federation
module "gh_oidc" {
  source      = "terraform-google-modules/github-actions-runners/google//modules/gh-oidc"
  project_id  = "project-id"
  pool_id     = "gha-gcp"
  provider_id = "gha-gcp"
  sa_mapping = {
    "gha_terraform" = {
      sa_name   = "projects/project-id/serviceAccounts/terraform@project-id.iam.gserviceaccount.com"
      attribute = "attribute.repository/GITHUB_ORG/REPOSITORY"
    }
  }
}

# Registry
resource "google_artifact_registry_repository" "go-boilerplate-containers" {
  location      = "europe-west1"
  repository_id = "go-boilerplate-containers"
  description   = "here be dev containers"
  format        = "DOCKER"
}

### Outputs ###
output "service_url" {
  value = google_cloud_run_service.go-boilerplate-server.status[0].url
}
