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
	_ datasource.DataSource              = &headDataSource{}
	_ datasource.DataSourceWithConfigure = &headDataSource{}
)

// NewHeadDataSource is a helper function to simplify the provider implementation.
func NewHeadDataSource() datasource.DataSource {
	return &headDataSource{}
}

// headDataSource is the data source implementation.
type headDataSource struct {
	repo *git.Repository
}

// headDataSourceModel maps the data source schema data.
type headDataSourceModel struct {
	Hash types.String `tfsdk:"hash"`
}

// coffeesIngredientsModel maps coffee ingredients data
// Metadata returns the data source type name.
func (d *headDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_head"
}

// Schema defines the schema for the data source.
func (d *headDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hash": schema.StringAttribute{
				Computed:    true,
				Description: "Hash of the commit",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *headDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state headDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	head, err := d.repo.Head()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Git Head",
			err.Error(),
		)
		return
	}

	state.Hash = types.StringValue(head.Hash().String())

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *headDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
