resource "airbyte_workspace" "simple" {
  name = "simple_test"
}

output "simple_airbyte_workspace" {
  value = airbyte_workspace.simple
}

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
  docker_repository = "airbyte/source-github"
  docker_image_tag = "0.3.7"
  documentation_url = "https://hub.docker.com/r/airbyte/source-github"
}

output "simple_airbyte_sourcedefinition" {
  value = airbyte_sourcedefinition.simple
}

#resource "airbyte_sourcedefinition" "complex" {
#  name = "simple"
#  docker_repository = "airbyte/source-postgres"
#  docker_image_tag = "1.0.17"
#  documentation_url = "https://example.com"
#  icon = "https://mir-s3-cdn-cf.behance.net/project_modules/max_1200/7a3ec529632909.55fc107b84b8c.png"
#}
#
#output "complex_airbyte_sourcedefinition" {
#  value = airbyte_sourcedefinition.complex
#}

resource "airbyte_source" "simple" {
  sourcedefinition_id = airbyte_sourcedefinition.simple.id
  workspace_id = airbyte_sourcedefinition.simple.id
  name = "simple_source"
  connection_configuration = {
    credentials = {
      personal_access_token = "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"
    }
    start_date = "2022-10-01"
    repository = "eabrouwer3/terraform-provider-airbyte "
  }
}

output "simple_airbyte_source" {
  value = airbyte_source.simple
}