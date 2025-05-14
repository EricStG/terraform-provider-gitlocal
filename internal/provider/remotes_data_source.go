// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &remotesDataSource{}
	_ datasource.DataSourceWithConfigure = &remotesDataSource{}
)

// NewRemotesDataSource is a helper function to simplify the provider implementation.
func NewRemotesDataSource() datasource.DataSource {
	return &remotesDataSource{}
}

// remotesDataSource is the data source implementation.
type remotesDataSource struct {
	repo *git.Repository
}

// coffeesDataSourceModel maps the data source schema data.
type remotesDataSourceModel struct {
	Remotes []remotesModel `tfsdk:"remotes"`
}

// remotesModel maps coffees schema data.
type remotesModel struct {
	Name types.String   `tfsdk:"name"`
	Urls []types.String `tfsdk:"urls"`
}

// Metadata returns the data source type name.
func (d *remotesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remotes"
}

// Schema defines the schema for the data source.
func (d *remotesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"remotes": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of remotes in the repository",
				NestedObject: schema.NestedAttributeObject{
					Attributes: RemoteSchema(false),
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *remotesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state remotesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	remotes, err := d.repo.Remotes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Git Remotes",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, remote := range remotes {
		config := remote.Config()
		remotesState := remotesModel{
			Name: types.StringValue(config.Name),
		}

		for _, url := range config.URLs {
			remotesState.Urls = append(remotesState.Urls, types.StringValue(url))
		}

		state.Remotes = append(state.Remotes, remotesState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *remotesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	repo, ok := req.ProviderData.(*git.Repository)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *git.Repository, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.repo = repo
}
