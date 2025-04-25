// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRemoteDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "gitlocal_remote" "test" { name = "origin" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlocal_remote.test", "name", "origin"),

					resource.TestCheckResourceAttr("data.gitlocal_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("data.gitlocal_remote.test", "urls.0", "https://github.com/EricStG/terraform-provider-gitlocal"),
				),
			},
		},
	})
}
