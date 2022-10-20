resource "airbyte_workspace" "simple" {
  name = "simple_test"
}

output "simple_airbyte_workspace" {
  value = airbyte_workspace.simple
}

resource "airbyte_workspace" "complex" {
  name = "complex_test"
  email = "test@example.com"
  display_setup_wizard = true
  anonymous_data_collection = false
  news = true
  security_updates = true
  notification_config {
    notification_type = "slack"
    send_on_success = true
    slack_webhook = "http://example.com/webhook"
  }
  notification_config {
    notification_type = "slack"
    send_on_failure = false
    slack_webhook = "https://example2.com/cooler-webhook"
  }
}

output "complex_airbyte_workspace" {
  value = airbyte_workspace.complex
}