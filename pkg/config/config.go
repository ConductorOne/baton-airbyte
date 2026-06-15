package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var Config = field.NewConfiguration(
	[]field.SchemaField{
		field.StringField(
			"hostname",
			field.WithRequired(true),
			field.WithDescription("The Airbyte hostname used to connect to the Airbyte API"),
		),
		field.StringField(
			"airbyte-client-id",
			field.WithRequired(true),
			field.WithDescription("The Airbyte client id used to connect to the Airbyte API."),
		),
		field.StringField(
			"airbyte-client-secret",
			field.WithRequired(true),
			field.WithDescription("The Airbyte client secret used to connect to the Airbyte API."),
			field.WithIsSecret(true),
		),
	},
	field.WithConstraints(
		field.FieldsRequiredTogether(
			field.StringField("airbyte-client-id"),
			field.StringField("airbyte-client-secret"),
		),
	),
)

func ValidateConfig(c *Airbyte) error {
	return nil
}
