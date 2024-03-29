package main

import (
	"github.com/cloudfoundry/terraform-provider-csbpg/csbpg"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: csbpg.Provider,
	})
}
