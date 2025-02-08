package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joho/godotenv"
)

func NewFileDataSource() datasource.DataSource {
	return &fileDataSource{}
}

type fileDataSource struct {
	data *fileDataSourceModel
}

type fileDataSourceModel struct {
	Path   types.String `tfsdk:"path"`
	Result types.Map    `tfsdk:"result"`
}

func (d *fileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *fileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Reads the environment variables from a file.",
		MarkdownDescription: "Use this data source to read environment variables from a file.",

		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Path to the file with the environment variables." +
					"Defaults to `.env`.",
				MarkdownDescription: "Path to the file with the environment variables." +
					"Defaults to `.env`.",
			},
			"result": schema.MapAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Sensitive:           true,
				Description:         "The environment variables read from the file.",
				MarkdownDescription: "The environment variables read from the file.",
			},
		},
	}
}

func (d *fileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	d.data = &fileDataSourceModel{}
	resp.Diagnostics.Append(req.Config.Get(ctx, d.data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.data.Path.ValueString() == "" {
		tflog.Warn(ctx, "path not set, using default")
		d.data.Path = types.StringValue(".env")
	}

	f, err := os.Open(d.data.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to open file", err.Error())
		return
	}
	defer f.Close()

	env, err := godotenv.Parse(f)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse file", err.Error())
		return
	}

	m := make(map[string]attr.Value, len(env))
	for k, v := range env {
		m[k] = types.StringValue(v)
	}
	d.setResult(ctx, resp, m)
}

func (d *fileDataSource) setResult(ctx context.Context, resp *datasource.ReadResponse, elements map[string]attr.Value) {
	var ds diag.Diagnostics
	d.data.Result, ds = types.MapValue(types.StringType, elements)
	resp.Diagnostics.Append(ds...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, d.data)...)
}
