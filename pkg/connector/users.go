package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-airbyte/pkg/airbyte"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *airbyte.Client
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// Create a new connector resource for an Airbyte user.
func userResource(user *airbyte.User) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	}

	userTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(v2.UserTrait_Status_STATUS_ENABLED),
		rs.WithEmail(user.Email, true),
	}

	resource, err := rs.NewUserResource(
		user.Email,
		userResourceType,
		user.ID,
		userTraitOptions,
	)

	if err != nil {
		return nil, err
	}

	return resource, nil
}

// List returns all the users as resource objects.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	// The only way found to list all users was the list users endpoint with access information per workspace, since the list users endpoint does not work as we might expect..
	// If we use the endpoint to list users by organization, we would lose the users who only have access to a single workspace.
	ListUserResponse, err := o.client.ListUsersWithAccessInfoByWorkspace(ctx, parentResourceID.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("airbyte-connector: failed to list users: %w", err)
	}

	resources := make([]*v2.Resource, 0, len(ListUserResponse))
	// Convert users to resources
	for _, userResponse := range ListUserResponse {
		user := airbyte.User{
			ID:    userResponse.UserID,
			Email: userResponse.UserEmail,
			Name:  userResponse.UserName,
		}
		ur, err := userResource(&user)

		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to create resource for user %s: %w", user.Email, err)
		}
		resources = append(resources, ur)
	}

	return resources, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *airbyte.Client) *userBuilder {
	return &userBuilder{
		resourceType: userResourceType,
		client:       client,
	}
}
