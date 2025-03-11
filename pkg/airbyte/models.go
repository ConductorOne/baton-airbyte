package airbyte

type Workspace struct {
	ID             string
	OrganizationId string
	Name           string
}

// ------------------------------------------------------------------------------------------------
// PUBLIC API responses
// ------------------------------------------------------------------------------------------------

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type JWTClaims struct {
	Issuer    string   `json:"iss"`
	Subject   string   `json:"sub"`
	ExpiresAt int64    `json:"exp"`
	Roles     []string `json:"roles"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Organization struct {
	ID    string `json:"organizationId"`
	Name  string `json:"organizationName"`
	Email string `json:"email"`
}

type Permission struct {
	ID             string `json:"permissionId"`
	PermissionType string `json:"permissionType"`
	UserID         string `json:"userId"`
	ScopeID        string `json:"scopeId"`
	Scope          string `json:"scope"`
}

// APIResponse is a generic wrapper for public API responses.
type APIResponse[T any] struct {
	Data     T      `json:"data"`
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type WorkspaceResponse struct {
	ID            string `json:"workspaceId"`
	Name          string `json:"name"`
	DataResidency string `json:"dataResidency"`
	Notifications struct {
		Failure struct {
			Email   NotificationSetting `json:"email"`
			Webhook NotificationSetting `json:"webhook"`
		} `json:"failure"`
		Success struct {
			Email   NotificationSetting `json:"email"`
			Webhook NotificationSetting `json:"webhook"`
		} `json:"success"`
		ConnectionUpdate struct {
			Email   NotificationSetting `json:"email"`
			Webhook NotificationSetting `json:"webhook"`
		} `json:"connectionUpdate"`
		ConnectionUpdateActionRequired struct {
			Email   NotificationSetting `json:"email"`
			Webhook NotificationSetting `json:"webhook"`
		} `json:"connectionUpdateActionRequired"`
		SyncDisabled struct {
			Email   NotificationSetting `json:"email"`
			Webhook NotificationSetting `json:"webhook"`
		} `json:"syncDisabled"`
		SyncDisabledWarning struct {
			Email   NotificationSetting `json:"email"`
			Webhook NotificationSetting `json:"webhook"`
		} `json:"syncDisabledWarning"`
	} `json:"notifications"`
}

// NotificationSetting represents the enabled/disabled state of a notification channel.
type NotificationSetting struct {
	Enabled bool `json:"enabled"`
}

// ------------------------------------------------------------------------------------------------
// PRIVATE API responses
// ------------------------------------------------------------------------------------------------

// WorkspaceReadList represents a list of workspace reads.
type WorkspaceReadListResponse struct {
	Workspaces []WorkspaceReadResponse `json:"workspaces"`
}

// WorkspaceRead represents a detailed workspace configuration.
type WorkspaceReadResponse struct {
	WorkspaceId             string               `json:"workspaceId"`
	CustomerId              string               `json:"customerId"`
	OrganizationId          string               `json:"organizationId"`
	Name                    string               `json:"name"`
	Slug                    string               `json:"slug"`
	InitialSetupComplete    bool                 `json:"initialSetupComplete"`
	DisplaySetupWizard      bool                 `json:"displaySetupWizard"`
	AnonymousDataCollection bool                 `json:"anonymousDataCollection"`
	News                    bool                 `json:"news"`
	SecurityUpdates         bool                 `json:"securityUpdates"`
	Notifications           []interface{}        `json:"notifications"`
	NotificationSettings    NotificationSettings `json:"notificationSettings"`
	DefaultGeography        string               `json:"defaultGeography"`
	WebhookConfigs          []interface{}        `json:"webhookConfigs"`
	Tombstone               bool                 `json:"tombstone"`
}

// NotificationSettings represents the notification configuration for different events.
type NotificationSettings struct {
	SendOnSuccess                        NotificationTypes `json:"sendOnSuccess"`
	SendOnFailure                        NotificationTypes `json:"sendOnFailure"`
	SendOnSyncDisabled                   NotificationTypes `json:"sendOnSyncDisabled"`
	SendOnSyncDisabledWarning            NotificationTypes `json:"sendOnSyncDisabledWarning"`
	SendOnConnectionUpdate               NotificationTypes `json:"sendOnConnectionUpdate"`
	SendOnConnectionUpdateActionRequired NotificationTypes `json:"sendOnConnectionUpdateActionRequired"`
	SendOnBreakingChangeWarning          NotificationTypes `json:"sendOnBreakingChangeWarning"`
	SendOnBreakingChangeSyncsDisabled    NotificationTypes `json:"sendOnBreakingChangeSyncsDisabled"`
}

// NotificationTypes represents the types of notifications for a specific event.
type NotificationTypes struct {
	NotificationType []string `json:"notificationType"`
}

type WorkspaceUserAccessInfoReadListResponse struct {
	UsersWithAccess []WorkspaceUserAccessInfoReadResponse `json:"usersWithAccess"`
}

type WorkspaceUserAccessInfoReadResponse struct {
	UserID                 string          `json:"userId"`
	UserEmail              string          `json:"userEmail"`
	UserName               string          `json:"userName"`
	WorkspaceID            string          `json:"workspaceId"`
	WorkspacePermission    *PermissionRead `json:"workspacePermission,omitempty"`
	OrganizationPermission *PermissionRead `json:"organizationPermission,omitempty"`
}

type PermissionRead struct {
	PermissionID   string `json:"permissionId"`
	PermissionType string `json:"permissionType"`
	UserID         string `json:"userId"`
	WorkspaceID    string `json:"workspaceId,omitempty"`
	OrganizationID string `json:"organizationId,omitempty"`
}
