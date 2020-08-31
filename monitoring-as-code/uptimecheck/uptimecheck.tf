provider "google" {
  project = "stack-doctor"
  region  = "us-west3"
  zone    = "us-west3-a"
}

resource "google_monitoring_uptime_check_config" "https" {
  display_name = "Test Uptime Check"
  timeout = "60s"

  http_check {
    path = "/"
    port = "443"
    use_ssl = true
    validate_ssl = true
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = "stack-doctor"
      host = "www.google.com"
    }
  }

  content_matchers {
    content = "lucky"
  }
}