package dto

// PPPSecret represents a PPP secret (user) from MikroTik API
// Mapping dari: /ppp/secret/print
type PPPSecret struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name"`
	Password             string `json:"password,omitempty"`
	Profile              string `json:"profile,omitempty"`
	Service              string `json:"service,omitempty"`
	Disabled             bool   `json:"disabled,omitempty"`
	CallerID             string `json:"callerID,omitempty"`
	LocalAddress         string `json:"localAddress,omitempty"`
	RemoteAddress        string `json:"remoteAddress,omitempty"`
	Routes               string `json:"routes,omitempty"`
	Comment              string `json:"comment,omitempty"`
	LimitBytesIn         int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut        int64  `json:"limitBytesOut,omitempty"`
	LastLoggedOut        string `json:"lastLoggedOut,omitempty"`
	LastCallerID         string `json:"lastCallerID,omitempty"`
	LastDisconnectReason string `json:"lastDisconnectReason,omitempty"`
}

// PPPProfile represents a PPP profile from MikroTik API
// Mapping dari: /ppp/profile/print
type PPPProfile struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name"`
	LocalAddress       string `json:"localAddress,omitempty"`
	RemoteAddress      string `json:"remoteAddress,omitempty"`
	DNSServer          string `json:"dnsServer,omitempty"`
	SessionTimeout     string `json:"sessionTimeout,omitempty"`
	IdleTimeout        string `json:"idleTimeout,omitempty"`
	OnlyOne            bool   `json:"onlyOne,omitempty"`
	Comment            string `json:"comment,omitempty"`
	RateLimit          string `json:"rateLimit,omitempty"`
	ParentQueue        string `json:"parentQueue,omitempty"`
	QueueType          string `json:"queueType,omitempty"`
	UseCompression     bool   `json:"useCompression,omitempty"`
	UseEncryption      bool   `json:"useEncryption,omitempty"`
	UseMPLS            bool   `json:"useMPLS,omitempty"`
	UseUPnP            bool   `json:"useUPnP,omitempty"`
	Bridge             string `json:"bridge,omitempty"`
	AddressList        string `json:"addressList,omitempty"`
	InterfaceList      string `json:"interfaceList,omitempty"`
	OnUp               string `json:"onUp,omitempty"`
	OnDown             string `json:"onDown,omitempty"`
	ChangeTCPMSS       bool   `json:"changeTCPMSS,omitempty"`
	IncomingFilter     string `json:"incomingFilter,omitempty"`
	OutgoingFilter     string `json:"outgoingFilter,omitempty"`
	InsertQueueBefore  string `json:"insertQueueBefore,omitempty"`
	WinsServer         string `json:"winsServer,omitempty"`
	BridgeHorizon      string `json:"bridgeHorizon,omitempty"`
	BridgeLearning     bool   `json:"bridgeLearning,omitempty"`
	BridgePathCost     int    `json:"bridgePathCost,omitempty"`
	BridgePortPriority int    `json:"bridgePortPriority,omitempty"`
}

// PPPActive represents an active PPP session from MikroTik API
// Mapping dari: /ppp/active/print
type PPPActive struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Service       string `json:"service,omitempty"`
	CallerID      string `json:"callerID,omitempty"`
	Address       string `json:"address,omitempty"`
	Uptime        string `json:"uptime,omitempty"`
	SessionID     string `json:"sessionID,omitempty"`
	Encoding      string `json:"encoding,omitempty"`
	BytesIn       int64  `json:"bytesIn,omitempty"`
	BytesOut      int64  `json:"bytesOut,omitempty"`
	PacketsIn     int64  `json:"packetsIn,omitempty"`
	PacketsOut    int64  `json:"packetsOut,omitempty"`
	LimitBytesIn  int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut int64  `json:"limitBytesOut,omitempty"`
}

// PPPSecretRequest represents a request to add a PPP secret
type PPPSecretRequest struct {
	Name          string `json:"name" validate:"required,max=50"`
	Password      string `json:"password,omitempty" validate:"max=50"`
	Profile       string `json:"profile,omitempty" validate:"max=50"`
	Service       string `json:"service,omitempty"`
	CallerID      string `json:"callerID,omitempty" validate:"max=100"`
	LocalAddress  string `json:"localAddress,omitempty" validate:"max=50"`
	RemoteAddress string `json:"remoteAddress,omitempty" validate:"max=50"`
	Routes        string `json:"routes,omitempty"`
	Comment       string `json:"comment,omitempty" validate:"max=200"`
	LimitBytesIn  int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut int64  `json:"limitBytesOut,omitempty"`
}

// PPPSecretUpdateRequest represents a request to update a PPP secret
type PPPSecretUpdateRequest struct {
	Name          string `json:"name,omitempty" validate:"max=50"`
	Password      string `json:"password,omitempty" validate:"max=50"`
	Profile       string `json:"profile,omitempty" validate:"max=50"`
	Service       string `json:"service,omitempty"`
	Disabled      bool   `json:"disabled,omitempty"`
	CallerID      string `json:"callerID,omitempty" validate:"max=100"`
	LocalAddress  string `json:"localAddress,omitempty" validate:"max=50"`
	RemoteAddress string `json:"remoteAddress,omitempty" validate:"max=50"`
	Routes        string `json:"routes,omitempty"`
	Comment       string `json:"comment,omitempty" validate:"max=200"`
	LimitBytesIn  int64  `json:"limitBytesIn,omitempty"`
	LimitBytesOut int64  `json:"limitBytesOut,omitempty"`
}

// PPPProfileRequest represents a request to add/update a PPP profile
type PPPProfileRequest struct {
	Name               string `json:"name" validate:"required,max=50"`
	LocalAddress       string `json:"localAddress,omitempty" validate:"max=50"`
	RemoteAddress      string `json:"remoteAddress,omitempty" validate:"max=50"`
	DNSServer          string `json:"dnsServer,omitempty" validate:"max=100"`
	SessionTimeout     string `json:"sessionTimeout,omitempty" validate:"max=20"`
	IdleTimeout        string `json:"idleTimeout,omitempty" validate:"max=20"`
	OnlyOne            bool   `json:"onlyOne,omitempty"`
	Comment            string `json:"comment,omitempty" validate:"max=200"`
	RateLimit          string `json:"rateLimit,omitempty" validate:"max=100"`
	ParentQueue        string `json:"parentQueue,omitempty" validate:"max=50"`
	QueueType          string `json:"queueType,omitempty" validate:"max=50"`
	UseCompression     bool   `json:"useCompression,omitempty"`
	UseEncryption      bool   `json:"useEncryption,omitempty"`
	UseMPLS            bool   `json:"useMPLS,omitempty"`
	UseUPnP            bool   `json:"useUPnP,omitempty"`
	Bridge             string `json:"bridge,omitempty" validate:"max=50"`
	AddressList        string `json:"addressList,omitempty" validate:"max=50"`
	InterfaceList      string `json:"interfaceList,omitempty" validate:"max=50"`
	OnUp               string `json:"onUp,omitempty"`
	OnDown             string `json:"onDown,omitempty"`
	ChangeTCPMSS       bool   `json:"changeTCPMSS,omitempty"`
	IncomingFilter     string `json:"incomingFilter,omitempty" validate:"max=50"`
	OutgoingFilter     string `json:"outgoingFilter,omitempty" validate:"max=50"`
	InsertQueueBefore  string `json:"insertQueueBefore,omitempty" validate:"max=50"`
	WinsServer         string `json:"winsServer,omitempty" validate:"max=100"`
}
