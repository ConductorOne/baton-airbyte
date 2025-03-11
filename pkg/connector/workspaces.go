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

// Define workspace permission type constants.
//
//	Reference link to permission types: "https://github.com/airbytehq/airbyte-api-python-sdk/blob/main/src/airbyte_api/models/publicpermissiontype.py".
const (
	WorkspaceAdmin  = "workspace_admin"
	WorkspaceEditor = "workspace_editor"
	WorkspaceRunner = "workspace_runner"
	WorkspaceReader = "workspace_reader"
)

var PublicWorkspacePermissionsTypes = []string{
	WorkspaceAdmin,
	WorkspaceEditor,
	WorkspaceRunner,
	WorkspaceReader,
}

// workspacesWithOrgIDMap maps workspace IDs to their corresponding organization IDs.
var workspacesWithOrgIDMap map[string]string

type workspaceBuilder struct {
	resourceType *v2.ResourceType
	client       *airbyte.Client
}

func (o *workspaceBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return workspaceResourceType
}

// Create a new connector resource for an airbyte workspace.
func workspaceResource(workspace airbyte.Workspace, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	resource, err := rs.NewResource(
		workspace.Name,
		workspaceResourceType,
		workspace.ID,
		rs.WithAnnotation(
			&v2.ChildResourceType{
				ResourceTypeId: userResourceType.Id,
			},
		),
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

// List returns all workspaces and their parent organization IDs (when available).
// The process requires two API calls:
//  1. GET /api/public/v1/workspaces
//     Returns all workspaces but without organization IDs
//  2. GET /api/v1/workspaces/list_by_organization_id
//     Returns workspaces into the accessible organizations only
//
// Workspaces belonging to organizations we can't access will be marked with an
// "unknown-parent" organization ID.
func (o *workspaceBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	// Initialize the map if we're starting a new list
	if pToken.Token == "" {
		// Initialize the map
		workspacesWithOrgIDMap = make(map[string]string)

		// Get workspaces with organizations
		allWorkspacesWithParentOrganizationID, err := o.getAllWorkspacesWithParentOrganizationID(ctx)
		if err != nil {
			return nil, "", nil, fmt.Errorf("airbyte-connector: getAllWorkspacesWithParentOrganizationID > failed to list workspaces: %w", err)
		}

		// Populate the map with workspace-to-organization relationships
		// The public workspaces endpoint doesn't include organization IDs,
		// so we maintain this mapping to associate workspaces with their parent organizations
		for _, w := range allWorkspacesWithParentOrganizationID {
			if w.OrganizationId != "" {
				workspacesWithOrgIDMap[w.ID] = w.OrganizationId
			}
		}
	}

	// pToken.Token is the offset for the current page
	bag, offsetForCurrentPage, err := parsePageToken(pToken, &v2.ResourceId{ResourceType: workspaceResourceType.Id})
	if err != nil {
		return nil, "", nil, err
	}

	listWorkspaceResponse, offsetForNextPage, err := o.client.ListAllWorkspaces(ctx, ResourcesPageSize, offsetForCurrentPage)
	if err != nil {
		return nil, "", nil, fmt.Errorf("airbyte-connector: ListAllWorkspaces > failed to list workspaces: %w", err)
	}

	next, err := bag.NextToken(offsetForNextPage)
	if err != nil {
		return nil, "", nil, err
	}

	// Process all workspaces
	resources := make([]*v2.Resource, 0, len(listWorkspaceResponse))
	for _, ws := range listWorkspaceResponse {
		workspace := airbyte.Workspace{
			ID:   ws.ID,
			Name: ws.Name,
		}

		var parentResourceID *v2.ResourceId
		// Only set parent resource ID if we have a valid organization ID
		if orgID, exists := workspacesWithOrgIDMap[ws.ID]; exists && orgID != "" {
			workspace.OrganizationId = orgID
			parentResourceID = &v2.ResourceId{
				ResourceType: organizationResourceType.Id,
				Resource:     orgID,
			}
		} else {
			// If we don't have a valid organization ID, we set the parent resource ID to "unknown-parent"
			// This is used to indicate that the workspace is associated with an organization that we don't have access to.
			parentResourceID = &v2.ResourceId{
				ResourceType: organizationResourceType.Id,
				Resource:     "unknown-parent",
			}
		}

		resource, err := workspaceResource(workspace, parentResourceID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to create resource for workspace %s: %w", workspace.Name, err)
		}

		resources = append(resources, resource)
	}

	return resources, next, nil, nil
}

// Entitlements returns a slice of entitlements for possible user roles under workspace (Viewer, Editor, Admin).
func (o *workspaceBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	// Preallocate slice for efficiency
	entitlements := make([]*v2.Entitlement, 0, len(PublicWorkspacePermissionsTypes))

	for _, permissionType := range PublicWorkspacePermissionsTypes {
		// Generate display name and description
		displayName := fmt.Sprintf("%s %s", resource.DisplayName, permissionType)
		description := fmt.Sprintf("%s role in %s Airbyte workspace", permissionType, resource.DisplayName)

		// Define entitlement options
		entitlementOptions := []ent.EntitlementOption{
			ent.WithGrantableTo(workspaceResourceType),
			ent.WithDisplayName(displayName),
			ent.WithDescription(description),
		}

		// Append new entitlement to the slice
		entitlements = append(entitlements, ent.NewPermissionEntitlement(resource, permissionType, entitlementOptions...))
	}

	return entitlements, "", nil, nil
}

// Grants returns a slice of grants for each user and their set role under workspace.
func (o *workspaceBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	listUserswithaccessInfoResponse, err := o.client.ListUsersWithAccessInfoByWorkspace(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("airbyte-connector: failed to list users under workspace %s: %w", resource.Id.Resource, err)
	}

	// Map organization permissions to workspace permissions.
	// We use this mapping because organization permissions propagate down to workspaces.
	// In some cases, users may only have organization-level permissions set without explicit
	// workspace permissions. This mapping ensures we correctly reflect the inherited
	// permissions at the workspace level.
	orgToWorkspacePermMap := map[string]string{
		"organization_admin":  WorkspaceAdmin,
		"organization_editor": WorkspaceEditor,
		"organization_runner": WorkspaceRunner,
		"organization_reader": WorkspaceReader,
	}

	var rv []*v2.Grant
	for _, userResponse := range listUserswithaccessInfoResponse {
		var permissionType string

		// Prefer workspace-level permission over organization-level
		if userResponse.WorkspacePermission != nil {
			permissionType = strings.ToLower(userResponse.WorkspacePermission.PermissionType)
		} else if userResponse.OrganizationPermission != nil {
			orgPermType := strings.ToLower(userResponse.OrganizationPermission.PermissionType)
			// Map organization permission to workspace permission
			if mappedPerm, exists := orgToWorkspacePermMap[orgPermType]; exists {
				permissionType = mappedPerm
			}
		}

		// Skip if no valid permission type found
		if !slices.Contains(PublicWorkspacePermissionsTypes, permissionType) && !slices.Contains(PublicOrganizationPermissionsTypes, permissionType) {
			continue
		}

		user := airbyte.User{
			ID:    userResponse.UserID,
			Email: userResponse.UserEmail,
			Name:  userResponse.UserName,
		}

		userResource, err := userResource(&user)
		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, grant.NewGrant(resource, permissionType, userResource.Id))
	}

	return rv, "", nil, nil
}

func newWorkspaceBuilder(client *airbyte.Client) *workspaceBuilder {
	return &workspaceBuilder{
		resourceType: workspaceResourceType,
		client:       client,
	}
}

// -------------------------------------------------------------------------------------------------
// PRIVATE HELPER FUNCTIONS
// -------------------------------------------------------------------------------------------------

// getAllWorkspacesWithParentOrganizationID retrieves all workspaces and their associated organization IDs
// by iterating through each organization and fetching its workspaces. This is necessary because the public
// workspace API endpoint doesn't provide organization information, but we need this relationship for proper
// resource hierarchy mapping.
func (o *workspaceBuilder) getAllWorkspacesWithParentOrganizationID(ctx context.Context) ([]*airbyte.Workspace, error) {
	allWorkspacesWithParentOrganizationID := make([]*airbyte.Workspace, 0)

	orgs, err := o.client.ListOrganizations(ctx)
	if err != nil {
		return nil, fmt.Errorf("airbyte-connector: failed to list organizations: %w", err)
	}

	for _, org := range orgs {
		var rowOffset uint64 = 0
		for {
			listWorkspaceReadResponse, nextRowOffset, err := o.client.ListWorkspacesByOrganization(ctx, org.ID, ResourcesPageSize, rowOffset)
			if err != nil {
				return nil, fmt.Errorf("airbyte-connector: failed to list workspaces: %w", err)
			}

			for _, workspaceReadResponse := range listWorkspaceReadResponse {
				workspace := &airbyte.Workspace{
					ID:             workspaceReadResponse.WorkspaceId,
					Name:           workspaceReadResponse.Name,
					OrganizationId: workspaceReadResponse.OrganizationId,
				}
				allWorkspacesWithParentOrganizationID = append(allWorkspacesWithParentOrganizationID, workspace)
			}

			if nextRowOffset == 0 {
				break
			}

			rowOffset = nextRowOffset
		}
	}

	return allWorkspacesWithParentOrganizationID, nil
}
