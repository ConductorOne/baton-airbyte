package connector

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/conductorone/baton-airbyte/pkg/airbyte"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

// Define organization permission type constants.
//
//	Reference link to permission types: https://github.com/airbytehq/airbyte-api-python-sdk/blob/main/src/airbyte_api/models/publicpermissiontype.py
const (
	OrganizationAdmin  = "organization_admin"
	OrganizationEditor = "organization_editor"
	OrganizationRunner = "organization_runner"
	OrganizationReader = "organization_reader"
	OrganizationMember = "organization_member"
)

var PublicOrganizationPermissionsTypes = []string{
	OrganizationAdmin,
	OrganizationEditor,
	OrganizationRunner,
	OrganizationReader,
	OrganizationMember,
}

type orgBuilder struct {
	resourceType *v2.ResourceType
	client       *airbyte.Client
}

func (o *orgBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return organizationResourceType
}

// Create a new connector resource for an airbyte organization.
func orgResource(org airbyte.Organization) (*v2.Resource, error) {
	resource, err := rs.NewResource(
		org.Name,
		organizationResourceType,
		org.ID,
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

// List returns all the organizations.
func (o *orgBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	orgs, err := o.client.ListOrganizations(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("airbyte-connector: failed to list organizations: %w", err)
	}

	// Iterate over organizations and filter valid ones
	resources := make([]*v2.Resource, 0, len(orgs))
	for _, org := range orgs {
		org := airbyte.Organization{
			ID:   org.ID,
			Name: org.Name,
		}
		// Convert organization to a v2.Resource
		resource, err := orgResource(org)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to create resource for organization %s: %w", org.Name, err)
		}

		resources = append(resources, resource)
	}

	return resources, "", nil, nil
}

// Entitlements returns a slice of entitlements for possible user roles under organization.
func (o *orgBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	// Preallocate slice for efficiency
	entitlements := make([]*v2.Entitlement, 0, len(PublicOrganizationPermissionsTypes))

	for _, permissionType := range PublicOrganizationPermissionsTypes {
		// Generate display name and description
		displayName := fmt.Sprintf("%s %s", resource.DisplayName, permissionType)
		description := fmt.Sprintf("%s role in %s Airbyte organization", permissionType, resource.DisplayName)

		// Define entitlement options
		entitlementOptions := []ent.EntitlementOption{
			ent.WithGrantableTo(userResourceType),
			ent.WithDisplayName(displayName),
			ent.WithDescription(description),
		}

		// Append new entitlement to the slice
		entitlements = append(entitlements, ent.NewPermissionEntitlement(resource, permissionType, entitlementOptions...))
	}

	return entitlements, "", nil, nil
}

// Grants returns a slice of grants for each user and their set role under organization.
func (o *orgBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, err := o.client.ListUsersByOrganization(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("airbyte-connector: failed to list users under organization %s: %w", resource.Id.Resource, err)
	}

	var rv []*v2.Grant
	for _, user := range users {
		// Get the permission type for the user under the organization
		permissionType, err := o.getOrganizationPermissionType(ctx, user.ID, resource.Id.Resource)
		if err != nil {
			return nil, "", nil, fmt.Errorf("airbyte-connector: failed to get permission type for user %s under organization %s: %w", user.ID, resource.Id.Resource, err)
		}

		// check for valid roles and skip if not
		if !slices.Contains(PublicOrganizationPermissionsTypes, permissionType) {
			continue
		}

		userResource, err := userResource(user)
		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, grant.NewGrant(resource, permissionType, userResource.Id))
	}

	return rv, "", nil, nil
}

func newOrgBuilder(client *airbyte.Client) *orgBuilder {
	return &orgBuilder{
		resourceType: organizationResourceType,
		client:       client,
	}
}

// -------------------------------------------------------------------------------------------------
// PRIVATE HELPER FUNCTIONS
// -------------------------------------------------------------------------------------------------

func (o *orgBuilder) getOrganizationPermissionType(ctx context.Context, userID, organizationID string) (string, error) {
	var allPermissions []*airbyte.Permission

	permissions, err := o.client.ListPermissionsByUserAndOrganization(ctx, userID, organizationID)
	if err != nil {
		return "", fmt.Errorf("airbyte-connector: failed to list permissions for user %s: %w", userID, err)
	}

	allPermissions = append(allPermissions, permissions...)

	// Find permission for this organization
	for _, permission := range allPermissions {
		if permission.Scope == "organization" && permission.ScopeID == organizationID {
			return strings.ToLower(permission.PermissionType), nil
		}
	}

	return "", nil
}
