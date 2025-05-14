// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCommitDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "gitlocal_commit" "test" { hash = "13981df24be0f47a94044d21bacdbc3d9132f162" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("data.gitlocal_commit.test", "date", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2025-05-14T09:08:28-04:00")
						parsed, err := time.Parse(time.RFC3339, value)

						if err != nil {
							return err
						}

						if !expected.Equal(parsed) {
							return errors.New("Date did not match expectation")
						}

						return nil
					}),
					resource.TestCheckResourceAttr("data.gitlocal_commit.test", "message", "Add head data source\n"),
				),
			},
		},
	})
}
