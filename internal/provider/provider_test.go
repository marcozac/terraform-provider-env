package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var (
	// Provider-defined functions require Terraform version 1.8+.
	minTfVersion = tfversion.Version1_8_0

	protoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		ProviderName: providerserver.NewProtocol6WithError(New()()),
	}
)
