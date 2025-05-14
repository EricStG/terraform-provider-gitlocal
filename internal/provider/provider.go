// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &gitlocalProvider{}
)

type gitlocalProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type gitlocalProviderModel struct {
	Path types.String `tfsdk:"path"`
}

// Metadata returns the provider type name.
func (p *gitlocalProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gitlocal"
	resp.Version = p.version
}

func (p *gitlocalProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description: "Path to the root of the local git repository",
				Required:    true,
			},
		},
	}
}

func (p *gitlocalProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config gitlocalProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Path.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Unknown git path",
			"The provider cannot open the git repository as there is an unknown configuration value for the path. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the GIT_LOCAL_PATH environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	gitPath := os.Getenv("GIT_LOCAL_PATH")

	if !config.Path.IsNull() {
		gitPath = config.Path.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if gitPath == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("path"),
			"Missing Git Local Path",
			"The provider cannot open the git repository as there is a missing or empty value for the path. "+
				"Set the path value in the configuration or use the GIT_LOCAL_PATH environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	repo, err := git.PlainOpen(gitPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Open Git Repository",
			"An unexpected error occurred when opening the git repository. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Git Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = repo
	resp.ResourceData = repo // todo: Do I need this?
}

func (p *gitlocalProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCommitDataSource,
		NewHeadDataSource,
		NewRemoteDataSource,
		NewRemotesDataSource,
	}
}

func (p *gitlocalProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &gitlocalProvider{
			version: version,
		}
	}
}
