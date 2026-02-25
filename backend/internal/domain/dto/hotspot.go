package dto

// HotspotUser represents a hotspot user from MikroTik API
// Mapping dari: /ip/hotspot/user/print
type HotspotUser struct {
	ID              string `json:"id,omitempty"`
	Server          string `json:"server,omitempty"`
	Name            string `json:"name"`
	Password        string `json:"password,omitempty"`
	Profile         string `json:"profile,omitempty"`
	MACAddress      string `json:"macAddress,omitempty"`
	IPAddress       string `json:"ipAddress,omitempty"`
	Uptime          string `json:"uptime,omitempty"`
	BytesIn         int64  `json:"bytesIn,omitempty"`
	BytesOut        int64  `json:"bytesOut,omitempty"`
	PacketsIn       int64  `json:"packetsIn,omitempty"`
	PacketsOut      int64  `json:"packetsOut,omitempty"`
	LimitUptime     string `json:"limitUptime,omitempty"`
	LimitBytesIn    int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut   int64  `json:"limitBytesOut,omitempty"`
	LimitBytesTotal int64  `json:"limitBytesTotal,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`
	Email           string `json:"email,omitempty"`
	Routes          string `json:"routes,omitempty"`
}

// HotspotActive represents an active hotspot session
// Mapping dari: /ip/hotspot/active/print
type HotspotActive struct {
	ID               string `json:"id,omitempty"`
	Server           string `json:"server,omitempty"`
	User             string `json:"user,omitempty"`
	Address          string `json:"address,omitempty"`
	MACAddress       string `json:"macAddress,omitempty"`
	LoginBy          string `json:"loginBy,omitempty"`
	Uptime           string `json:"uptime,omitempty"`
	SessionTimeLeft  string `json:"sessionTimeLeft,omitempty"`
	IdleTime         string `json:"idleTime,omitempty"`
	IdleTimeout      string `json:"idleTimeout,omitempty"`
	KeepaliveTimeout string `json:"keepaliveTimeout,omitempty"`
	BytesIn          int64  `json:"bytesIn,omitempty"`
	BytesOut         int64  `json:"bytesOut,omitempty"`
	PacketsIn        int64  `json:"packetsIn,omitempty"`
	PacketsOut       int64  `json:"packetsOut,omitempty"`
	Radius           bool   `json:"radius,omitempty"`
}

// HotspotHost represents a hotspot host
// Mapping dari: /ip/hotspot/host/print
type HotspotHost struct {
	ID           string `json:"id,omitempty"`
	MACAddress   string `json:"macAddress,omitempty"`
	Address      string `json:"address,omitempty"`
	ToAddress    string `json:"toAddress,omitempty"`
	Server       string `json:"server,omitempty"`
	Authorized   bool   `json:"authorized,omitempty"`
	Bypassed     bool   `json:"bypassed,omitempty"`
	Blocked      bool   `json:"blocked,omitempty"`
	FoundBy      string `json:"foundBy,omitempty"`
	HostDeadTime string `json:"hostDeadTime,omitempty"`
	Comment      string `json:"comment,omitempty"`
}

// UserProfile represents a hotspot user profile
// Mapping dari: /ip/hotspot/user/profile/print
type UserProfile struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name"`
	AddressPool          string `json:"addressPool,omitempty"`
	SharedUsers          int    `json:"sharedUsers,omitempty"`
	RateLimit            string `json:"rateLimit,omitempty"`
	ParentQueue          string `json:"parentQueue,omitempty"`
	StatusAutorefresh    string `json:"statusAutorefresh,omitempty"`
	OnLogin              string `json:"onLogin,omitempty"`
	OnLogout             string `json:"onLogout,omitempty"`
	OnUp                 string `json:"onUp,omitempty"`
	OnDown               string `json:"onDown,omitempty"`
	TransparentProxy     bool   `json:"transparentProxy,omitempty"`
	OpenStatusPage       string `json:"openStatusPage,omitempty"`
	Advertise            bool   `json:"advertise,omitempty"`
	AdvertiseInterval    string `json:"advertiseInterval,omitempty"`
	AdvertiseTimeout     string `json:"advertiseTimeout,omitempty"`
	AdvertiseURL         string `json:"advertiseURL,omitempty"`

	// Mikhmon-specific fields (parsed from on-login script output)
	ExpireMode   string  `json:"expireMode,omitempty"`
	Validity     string  `json:"validity,omitempty"`
	Price        float64 `json:"price,omitempty"`
	SellingPrice float64 `json:"sellingPrice,omitempty"`
	LockUser     string  `json:"lockUser,omitempty"`
	LockServer   string  `json:"lockServer,omitempty"`
}

// AddUserRequest represents a request to add a hotspot user
type AddUserRequest struct {
	RouterID        uint   `json:"routerId" validate:"required"`
	Server          string `json:"server,omitempty"`
	Name            string `json:"name" validate:"required,max=50"`
	Password        string `json:"password,omitempty" validate:"max=50"`
	Profile         string `json:"profile" validate:"required,max=50"`
	MACAddress      string `json:"macAddress,omitempty" validate:"omitempty,mac"`
	TimeLimit       string `json:"timeLimit,omitempty" validate:"max=20"`
	DataLimit       string `json:"dataLimit,omitempty" validate:"max=20"`
	Comment         string `json:"comment,omitempty" validate:"max=100"`
}

// CreateUserRequest represents a request to create a hotspot user (internal use)
type CreateUserRequest struct {
	Server          string `json:"server,omitempty"`
	Name            string `json:"name" validate:"required,max=50"`
	Password        string `json:"password,omitempty" validate:"max=50"`
	Profile         string `json:"profile" validate:"required,max=50"`
	MACAddress      string `json:"macAddress,omitempty"`
	LimitUptime     string `json:"limitUptime,omitempty"`
	LimitBytesTotal int64  `json:"limitBytesTotal,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`
}

// UpdateUserRequest represents a request to update a hotspot user
type UpdateUserRequest struct {
	Server          string `json:"server,omitempty" validate:"max=50"`
	Name            string `json:"name,omitempty" validate:"max=50"`
	Password        string `json:"password,omitempty" validate:"max=50"`
	Profile         string `json:"profile,omitempty" validate:"max=50"`
	MACAddress      string `json:"macAddress,omitempty" validate:"omitempty,mac"`
	TimeLimit       string `json:"timeLimit,omitempty" validate:"max=20"`
	DataLimit       string `json:"dataLimit,omitempty" validate:"max=20"`
	Comment         string `json:"comment,omitempty" validate:"max=100"`
	Disabled        bool   `json:"disabled,omitempty"`
}

// RemoveUserRequest represents a request to remove users by comment (for vouchers)
type RemoveUserRequest struct {
	IDs     []string `json:"ids,omitempty"`      // Remove by IDs
	Comment string   `json:"comment,omitempty"`  // Remove by comment (for vouchers)
}

// GetUsersRequest represents a request to get users
type GetUsersRequest struct {
	Profile string `json:"profile,omitempty"`
	Comment string `json:"comment,omitempty"` // For printing vouchers
}

// ProfileRequest represents a request to create/update profile
type ProfileRequest struct {
	Name         string  `json:"name" validate:"required,max=50"`
	AddressPool  string  `json:"addressPool,omitempty" validate:"max=50"`
	SharedUsers  int     `json:"sharedUsers,omitempty" validate:"min=1,max=999"`
	RateLimit    string  `json:"rateLimit,omitempty" validate:"max=100"`
	ParentQueue  string  `json:"parentQueue,omitempty" validate:"max=50"`
	ExpireMode   string  `json:"expireMode,omitempty" validate:"omitempty,oneof=0 rem ntf remc ntfc"`
	Validity     string  `json:"validity,omitempty" validate:"max=20"`
	Price        float64 `json:"price,omitempty" validate:"min=0"`
	SellingPrice float64 `json:"sellingPrice,omitempty" validate:"min=0"`
	LockUser     string  `json:"lockUser,omitempty" validate:"omitempty,oneof=Disable Enable"`
	LockServer   string  `json:"lockServer,omitempty" validate:"omitempty,oneof=Disable Enable"`
}

// ProfileUpdateRequest represents a request to update profile
type ProfileUpdateRequest struct {
	Name         string  `json:"name,omitempty" validate:"max=50"`
	AddressPool  string  `json:"addressPool,omitempty" validate:"max=50"`
	SharedUsers  int     `json:"sharedUsers,omitempty" validate:"min=1,max=999"`
	RateLimit    string  `json:"rateLimit,omitempty" validate:"max=100"`
	ParentQueue  string  `json:"parentQueue,omitempty" validate:"max=50"`
	ExpireMode   string  `json:"expireMode,omitempty" validate:"omitempty,oneof=0 rem ntf remc ntfc"`
	Validity     string  `json:"validity,omitempty" validate:"max=20"`
	Price        float64 `json:"price,omitempty" validate:"min=0"`
	SellingPrice float64 `json:"sellingPrice,omitempty" validate:"min=0"`
	LockUser     string  `json:"lockUser,omitempty" validate:"omitempty,oneof=Disable Enable"`
	LockServer   string  `json:"lockServer,omitempty" validate:"omitempty,oneof=Disable Enable"`
}

// UserFilter represents filter options for getting users
type UserFilter struct {
	Profile string
	Comment string
}
