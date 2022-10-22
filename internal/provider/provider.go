package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/eabrouwer3/terraform-provider-airbyte/internal/apiclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"host_url": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("AIRBYTE_URL", "http://localhost:8000"),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"airbyte_workspace":        dataSourceWorkspace(),
				"airbyte_sourcedefinition": dataSourceSourceDefinition(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"airbyte_workspace":        resourceWorkspace(),
				"airbyte_sourcedefinition": resourceSourceDefinition(),
				"airbyte_source":           resourceSource(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		host := d.Get("host_url").(string)

		return &apiclient.ApiClient{
			HTTPClient: &http.Client{Timeout: 120 * time.Second},
			HostURL:    host,
		}, nil
	}
}
