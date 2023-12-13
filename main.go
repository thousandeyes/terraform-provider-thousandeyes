package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugModeTest2 bool

	flag.BoolVar(&debugModeTest2, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := &plugin.ServeOpts{
		Debug:        debugModeTest2,
		ProviderAddr: "registry.terraform.io/thousandeyes/thousandeyes",
		ProviderFunc: thousandeyes.New(),
	}

	plugin.Serve(opts)
}
