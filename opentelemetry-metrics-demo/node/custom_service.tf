provider "google" {
  project = "stack-doctor"
  region  = "us-west3"
  zone    = "us-west3-a"
}

resource "google_monitoring_custom_service" "terraform-service" {
  service_id = "terraform-service"
  display_name = "Service Created by Terraform"
}

resource "google_monitoring_slo" "request_based_slo" {
  service = google_monitoring_custom_service.terraform-service.service_id
  slo_id = "terraform-slo"
  display_name = "99% of 10-min windows in rolling day have mean latency under 8s"

  goal = 0.99
  rolling_period_days = 1

  windows_based_sli {
    window_period = "600s"
    metric_mean_in_range {
      time_series = join(" AND ", [
        "metric.type=\"external.googleapis.com/prometheus/response_latency\"",
        "resource.type=\"k8s_container\"",
        "resource.label.\"container_name\"=\"opentelemetry-server\""
      ])

      range {
        max = 8000
      }
    }
  }
}