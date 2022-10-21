data "airbyte_workspace" "by_id" {
  id = "99cfa08b-daff-4516-b494-a86f6ab0c120"
}

data "airbyte_workspace" "by_slug" {
  slug = "99cfa08b-daff-4516-b494-a86f6ab0c120"
}

output "workspace_by_id" {
  value = data.airbyte_workspace.by_id
}

output "workspace_by_slug" {
  value = data.airbyte_workspace.by_slug
}

data "airbyte_sourcedefinition" "zendesk" {
  id = "c8630570-086d-4a40-99ae-ea5b18673071"
}

output "sourcedefinition_zendesk" {
  value = data.airbyte_sourcedefinition.zendesk
}