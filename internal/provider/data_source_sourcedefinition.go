package provider

import (
	"context"
	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSourceDefinition() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Airbyte Source Definition by id",
		ReadContext: dataSourceSourceDefinitionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Source Definition ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Source Definition Name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"docker_repository": {
				Description: "Docker Repository URL (e.g. 112233445566.dkr.ecr.us-east-1.amazonaws.com/source-custom) or DockerHub identifier (e.g. airbyte/source-postgres)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"docker_image_tag": {
				Description: "Docker image tag",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"documentation_url": {
				Description: "Documentation URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"icon": {
				Description: "URL for the icon displayed in the UI",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"protocol_version": {
				Description: "The Airbyte Protocol version supported by the connector",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"release_stage": {
				Description: "Allowed: alpha | beta | generally_available | custom",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"release_date": {
				Description: "The date when this connector was first released, in yyyy-mm-dd format",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"source_type": {
				Description: "Allowed: api | file | database | custom",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default_resource_requirements": {
				Description: "Actor definition specific resource requirements. Ff default is set, these are the requirements " +
					"that should be set for ALL jobs run for this actor definition. It is overridden by the job type specific " +
					"configurations. If not set, the platform will use defaults. These values will be overridden by configuration " +
					"at the connection level.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_request": {
							Description: "CPU Requested",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cpu_limit": {
							Description: "CPU Limit",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"memory_request": {
							Description: "Memory Requested",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"memory_limit": {
							Description: "Memory Limit",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"job_specific_resource_requirements": {
				Description: "Sets resource requirements for a specific job type for an actor definition. These values override " +
					"the default, if both are set. These values will be overridden by configuration at the connection level.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_type": {
							Description: "Allowed: get_spec | check_connection | discover_schema | sync | reset_connection | connection_updater | replicate",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cpu_request": {
							Description: "CPU Requested",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cpu_limit": {
							Description: "CPU Limit",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"memory_request": {
							Description: "Memory Requested",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"memory_limit": {
							Description: "Memory Limit",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSourceDefinitionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	sdId := d.Get("id").(string)

	sd, err := client.GetSourceDefinitionById(sdId)
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
