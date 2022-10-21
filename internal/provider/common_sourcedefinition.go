package provider

import (
	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FlattenSourceDefinition(d *schema.ResourceData, sd *apiclient.SourceDefinition) error {
	if err := d.Set("id", sd.SourceDefinitionId); err != nil {
		return err
	}
	if err := d.Set("name", sd.Name); err != nil {
		return err
	}
	if err := d.Set("docker_repository", sd.DockerRepository); err != nil {
		return err
	}
	if err := d.Set("docker_image_tag", sd.DockerImageTag); err != nil {
		return err
	}
	if err := d.Set("documentation_url", sd.DocumentationUrl); err != nil {
		return err
	}
	if sd.Icon != "" {
		if err := d.Set("icon", sd.Icon); err != nil {
			return err
		}
	}
	if sd.ProtocolVersion != "" {
		if err := d.Set("protocol_version", sd.ProtocolVersion); err != nil {
			return err
		}
	}
	if sd.ReleaseStage != "" {
		if err := d.Set("release_stage", sd.ReleaseStage); err != nil {
			return err
		}
	}
	if sd.ReleaseDate != "" {
		if err := d.Set("release_date", sd.ReleaseDate); err != nil {
			return err
		}
	}
	if sd.SourceType != "" {
		if err := d.Set("source_type", sd.SourceType); err != nil {
			return err
		}
	}
	if err := d.Set("default_resource_requirements", flattenDefaultReqs(sd.ResourceRequirements)); err != nil {
		return err
	}
	if err := d.Set("job_specific_resource_requirements", flattenJobSpecReqs(sd.ResourceRequirements)); err != nil {
		return err
	}

	return nil
}

func flattenDefaultReqs(reqs *apiclient.ResourceRequirements) []interface{} {
	if reqs != nil {
		rawDefaultReqs := reqs.Default
		if rawDefaultReqs != nil {
			defaultReqs := make([]interface{}, 1, 1)
			req := make(map[string]interface{})

			if rawDefaultReqs.CPURequest != "" {
				req["cpu_request"] = rawDefaultReqs.CPURequest
			}
			if rawDefaultReqs.CPULimit != "" {
				req["cpu_limit"] = rawDefaultReqs.CPULimit
			}
			if rawDefaultReqs.MemoryRequest != "" {
				req["memory_request"] = rawDefaultReqs.MemoryRequest
			}
			if rawDefaultReqs.MemoryLimit != "" {
				req["memory_limit"] = rawDefaultReqs.MemoryLimit
			}

			defaultReqs[0] = req
			return defaultReqs
		}
	}
	return make([]interface{}, 0)
}

func flattenJobSpecReqs(reqs *apiclient.ResourceRequirements) []interface{} {
	if reqs != nil {
		rawJobSpecReqs := reqs.JobSpecific
		if rawJobSpecReqs != nil {
			reqs := make([]interface{}, len(*rawJobSpecReqs), len(*rawJobSpecReqs))

			for i, rawJobSpecReq := range *rawJobSpecReqs {
				req := make(map[string]interface{})

				req["job_type"] = rawJobSpecReq.JobType
				if rawJobSpecReq.ResourceRequirements.CPURequest != "" {
					req["cpu_request"] = rawJobSpecReq.ResourceRequirements.CPURequest
				}
				if rawJobSpecReq.ResourceRequirements.CPULimit != "" {
					req["cpu_limit"] = rawJobSpecReq.ResourceRequirements.CPULimit
				}
				if rawJobSpecReq.ResourceRequirements.MemoryRequest != "" {
					req["memory_request"] = rawJobSpecReq.ResourceRequirements.MemoryRequest
				}
				if rawJobSpecReq.ResourceRequirements.MemoryLimit != "" {
					req["memory_limit"] = rawJobSpecReq.ResourceRequirements.MemoryLimit
				}

				reqs[i] = req
			}

			return reqs
		}
	}
	return make([]interface{}, 0)
}
