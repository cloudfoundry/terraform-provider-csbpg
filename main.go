package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/cloudfoundry/terraform-provider-csbpg/csbpg"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: csbpg.Provider,
	})
}
