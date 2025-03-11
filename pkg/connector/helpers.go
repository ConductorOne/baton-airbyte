package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
)

const ResourcesPageSize uint64 = 50

// parsePageToken unmarshals the token and initializes pagination state,
// creating a new state if none exists. Returns bag, token string, and error.
func parsePageToken(pagToken *pagination.Token, resourceID *v2.ResourceId) (*pagination.Bag, string, error) {
	bag := &pagination.Bag{}
	err := bag.Unmarshal(pagToken.Token)
	if err != nil {
		return nil, "", err
	}

	if bag.Current() == nil {
		// If no current page state, push a new one for the provided resource.
		bag.Push(pagination.PageState{
			ResourceTypeID: resourceID.ResourceType,
			ResourceID:     resourceID.Resource,
		})
	}

	return bag, bag.PageToken(), nil
}
