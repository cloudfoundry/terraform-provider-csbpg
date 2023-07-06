testing-github-actions
# terraform-provider-csbpg

This is a highly specialised Terraform provider designed to be used exclusively with the [Cloud Service Broker](https://github.com/cloudfoundry-incubator/cloud-service-broker) ("CSB"). Initially CSB brokerpaks used other PostgreSQL Terraform providers, but some features proved problematic. In particular:
- This provider solves the problem of mapping service bindings to PostgreSQL roles. When a binding is deleted, all the object that it owns are re-assigned to a persistent Data Owner Role.
- This provider handles TLS certificates as Terraform strings and not as files. This means that it can be used without creating local files to hold certificates, which solves some circular lifecycle issues around provider configuration.

## Usage
```terraform
terraform {
  required_providers {
    csbpg = {
      source  = ".../csbpg"
      version = "1.0.0"
    }
  }
}

provider "csbpg" {
  host            = "localhost"
  port            = 5432
  username        = "postgres"
  password        = "fill-me-in"
  database        = "postgres"
  data_owner_role = "dataowner"
}

resource "csbpg_binding_user" "binding_user" {
  username = "foo"
  password = "bar"
}
```

## Releasing
To create a new GitHub release, decide on a new version number [according to Semanitc Versioning](https://semver.org/), and then:
1. Create a tag on the main branch with a leading `v`:
   `git tag vX.Y.X`
1. Push the tag:
   `git push --tags`
1. Wait for the GitHub action to run GoReleaser and create the new GitHub release

