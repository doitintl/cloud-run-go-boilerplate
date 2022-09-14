terraform {
  backend "gcs" {
    bucket = "project-id-tf"
    prefix = "state"
  }
}
# for now the provider versions are handled through the lockfile
provider "google" {
  project = "project-id"
  region  = "europe-west1"
  zone    = "europe-west1-b"
}

