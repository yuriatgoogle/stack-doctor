provider "google" {
  project = "next-2020-ops302"
  region  = "us-west3"
  zone    = "us-west3-a"
}

resource "google_monitoring_uptime_check_config" "https" {
  display_name = "GAE Uptime Check"
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
      project_id = "next-2020-ops302"
      host = "next-2020-ops302.appspot.com"
    }
  }

  content_matchers {
    content = "Hello World"
  }
}