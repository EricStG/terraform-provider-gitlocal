// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &commitDataSource{}
	_ datasource.DataSourceWithConfigure = &commitDataSource{}
)

// NewCommitDataSource is a helper function to simplify the provider implementation.
func NewCommitDataSource() datasource.DataSource {
	return &commitDataSource{}
}

// commitDataSource is the data source implementation.
type commitDataSource struct {
	repo *git.Repository
}

// commitDataSourceModel maps the data source schema data.
type commitDataSourceModel struct {
	Date    types.String `tfsdk:"date"`
	Hash    types.String `tfsdk:"hash"`
	Message types.String `tfsdk:"message"`
}

// Metadata returns the data source type name.
func (d *commitDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_commit"
}

// Schema defines the schema for the data source.
func (d *commitDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hash": schema.StringAttribute{
				Description: "Hash of the commit",
				Required:    true,
			},
			"date": schema.StringAttribute{
				Computed:    true,
				Description: "Date of the commit in RFC 3339",
			},
			"message": schema.StringAttribute{
				Computed:    true,
				Description: "Message of the commit",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *commitDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state commitDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	hashArg := state.Hash.ValueString()
	hash := plumbing.NewHash((hashArg))

	commit, err := d.repo.CommitObject(hash)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Commit `"+hashArg+"`",
			err.Error(),
		)
		return
	}

	state.Date = types.StringValue(commit.Author.When.Format(time.RFC3339))
	state.Message = types.StringValue(commit.Message)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *commitDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
