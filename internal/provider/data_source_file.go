package provider

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	Path     types.String `tfsdk:"path"`
	Required types.Bool   `tfsdk:"required"`
	Result   types.Map    `tfsdk:"result"`
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
			"required": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the file is required.",
				MarkdownDescription: "Whether the file is required." +
					"Prevents the data source from returning an error if the file is not found.",
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

	f, ok := d.open(resp, d.data.Path.ValueString())
	if !ok {
		d.setResult(ctx, resp, make(map[string]attr.Value))
		return
	}
	defer f.Close()

	env, err := godotenv.Parse(f)
	if err != nil {
		resp.Diagnostics.AddError("Failed To Parse File", err.Error())
		return
	}

	m := make(map[string]attr.Value, len(env))
	for k, v := range env {
		m[k] = types.StringValue(v)
	}
	d.setResult(ctx, resp, m)
}

func (d *fileDataSource) open(resp *datasource.ReadResponse, filename string) (f *os.File, ok bool) {
	f, err := os.Open(filename)
	if err == nil {
		return f, true
	}
	switch {
	case !errors.Is(err, fs.ErrNotExist):
		resp.Diagnostics.AddError("Failed To Open File", err.Error())
	case d.data.Required.ValueBool():
		resp.Diagnostics.AddAttributeError(path.Root("path"),
			"File Not Found",
			fmt.Sprintf("File %q not found", filename),
		)
	default:
		resp.Diagnostics.AddAttributeWarning(path.Root("path"),
			"File Not Found",
			fmt.Sprintf("File %q not found, returning empty result", filename),
		)
	}
	return
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
