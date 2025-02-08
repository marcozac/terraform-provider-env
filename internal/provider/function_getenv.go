package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// NewGetenvFunction creates a new getenv function.
func NewGetenvFunction() function.Function {
	return getenvFunction{}
}

type getenvFunction struct{}

func (f getenvFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "getenv"
}

func (f getenvFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Retrieves the value of an environment variable.",
		MarkdownDescription: "This function reads the value of the specified environment variable." +
			"If the environment variable is not set - or has an empty value - and marked as required, an error will be returned." +
			"Otherwise, it returns an empty string when the variable is unset.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "name",
				MarkdownDescription: "The name of the environment variable to read",
			},
			function.BoolParameter{
				Name:                "required",
				MarkdownDescription: "Whether the environment variable is required",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f getenvFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var name string
	var required bool

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &name, &required))
	if resp.Error != nil {
		return
	}

	v := os.Getenv(name)
	if required && v == "" {
		resp.Error = function.NewArgumentFuncError(1, "Environment variable not found")
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, v))
}
