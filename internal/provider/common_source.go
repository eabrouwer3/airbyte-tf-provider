package provider

import (
	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FlattenSource(d *schema.ResourceData, s *apiclient.Source) error {
	if err := d.Set("id", s.SourceId); err != nil {
		return err
	}
	if err := d.Set("name", s.Name); err != nil {
		return err
	}
	if err := d.Set("sourcedefinition_id", s.SourceDefinitionId); err != nil {
		return err
	}
	if err := d.Set("workspace_id", s.WorkspaceId); err != nil {
		return err
	}
	if err := d.Set("source_name", s.SourceName); err != nil {
		return err
	}
	if err := d.Set("icon", s.Icon); err != nil {
		return err
	}
	if err := d.Set("connection_configuration", s.ConnectionConfiguration); err != nil {
		return err
	}

	return nil
}
