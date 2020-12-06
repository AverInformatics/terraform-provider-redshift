package main

import (
	"github.com/joshuarose/terraform-provider-redshift/redshift"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: redshift.Provider})
}
