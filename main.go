package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/joreichhardt/terraform-provider-centron/internal/provider"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/joreichhardt/centron",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
