terraform {
  required_providers {
    gitlocal = {
      source = "registry.terraform.io/ericstg/gitlocal"
    }
  }
}

provider "gitlocal" {
  path = "../../"
}

data "gitlocal_remotes" "example" {}

output "git_remotes" {
  value = data.gitlocal_remotes.example
}
