package connector

import (
	"context"
	"io"

	"github.com/conductorone/baton-airbyte/pkg/airbyte"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

// Airbyte represents the Baton connector for Airbyte.
type Airbyte struct {
	client *airbyte.Client
}

// ResourceSyncers returns a list of syncers for different resource types.
func (a *Airbyte) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newOrgBuilder(a.client),
		newUserBuilder(a.client),
		newWorkspaceBuilder(a.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Airbyte) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Airbyte) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Airbyte Baton Connector",
		Description: "Connector syncing Airbyte organizations and users to Baton",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Airbyte) Validate(ctx context.Context) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	_, err := d.client.ListOrganizations(ctx)
	if err != nil {
		l.Error("Error listing organizations", zap.Error(err))
		return nil, err
	}

	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, hostname string, clientId string, clientSecret string) (*Airbyte, error) {
	airbyteClient, err := airbyte.NewClient(ctx, hostname, clientId, clientSecret)
	if err != nil {
		l := ctxzap.Extract(ctx)
		l.Error("Error creating Airbyte client", zap.Error(err))
		return nil, err
	}

	return &Airbyte{
		client: airbyteClient,
	}, nil
}
