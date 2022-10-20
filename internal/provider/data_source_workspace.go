package provider

import (
	"context"
	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Airbyte Workspace by id or slug",
		ReadContext: dataSourceWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description:  "Workspace ID",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"id", "slug"},
			},
			"customer_id": {
				Description: "Customer ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "Customer Email",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Workspace Name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"slug": {
				Description:  "Workspace Slug",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"id", "slug"},
			},
			"initial_setup_complete": {
				Description: "Is the initial setup complete",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"display_setup_wizard": {
				Description: "Should the UI display the setup wizard",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"anonymous_data_collection": {
				Description: "Is anonymous data collection turned on",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"news": {
				Description: "Should the UI show news updates",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"security_updates": {
				Description: "Should the UI show security updates",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"notification_config": {
				Description: "Notification systems set up",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notification_type": {
							Description: "Possible value: slack",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"send_on_success": {
							Description: "Should the notification be sent for successes",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"send_on_failure": {
							Description: "Should the notification be sent for failures",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"slack_webhook": {
							Description: "Configuration for Slack notifications - See https://slack.com/help/articles/115005265063-Incoming-webhooks-for-Slack",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"fist_completed_sync": {
				Description: "Has a first sync completed",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"feedback_done": {
				Description: "",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"default_geography": {
				Description: "Possible values: auto | us | eu",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	workspaceId := d.Get("id").(string)
	slug := d.Get("slug").(string)

	if workspaceId != "" && slug != "" {
		return diag.Errorf("Only one of `id` and `slug` can be set")
	}

	var workspace *apiclient.Workspace
	var err error
	if workspaceId != "" {
		workspace, err = client.GetWorkspaceById(workspaceId)
	} else if slug != "" {
		workspace, err = client.GetWorkspaceBySlug(slug)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	// Flatten workspace to schema
	err = FlattenWorkspace(d, workspace)
	if err != nil {
		return diag.FromErr(err)
	}

	if workspaceId != "" {
		d.SetId(workspace.WorkspaceId)
	} else if slug != "" {
		d.SetId(workspace.Slug)
	}

	return diags
}
