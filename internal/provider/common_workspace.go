package provider

import (
	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FlattenWorkspace(d *schema.ResourceData, workspace *apiclient.Workspace) error {
	if err := d.Set("id", workspace.WorkspaceId); err != nil {
		return err
	}
	if err := d.Set("customer_id", workspace.CustomerId); err != nil {
		return err
	}
	if workspace.Email != "" {
		if err := d.Set("email", workspace.Email); err != nil {
			return err
		}
	}
	if err := d.Set("name", workspace.Name); err != nil {
		return err
	}
	if err := d.Set("slug", workspace.Slug); err != nil {
		return err
	}
	if workspace.InitialSetupComplete != nil {
		if err := d.Set("initial_setup_complete", workspace.InitialSetupComplete); err != nil {
			return err
		}
	}
	if workspace.DisplaySetupWizard != nil {
		if err := d.Set("display_setup_wizard", workspace.DisplaySetupWizard); err != nil {
			return err
		}
	}
	if workspace.AnonymousDataCollection != nil {
		if err := d.Set("anonymous_data_collection", workspace.AnonymousDataCollection); err != nil {
			return err
		}
	}
	if workspace.News != nil {
		if err := d.Set("news", workspace.News); err != nil {
			return err
		}
	}
	if workspace.SecurityUpdates != nil {
		if err := d.Set("security_updates", workspace.SecurityUpdates); err != nil {
			return err
		}
	}
	if err := d.Set("notification_config", flattenNotifications(&workspace.Notifications)); err != nil {
		return err
	}
	if workspace.FirstCompletedSync != nil {
		if err := d.Set("fist_completed_sync", workspace.FirstCompletedSync); err != nil {
			return err
		}
	}
	if workspace.FeedbackDone != nil {
		if err := d.Set("feedback_done", workspace.FeedbackDone); err != nil {
			return err
		}
	}
	if workspace.DefaultGeography != "" {
		if err := d.Set("default_geography", workspace.DefaultGeography); err != nil {
			return err
		}
	}

	return nil
}

func flattenNotifications(rawNotifs *[]apiclient.Notification) []interface{} {
	if rawNotifs != nil {
		notifs := make([]interface{}, len(*rawNotifs), len(*rawNotifs))

		for i, rawNotif := range *rawNotifs {
			n := make(map[string]interface{})

			n["notification_type"] = rawNotif.NotificationType
			n["send_on_success"] = rawNotif.SendOnSuccess
			n["send_on_failure"] = rawNotif.SendOnFailure
			n["slack_webhook"] = rawNotif.SlackConfiguration.Webhook

			notifs[i] = n
		}

		return notifs
	}

	return make([]interface{}, 0)
}
