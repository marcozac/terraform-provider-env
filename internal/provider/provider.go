package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const ProviderName = "env"

func New() func() provider.Provider {
	return func() provider.Provider {
		return &envProvider{}
	}
}

var (
	_ provider.Provider              = (*envProvider)(nil)
	_ provider.ProviderWithFunctions = (*envProvider)(nil)
)

type envProvider struct{}

func (p *envProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = ProviderName
}

func (p *envProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A provider for interacting with environment variables.",
	}
}

func (p *envProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *envProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *envProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFileDataSource,
	}
}

func (p *envProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewGetenvFunction,
	}
}
