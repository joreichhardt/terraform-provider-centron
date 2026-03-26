package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = (*centronProvider)(nil)

type centronProvider struct{}

func New() provider.Provider {
	return &centronProvider{}
}

// Provider provider configuration data model.
type ProviderModel struct {
	BaseURL      types.String `tfsdk:"base_url"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (p *centronProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "centron"
}

func (p *centronProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Centron Cloud API.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "API Base URL. Defaults to 'https://ccenter.centron.de/api/v1'.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "Centron API Client ID. Can also be provided via CENTRON_CLIENT_ID.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Centron API Client Secret. Can also be provided via CENTRON_CLIENT_SECRET.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *centronProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	baseURL := data.BaseURL.ValueString()
	clientID := data.ClientID.ValueString()
	clientSecret := data.ClientSecret.ValueString()

	if clientID == "" {
		clientID = os.Getenv("CENTRON_CLIENT_ID")
	}

	if clientSecret == "" {
		clientSecret = os.Getenv("CENTRON_CLIENT_SECRET")
	}

	if clientID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Client ID",
			"The provider cannot create the Centron API client as there is a missing or empty value for the Centron API client_id. "+
				"Set the client_id value in the configuration or use the CENTRON_CLIENT_ID environment variable.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Client Secret",
			"The provider cannot create the Centron API client as there is a missing or empty value for the Centron API client_secret. "+
				"Set the client_secret value in the configuration or use the CENTRON_CLIENT_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := NewClient(baseURL, clientID, clientSecret)
	
	// Pass the client to resources and data sources
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *centronProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServerResource,
	}
}

func (p *centronProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
//some work}
