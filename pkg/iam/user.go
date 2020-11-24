package iam

type (
	// User encapsulates information about each user.
	User struct {
		Actions            *UserActions      `json:"actions,omitempty"`
		AuthGrants         []*AuthGrant      `json:"authGrants,omitempty"`
		Notifications      *UserNotification `json:"notifications,omitempty"`
		IdentityID         string            `json:"uiIdentityId"`
		FirstName          string            `json:"firstName"`
		LastName           string            `json:"lastName"`
		UserName           string            `json:"uiUserName,omitempty"`
		Email              string            `json:"email"`
		IsLocked           bool              `json:"isLocked"`
		Phone              string            `json:"phone,omitempty"`
		TimeZone           string            `json:"timeZone,omitempty"`
		LastLoginDate      string            `json:"lastLoginDate,omitempty"`
		ContactType        string            `json:"contactType"`
		PreferredLanguage  string            `json:"preferredLanguage,omitempty"`
		SessionTimeOut     *int64            `json:"sessionTimeOut,omitempty"`
		PasswordExpiryDate string            `json:"passwordExpiryDate,omitempty"`
		SecondaryEmail     string            `json:"secondaryEmail,omitempty"`
		MobilePhone        string            `json:"mobilePhone,omitempty"`
		Street             string            `json:"street,omitempty"`
		City               string            `json:"city,omitempty"`
		State              string            `json:"state,omitempty"`
		ZipCode            string            `json:"zipCode,omitempty"`
		Country            string            `json:"country"`
		JobTitle           string            `json:"jobTitle"`
		TFAEnabled         bool              `json:"tfaEnabled"`
		TFAConfigured      bool              `json:"tfaConfigured"`
		EmailUpdatePending bool              `json:"emailUpdatePending"`
	}

	// UserActions encapsulates permissions available to the user for this group.
	UserActions struct {
		APIClient        bool `json:"apiClient"`
		Delete           bool `json:"delete"`
		Edit             bool `json:"edit"`
		IsCloneable      bool `json:"isCloneable"`
		ResetPassword    bool `json:"resetPassword"`
		ThirdPartyAccess bool `json:"thirdPartyAccess"`
	}

	// AuthGrant is user’s role assignments, per group.
	AuthGrant struct {
		GroupID         int64  `json:"groupId"`
		GroupName       string `json:"groupName"`
		IsBlocked       bool   `json:"isBlocked"`
		RoleDescription string `json:"roleDescription"`
		RoleID          int64  `json:"roleId"`
		RoleName        string `json:"roleName"`
	}

	// UserNotification types of notification emails the user receives.
	UserNotification struct {
		EnableEmail bool                     `json:"enableEmailNotifications"`
		Options     *UserNotificationOptions `json:"options,omitempty"`
	}

	// UserNotificationOptions types of notification emails the user receives.
	UserNotificationOptions struct {
		NewUser        bool     `json:"newUserNotification"`
		PasswordExpiry bool     `json:"passwordExpiry"`
		Proactive      []string `json:"proactive,omitempty"`
		Upgrade        []string `json:"upgrade,omitempty"`
	}
)
