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
	_ datasource.DataSource              = &remoteDataSource{}
	_ datasource.DataSourceWithConfigure = &remoteDataSource{}
)

// NewRemoteDataSource is a helper function to simplify the provider implementation.
func NewRemoteDataSource() datasource.DataSource {
	return &remoteDataSource{}
}

// remoteDataSource is the data source implementation.
type remoteDataSource struct {
	repo *git.Repository
}

// coffeesDataSourceModel maps the data source schema data.
type remoteDataSourceModel struct {
	Name types.String   `tfsdk:"name"`
	Urls []types.String `tfsdk:"urls"`
}

// coffeesIngredientsModel maps coffee ingredients data
// Metadata returns the data source type name.
func (d *remoteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote"
}

// Schema defines the schema for the data source.
// todo: figure out how to share schema
func (d *remoteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: RemoteSchema(true),
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *remoteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state remoteDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	remoteName := state.Name.ValueString()

	remote, err := d.repo.Remote(remoteName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Git Remote `"+remoteName+"`",
			err.Error(),
		)
		return
	}

	for _, url := range remote.Config().URLs {
		state.Urls = append(state.Urls, types.StringValue(url))
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *remoteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
