package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/william20111/terraform-provider-thousandeyes/thousandeyes"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: thousandeyes.Provider,
	})
}
