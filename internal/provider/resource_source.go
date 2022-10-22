package provider

import (
	"context"
	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Airbyte Source Definition",

		CreateContext: resourceSourceCreate,
		ReadContext:   resourceSourceRead,
		UpdateContext: resourceSourceUpdate,
		DeleteContext: resourceSourceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Source ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"sourcedefinition_id": {
				Description: "Source Definition ID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"workspace_id": {
				Description: "Workspace ID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Name of the Source",
				Type:        schema.TypeString,
				Required:    true,
			},
			"source_name": {
				Description: "Name of the Source Definition",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"icon": {
				Description: "URL for the icon displayed in the UI",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"connection_configuration": {
				Description: "Map of Credentials for the source",
				Type:        schema.TypeMap,
				Required:    true,
				Elem:        any,
			},
		},
	}
}

func setSourceFields(d *schema.ResourceData) apiclient.CommonSourceFields {
	s := apiclient.CommonSourceFields{}

	if v, ok := d.GetOk("name"); ok {
		s.Name = v.(string)
	}
	if v, ok := d.GetOk("connection_configuration"); ok {
		s.ConnectionConfiguration = v.(map[string]any)
	}

	return s
}

func resourceSourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	newSource := apiclient.NewSource{
		SourceDefinitionIdBody: apiclient.SourceDefinitionIdBody{
			SourceDefinitionId: d.Get("sourcedefinition_id").(string),
		},
		WorkspaceIdBody: apiclient.WorkspaceIdBody{
			WorkspaceId: d.Get("workspace_id").(string),
		},
		CommonSourceFields: setSourceFields(d),
	}

	s, err := client.CreateSource(newSource)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(s.SourceId)

	resourceSourceRead(ctx, d, meta)

	return diags
}

func resourceSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*apiclient.ApiClient)

	var diags diag.Diagnostics

	sdId := d.Id()

	s, err := c.GetSourceById(sdId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = FlattenSource(d, s)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	updatedSource := apiclient.UpdatedSource{
		SourceIdBody: apiclient.SourceIdBody{
			SourceId: d.Get("id").(string),
		},
		CommonSourceFields: setSourceFields(d),
	}

	s, err := client.UpdateSource(updatedSource)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(s.SourceId)

	resourceSourceRead(ctx, d, meta)

	return diags
}

func resourceSourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	sourceId := d.Id()

	err := client.DeleteSource(sourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
