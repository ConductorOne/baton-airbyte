package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	HostnameField = field.StringField(
		"hostname",
		field.WithRequired(true),
		field.WithDescription("The Airbyte hostname used to connect to the Airbyte API"),
	)
	ClientIdField = field.StringField(
		"airbyte-client-id",
		field.WithRequired(true),
		field.WithDescription("The Airbyte client id used to connect to the Airbyte API."),
	)
	ClientSecretField = field.StringField(
		"airbyte-client-secret",
		field.WithRequired(true),
		field.WithDescription("The Airbyte client secret used to connect to the Airbyte API."),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run. Note: these fields can be marked as optional or
	// required.
	ConfigurationFields = []field.SchemaField{HostnameField, ClientIdField, ClientSecretField}

	// FieldRelationships defines relationships between the fields listed in
	// ConfigurationFields that can be automatically validated. For example, a
	// username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{
		field.FieldsRequiredTogether(ClientIdField, ClientSecretField),
	}

	// Config is the configuration schema for the connector.
	Config = field.Configuration{
		Fields:      ConfigurationFields,
		Constraints: FieldRelationships,
	}
)
