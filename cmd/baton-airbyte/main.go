package main

import (
	"context"

	cfg "github.com/conductorone/baton-airbyte/pkg/config"
	"github.com/conductorone/baton-airbyte/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
)

var version = "dev"

func main() {
	ctx := context.Background()

	config.RunConnector(
		ctx,
		"baton-airbyte",
		version,
		cfg.Config,
		connector.NewLambdaConnector,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilderV2(&connector.Airbyte{}),
	)
}
