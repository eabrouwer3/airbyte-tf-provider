package provider

import (
	"context"
	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO: Implement this eventually - not super important right no

func dataSourceSourceDefinitionSpec() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Airbyte Source Definition by id",
		ReadContext: dataSourceSourceDefinitionRead,
		Schema:      map[string]*schema.Schema{},
	}
}

func dataSourceSourceDefinitionSpecRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	sdId := d.Get("id").(string)

	sd, err := client.GetSourceDefinitionSpec(sdId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = FlattenSourceDefinition(d, sd)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sdId)

	return diags
}
