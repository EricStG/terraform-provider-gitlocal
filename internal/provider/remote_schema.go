package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func RemoteSchema(isSingle bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed:    !isSingle,
			Description: "Name of the remote",
			Required:    isSingle,
		},
		"urls": schema.ListAttribute{
			Computed:    true,
			Description: "List of remote URLs",
			ElementType: types.StringType,
		},
	}
}
