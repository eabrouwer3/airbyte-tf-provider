#resource "airbyte_workspace" "simple" {
#  name = "simple_test"
#}
#
#output "simple_airbyte_workspace" {
#  value = airbyte_workspace.simple
#}
#
#resource "airbyte_workspace" "complex" {
#  name = "complex_test"
#  email = "test@example.com"
#  display_setup_wizard = true
#  anonymous_data_collection = false
#  news = true
#  security_updates = true
#  notification_config {
#    notification_type = "slack"
#    send_on_success = true
#    slack_webhook = "http://example.com/webhook"
#  }
#  notification_config {
#    notification_type = "slack"
#    send_on_failure = false
#    slack_webhook = "https://example2.com/cooler-webhook"
#  }
#}
#
#output "complex_airbyte_workspace" {
#  value = airbyte_workspace.complex
#}

resource "airbyte_sourcedefinition" "simple" {
  name = "simple"
  docker_repository = "airbyte/source-postgres"
  docker_image_tag = "1.0.18"
  documentation_url = "https://example.com"
}

output "simple_airbyte_sourcedefinition" {
  value = airbyte_sourcedefinition.simple
}

resource "airbyte_sourcedefinition" "complex" {
  name = "simple"
  docker_repository = "airbyte/source-postgres"
  docker_image_tag = "1.0.17"
  documentation_url = "https://example.com"
  icon = "https://mir-s3-cdn-cf.behance.net/project_modules/max_1200/7a3ec529632909.55fc107b84b8c.png"
}

output "complex_airbyte_sourcedefinition" {
  value = airbyte_sourcedefinition.complex
}