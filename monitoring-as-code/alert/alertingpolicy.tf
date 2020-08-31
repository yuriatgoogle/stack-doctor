provider "google" {
  project = "stack-doctor"
  region  = "us-west3"
  zone    = "us-west3-a"
}
resource "google_monitoring_notification_channel" "email0" {
  display_name = "GAE Service Oncall"
  type = "email"
  labels = {
    email_address = "website-oncall@example.com"
  }
}

locals {
  email0_id = "${google_monitoring_notification_channel.email0.name}"
}
resource "google_monitoring_alert_policy" "alert_policy" {
  display_name = "GAE Default Service Down"
  combiner     = "OR"
  conditions {
    display_name = "Uptime Check Failure"
    condition_threshold {
      filter     = "metric.type=\"monitoring.googleapis.com/uptime_check/check_passed\" resource.type=\"uptime_url\" metric.label.\"check_id\"=\"gae-uptime-check\""
      duration   = "300s"
      comparison = "COMPARISON_GT"
      aggregations {
        per_series_aligner = "ALIGN_NEXT_OLDER"
        cross_series_reducer =  "REDUCE_COUNT_FALSE"
        alignment_period = "60s"
      }
    }
  }
  notification_channels = [
    "${local.email0_id}"
  ]
}