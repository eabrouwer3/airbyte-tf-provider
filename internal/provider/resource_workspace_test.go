package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkspace_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkspace_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_workspace.basic", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("airbyte_workspace.basic", "slug", regexp.MustCompile("^basic_test")),
					resource.TestCheckResourceAttr("airbyte_workspace.basic", "name", "basic_test"),
					resource.TestCheckResourceAttr("airbyte_workspace.basic", "notification_config.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceWorkspace_complex(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkspace_complex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_workspace.complex", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("airbyte_workspace.complex", "slug", regexp.MustCompile("^complex_test")),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "name", "complex_test"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "email", "test@example.com"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "display_setup_wizard", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "anonymous_data_collection", "false"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "news", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "security_updates", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.#", "2"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.0.notification_type", "slack"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.0.send_on_success", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.0.send_on_failure", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.0.slack_webhook", "http://example.com/webhook"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.1.notification_type", "slack"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.1.send_on_success", "false"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.1.send_on_failure", "false"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.1.slack_webhook", "https://example2.com/cooler-webhook"),
				),
			},
			{
				Config: testAccResourceWorkspace_complexChange,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.#", "1"),
				),
			},
		},
	})
}

const testAccResourceWorkspace_basic = `
resource "airbyte_workspace" "basic" {
  name = "basic_test"
}
`

const testAccResourceWorkspace_complex = `
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
`

const testAccResourceWorkspace_complexChange = `
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
}
`
