package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	hostnameField = field.StringField(
		"hostname",
		field.WithRequired(true),
		field.WithDisplayName("Airbyte hostname"),
		field.WithDescription("The Airbyte hostname used to connect to the Airbyte API"),
		field.WithPlaceholder("Your Airbyte hostname"),
	)
	clientIDField = field.StringField(
		"airbyte-client-id",
		field.WithRequired(true),
		field.WithDisplayName("Airbyte client ID"),
		field.WithDescription("The Airbyte client ID used to connect to the Airbyte API."),
		field.WithPlaceholder("Your Airbyte client ID"),
		field.WithIsSecret(true),
	)
	clientSecretField = field.StringField(
		"airbyte-client-secret",
		field.WithRequired(true),
		field.WithDisplayName("Airbyte client secret"),
		field.WithDescription("The Airbyte client secret used to connect to the Airbyte API."),
		field.WithPlaceholder("Your Airbyte client secret"),
		field.WithIsSecret(true),
	)
)

var Config = field.NewConfiguration(
	[]field.SchemaField{
		hostnameField,
		clientIDField,
		clientSecretField,
	},
	field.WithConnectorDisplayName("Airbyte"),
	field.WithIconUrl("/static/app-icons/airbyte.svg"),
	field.WithHelpUrl("/docs/baton/airbyte"),
)
