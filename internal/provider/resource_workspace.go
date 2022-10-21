package provider

import (
	"context"
	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkspace() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Airbyte Workspace",

		CreateContext: resourceWorkspaceCreate,
		ReadContext:   resourceWorkspaceRead,
		UpdateContext: resourceWorkspaceUpdate,
		DeleteContext: resourceWorkspaceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Workspace ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"customer_id": {
				Description: "Customer ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "Customer Email",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "Workspace Name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"slug": {
				Description: "Workspace Slug",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"initial_setup_complete": {
				Description: "Is the initial setup complete",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"display_setup_wizard": {
				Description: "Should the UI display the setup wizard",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"anonymous_data_collection": {
				Description: "Is anonymous data collection turned on",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"news": {
				Description: "Should the UI show news updates",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"security_updates": {
				Description: "Should the UI show security updates",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"notification_config": {
				Description: "Notification systems set up",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notification_type": {
							Description:  "Possible value: slack",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"slack"}, false),
						},
						"send_on_success": {
							Description: "Should the notification be sent for successes",
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
						},
						"send_on_failure": {
							Description: "Should the notification be sent for failures",
							Type:        schema.TypeBool,
							Default:     true,
							Optional:    true,
						},
						"slack_webhook": {
							Description: "Configuration for Slack notifications - See https://slack.com/help/articles/115005265063-Incoming-webhooks-for-Slack",
							Type:        schema.TypeString,
							Required:    true,
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

func resourceToWorkspace(d *schema.ResourceData) apiclient.CommonWorkspaceFields {
	workspace := apiclient.CommonWorkspaceFields{}

	if v, ok := d.GetOk("email"); ok {
		workspace.Email = v.(string)
	}
	if v, ok := d.GetOk("anonymous_data_collection"); ok {
		val := v.(bool)
		workspace.AnonymousDataCollection = &val
	}
	if v, ok := d.GetOk("news"); ok {
		val := v.(bool)
		workspace.News = &val
	}
	if v, ok := d.GetOk("security_updates"); ok {
		val := v.(bool)
		workspace.SecurityUpdates = &val
	}
	if v, ok := d.GetOk("display_setup_wizard"); ok {
		val := v.(bool)
		workspace.DisplaySetupWizard = &val
	}

	notifInput, notifOk := d.GetOk("notification_config")
	if notifOk {
		var notifs []apiclient.Notification

		for _, rawNotif := range notifInput.([]interface{}) {
			rn := rawNotif.(map[string]interface{})

			n := apiclient.Notification{
				NotificationType: rn["notification_type"].(string),
				SendOnSuccess:    rn["send_on_success"].(bool),
				SendOnFailure:    rn["send_on_failure"].(bool),
				SlackConfiguration: apiclient.SlackConfiguration{
					Webhook: rn["slack_webhook"].(string),
				},
			}

			notifs = append(notifs, n)
		}
		workspace.Notifications = notifs
	}

	return workspace
}

func resourceWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	newWorkspace := apiclient.NewWorkspace{
		WorkspaceNameBody: apiclient.WorkspaceNameBody{
			Name: d.Get("name").(string),
		},
		CommonWorkspaceFields: resourceToWorkspace(d),
	}

	w, err := client.CreateWorkspace(newWorkspace)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(w.WorkspaceId)

	resourceWorkspaceRead(ctx, d, meta)

	return diags
}

func resourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*apiclient.ApiClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	workspaceId := d.Id()

	w, err := c.GetWorkspaceById(workspaceId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = FlattenWorkspace(d, w)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	updatedWorkspace := apiclient.UpdatedWorkspace{
		WorkspaceIdBody: apiclient.WorkspaceIdBody{
			WorkspaceId: d.Get("id").(string),
		},
		CommonWorkspaceFields: resourceToWorkspace(d),
	}

	w, err := client.UpdateWorkspace(updatedWorkspace)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(w.WorkspaceId)

	resourceWorkspaceRead(ctx, d, meta)

	return diags
}

func resourceWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	workspaceId := d.Id()

	err := client.DeleteWorkspace(workspaceId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
