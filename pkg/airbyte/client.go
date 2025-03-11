package airbyte

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type Client struct {
	baseURL      *url.URL
	accessToken  string
	clientID     string
	clientSecret string
	httpClient   *uhttp.BaseHttpClient
	tokenExpiry  time.Time
}

const (
	getAccessTokenPath               = "/api/v1/applications/token" // #nosec G101
	getWorkspacePath                 = "/api/public/v1/workspaces/{workspaceId}"
	listWorkspacesPath               = "/api/public/v1/workspaces"
	listUsersPath                    = "/api/public/v1/users"
	listOrganizationsPath            = "/api/public/v1/organizations"
	listPermissionsPath              = "/api/public/v1/permissions"
	listWorkspacesByOrganizationPath = "/api/v1/workspaces/list_by_organization_id"
	listUsersWithAccessInfoPath      = "/api/v1/users/list_access_info_by_workspace_id"
)

func NewClient(ctx context.Context, hostname string, clientID string, clientSecret string) (*Client, error) {
	baseURL, err := url.Parse(hostname)
	if err != nil {
		return nil, err
	}

	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	wrapper, err := uhttp.NewBaseHttpClientWithContext(ctx, httpClient)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient:   wrapper,
		baseURL:      baseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}

// Access token lifetimes vary by Airbyte deployment type/version:
// • Open Source/Cloud: 3 minutes
// • Enterprise: 24 hours
//
// This function ensures that the access token is valid before making requests.
// It refreshes the token if it's expired or about to expire.
//
// The token is refreshed when:
// • The token is not set (first time access)
// • The token is expired (3 minutes/24 hours)
// • The token expires in the next 30 seconds
//
// This ensures that the token is always fresh when needed.
//
// Reference: https://reference.airbyte.com/reference/authentication
func (c *Client) ensureValidToken(ctx context.Context) error {
	// Check if token needs refresh (with 30s buffer).
	if c.accessToken == "" || time.Now().Add(30*time.Second).After(c.tokenExpiry) {
		// Get new token.
		token, expiry, err := c.GetAccessToken(ctx)
		if err != nil {
			return err
		}

		c.accessToken = token
		c.tokenExpiry = expiry
	}

	return nil
}

// -------------------------------------------------------------------------------------------------
// PUBLIC API ENDPOINTS
// -------------------------------------------------------------------------------------------------

// GetAccessToken fetches a new access token from Airbyte.
//
// This function handles the process of obtaining a new access token from Airbyte.
// It constructs the request URL, sets up the request body, and makes the HTTP POST request.
//
// The token response is parsed to extract the access token and its expiration time.
//
// The function returns the new access token and its expiration time.
func (c *Client) GetAccessToken(ctx context.Context) (string, time.Time, error) {
	tokenResp := &TokenResponse{}

	body := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
	}

	err := c.doRequest(ctx, http.MethodPost, c.buildResourceURL(getAccessTokenPath, nil, nil), tokenResp, body, true)
	if err != nil {
		return "", time.Time{}, err
	}

	// Parse JWT token to get expiry
	parts := strings.Split(tokenResp.AccessToken, ".")
	if len(parts) != 3 {
		return "", time.Time{}, fmt.Errorf("invalid JWT token format")
	}

	// Decode the claims (middle part)
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error decoding JWT claims: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return "", time.Time{}, fmt.Errorf("error parsing JWT claims: %w", err)
	}

	expiry := time.Unix(claims.ExpiresAt, 0)
	return tokenResp.AccessToken, expiry, nil
}

// ListAllWorkspaces fetches all workspaces from Airbyte.
//
// This function retrieves all workspaces available in the Airbyte system.
// It uses pagination to handle large datasets efficiently.
//
// The function returns a list of workspaces and the offset for the next page of workspaces.
func (c *Client) ListAllWorkspaces(ctx context.Context, limit uint64, offset string) ([]*WorkspaceResponse, string, error) {
	resp := &APIResponse[[]*WorkspaceResponse]{}

	// If offset is empty, set it to 0.
	if offset == "" {
		offset = "0"
	}

	queryParams := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": offset,
	}

	err := c.doRequest(ctx, http.MethodGet, c.buildResourceURL(listWorkspacesPath, nil, queryParams), resp, nil, false)
	if err != nil {
		return nil, "", err
	}

	return resp.Data, GetOffsetForTheNextPageFromURL(resp.Next), nil
}

// ListUsersByOrganization fetches users by organization from Airbyte.
//
// This function retrieves users associated with a specific organization.
//
// The function returns a list of users.
func (c *Client) ListUsersByOrganization(ctx context.Context, orgId string) ([]*User, error) {
	resp := &APIResponse[[]*User]{}

	queryParams := map[string]string{
		"organizationId": orgId,
	}

	// This endpoint doesn't support pagination.
	err := c.doRequest(ctx, http.MethodGet, c.buildResourceURL(listUsersPath, nil, queryParams), resp, nil, false)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// ListPermissionsByUserAndOrganization fetches permissions by user and organization from Airbyte.
//
// This function retrieves permissions associated with a specific user and organization.
//
// The function returns a list of permissions.
func (c *Client) ListPermissionsByUserAndOrganization(ctx context.Context, userId string, orgId string) ([]*Permission, error) {
	resp := &APIResponse[[]*Permission]{}

	pathParams := map[string]string{
		"userId": userId,
		"orgId":  orgId,
	}

	// This endpoint doesn't support pagination.
	err := c.doRequest(ctx, http.MethodGet, c.buildResourceURL(listPermissionsPath, pathParams, nil), resp, nil, false)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// ListOrganizations fetches all organizations from Airbyte.
//
// This function retrieves all organizations available in the Airbyte.
//
// The function returns a list of organizations.
func (c *Client) ListOrganizations(ctx context.Context) ([]*Organization, error) {
	resp := &APIResponse[[]*Organization]{}

	// This endpoint doesn't support pagination.
	err := c.doRequest(ctx, http.MethodGet, c.buildResourceURL(listOrganizationsPath, nil, nil), resp, nil, false)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// -------------------------------------------------------------------------------------------------
// PRIVATE API ENDPOINTS
// -------------------------------------------------------------------------------------------------

// ListWorkspacesByOrganization fetches workspaces by organization from Airbyte.
//
// This function retrieves workspaces associated with a specific organization.
// It uses pagination to handle large datasets efficiently.
//
// The function returns a list of workspaces.
func (c *Client) ListWorkspacesByOrganization(ctx context.Context, orgId string, pageSize uint64, rowOffset uint64) ([]WorkspaceReadResponse, uint64, error) {
	resp := &WorkspaceReadListResponse{}

	body := map[string]interface{}{
		"organizationId": orgId,
		"pagination": map[string]interface{}{
			"pageSize":  pageSize,
			"rowOffset": rowOffset,
		},
	}

	err := c.doRequest(ctx, http.MethodPost, c.buildResourceURL(listWorkspacesByOrganizationPath, nil, nil), resp, body, false)
	if err != nil {
		return nil, 0, err
	}

	if uint64(len(resp.Workspaces)) < pageSize {
		return resp.Workspaces, 0, nil
	}

	nextRowOffset := rowOffset + pageSize

	return resp.Workspaces, nextRowOffset, nil
}

// ListUsersWithAccessInfoByWorkspace fetches users with access info by workspace from Airbyte.
//
// This function retrieves users with access info (workspace and organization permission type) associated with a particular workspace.
//
// The function returns a list of users with access info.
func (c *Client) ListUsersWithAccessInfoByWorkspace(ctx context.Context, workspaceId string) ([]WorkspaceUserAccessInfoReadResponse, error) {
	resp := &WorkspaceUserAccessInfoReadListResponse{}

	body := map[string]string{
		"workspaceId": workspaceId,
	}

	// This endpoint doesn't support pagination.
	err := c.doRequest(ctx, http.MethodPost, c.buildResourceURL(listUsersWithAccessInfoPath, nil, nil), resp, body, false)
	if err != nil {
		return nil, err
	}

	return resp.UsersWithAccess, nil
}

// -------------------------------------------------------------------------------------------------
// PRIVATE HELPER FUNCTIONS
// -------------------------------------------------------------------------------------------------

// doRequest handles HTTP requests with authentication and optional pagination.
//
// This function constructs a request with the specified HTTP method, URL, and optional data.
// It also handles authentication by adding an authorization header if not skipping authentication.
//
// The function returns an error if the request fails or if the response cannot be parsed.
func (c *Client) doRequest(
	ctx context.Context,
	method string,
	urlAddress *url.URL,
	response interface{},
	data interface{},
	skipAuth bool,
) error {
	reqOptions := []uhttp.RequestOption{
		uhttp.WithContentType("application/json"),
		uhttp.WithAccept("application/json"),
	}

	// Only add authorization header if not skipping auth.
	if !skipAuth {
		if err := c.ensureValidToken(ctx); err != nil {
			return err
		}
		reqOptions = append(reqOptions, uhttp.WithHeader("Authorization", "Bearer "+c.accessToken))
	}

	if data != nil {
		reqOptions = append(reqOptions, uhttp.WithJSONBody(data))
	}

	req, err := c.httpClient.NewRequest(ctx, method, urlAddress, reqOptions...)
	if err != nil {
		return err
	}

	doOptions := []uhttp.DoOption{}
	if response != nil {
		doOptions = append(doOptions, uhttp.WithJSONResponse(response))
	}

	resp, err := c.httpClient.Do(req, doOptions...)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

// The buildResourceURL function constructs an absolute URL by formatting a resource path.
//
// This function constructs a URL by replacing path parameters with their actual values and adding query parameters.
//
// The function returns the constructed URL.
// Example:
// pathTemplate: "/api/v1/workspaces/{workspaceId}"
// pathParams: map[string]string{"workspaceId": "123"}
// queryParams: map[string]string{"limit": "10", "offset": "0"}
// The function returns the constructed URL: "/api/v1/workspaces/123?limit=10&offset=0".
func (c *Client) buildResourceURL(pathTemplate string, pathParams map[string]string, queryParams map[string]string) *url.URL {
	finalPath := pathTemplate

	// Replace path parameters using named placeholders
	if len(pathParams) > 0 {
		for key, value := range pathParams {
			placeholder := fmt.Sprintf("{%s}", key)
			finalPath = strings.ReplaceAll(finalPath, placeholder, url.PathEscape(value))
		}
	}

	// Create URL from base and path
	u := c.baseURL.ResolveReference(&url.URL{Path: finalPath})

	// Add query parameters if provided
	if len(queryParams) > 0 {
		q := u.Query()
		for key, value := range queryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	return u
}

// GetOffsetForTheNextPageFromURL extracts the offset from a URL.
//
// This function parses a URL and extracts the offset parameter from the query string.
//
// The function returns the extracted offset value.
// Example:
// urlPayload: "https://api.airbyte.com/v1/workspaces?includeDeleted=false&limit=10&offset=0"
// The function returns the extracted offset value: "10" in case that the next page URL is: "https://api.airbyte.com/v1/workspaces?includeDeleted=false&limit=10&offset=10"
func GetOffsetForTheNextPageFromURL(urlPayload string) string {
	if urlPayload == "" {
		return ""
	}

	u, err := url.Parse(urlPayload)
	if err != nil {
		return ""
	}

	offset := u.Query().Get("offset")
	if offset == "" {
		return ""
	}

	return offset
}
