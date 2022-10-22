package provider

import (
	"context"
	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

func resourceSourceDefinition() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Airbyte Source Definition",

		CreateContext: resourceSourceDefinitionCreate,
		ReadContext:   resourceSourceDefinitionRead,
		UpdateContext: resourceSourceDefinitionUpdate,
		DeleteContext: resourceSourceDefinitionDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Source Definition ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Source Definition Name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"docker_repository": {
				Description: "Docker Repository URL (e.g. 112233445566.dkr.ecr.us-east-1.amazonaws.com/source-custom) or DockerHub identifier (e.g. airbyte/source-postgres)",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"docker_image_tag": {
				Description: "Docker image tag",
				Type:        schema.TypeString,
				Required:    true,
			},
			"documentation_url": {
				Description: "Documentation URL",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"icon": {
				Description: "URL for the icon displayed in the UI",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
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
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_request": {
							Description: "CPU Requested",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"cpu_limit": {
							Description: "CPU Limit",
							Type:        schema.TypeString,
							Optional:    false,
							Computed:    true,
						},
						"memory_request": {
							Description: "Memory Requested",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"memory_limit": {
							Description: "Memory Limit",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"job_specific_resource_requirements": {
				Description: "Sets resource requirements for a specific job type for an actor definition. These values override " +
					"the default, if both are set. These values will be overridden by configuration at the connection level.",
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_type": {
							Description:  "Allowed: get_spec | check_connection | discover_schema | sync | reset_connection | connection_updater | replicate",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"get_spec", "check_connection", "discover_schema", "sync", "reset_connection", "connection_updater", "replicate"}, false),
						},
						"cpu_request": {
							Description: "CPU Requested",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"cpu_limit": {
							Description: "CPU Limit",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"memory_request": {
							Description: "Memory Requested",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"memory_limit": {
							Description: "Memory Limit",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func setSourceDefinitionFields(d *schema.ResourceData) apiclient.CommonSourceDefinitionFields {
	sd := apiclient.CommonSourceDefinitionFields{}

	if v, ok := d.GetOk("name"); ok {
		sd.Name = v.(string)
	}
	if v, ok := d.GetOk("docker_repository"); ok {
		sd.DockerRepository = v.(string)
	}
	if v, ok := d.GetOk("docker_image_tag"); ok {
		sd.DockerImageTag = v.(string)
	}
	if v, ok := d.GetOk("documentation_url"); ok {
		sd.DocumentationUrl = v.(string)
	}
	if v, ok := d.GetOk("icon"); ok {
		sd.Icon = v.(string)
	}

	_, defaultReqOk := d.GetOk("default_resource_requirements")
	_, jobSpecReqOk := d.GetOk("job_specific_resource_requirements")
	if defaultReqOk || jobSpecReqOk {
		sd.ResourceRequirements = setReqFields(d)
	}

	return sd
}

func setReqFields(d *schema.ResourceData) *apiclient.ResourceRequirements {
	reqs := apiclient.ResourceRequirements{}

	defaultReqs, defaultReqOk := d.GetOk("default_resource_requirements")
	if defaultReqOk && len(defaultReqs.([]interface{})) > 0 {
		defaultReq := defaultReqs.([]map[string]interface{})[0]
		if v := defaultReq["cpu_request"].(string); v != "" {
			reqs.Default.CPURequest = v
		}
		if v := defaultReq["cpu_limit"].(string); v != "" {
			reqs.Default.CPULimit = v
		}
		if v := defaultReq["memory_request"].(string); v != "" {
			reqs.Default.MemoryRequest = v
		}
		if v := defaultReq["memory_limit"].(string); v != "" {
			reqs.Default.MemoryLimit = v
		}
	}

	jobSpecInput, jobSpecOk := d.GetOk("job_specific_resource_requirements")
	if jobSpecOk {
		var jobSpecReqs []apiclient.JobSpecificResourceRequirements

		for _, rawJobSpec := range jobSpecInput.([]interface{}) {
			rjs := rawJobSpec.(map[string]interface{})

			js := apiclient.JobSpecificResourceRequirements{
				JobType: rjs["job_type"].(string),
			}

			if v := rjs["cpu_request"].(string); v != "" {
				js.ResourceRequirements.CPURequest = v
			}
			if v := rjs["cpu_limit"].(string); v != "" {
				js.ResourceRequirements.CPULimit = v
			}
			if v := rjs["memory_request"].(string); v != "" {
				js.ResourceRequirements.MemoryRequest = v
			}
			if v := rjs["memory_limit"].(string); v != "" {
				js.ResourceRequirements.MemoryLimit = v
			}

			jobSpecReqs = append(jobSpecReqs, js)
		}
		reqs.JobSpecific = &jobSpecReqs
	}

	return &reqs
}

func resourceSourceDefinitionCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	newSD := setSourceDefinitionFields(d)

	sd, err := client.CreateSourceDefinition(newSD)
	if err != nil {
		if strings.Contains(err.Error(), "status: 500") {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create sourceDefinition. Airbyte likely unable to find/access specified docker_repository or docker_image_tag.",
				Detail:   err.Error(),
			})
			return diags
		}
		return diag.FromErr(err)
	}

	d.SetId(sd.SourceDefinitionId)

	resourceSourceDefinitionRead(ctx, d, meta)

	return diags
}

func resourceSourceDefinitionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*apiclient.ApiClient)

	var diags diag.Diagnostics

	sdId := d.Id()

	sd, err := c.GetSourceDefinitionById(sdId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = FlattenSourceDefinition(d, sd)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSourceDefinitionUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	updatedSourceDefinition := apiclient.UpdatedSourceDefinition{
		SourceDefinitionIdBody: apiclient.SourceDefinitionIdBody{
			SourceDefinitionId: d.Get("id").(string),
		},
		DockerImageTag:       d.Get("docker_image_tag").(string),
		ResourceRequirements: nil,
	}
	if v, ok := d.GetOk("docker_image_tag"); ok {
		updatedSourceDefinition.DockerImageTag = v.(string)
	}
	_, defaultReqOk := d.GetOk("default_resource_requirements")
	_, jobSpecReqOk := d.GetOk("job_specific_resource_requirements")
	if defaultReqOk || jobSpecReqOk {
		updatedSourceDefinition.ResourceRequirements = setReqFields(d)
	}

	sd, err := client.UpdateSourceDefinition(updatedSourceDefinition)
	if err != nil {
		if strings.Contains(err.Error(), "status: 500") {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create sourceDefinition. Airbyte likely unable to find/access specified docker_repository or docker_image_tag.",
				Detail:   err.Error(),
			})
			return diags
		}
		return diag.FromErr(err)
	}

	d.SetId(sd.SourceDefinitionId)

	resourceSourceDefinitionRead(ctx, d, meta)

	return diags
}

func resourceSourceDefinitionDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*apiclient.ApiClient)
	var diags diag.Diagnostics

	sdId := d.Id()

	err := client.DeleteSourceDefinition(sdId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
