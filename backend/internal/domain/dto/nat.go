package dto

// NATRule represents a firewall NAT rule
// Mapping from: /ip/firewall/nat/print
type NATRule struct {
	ID              string `json:"id,omitempty"`
	Chain           string `json:"chain,omitempty"`
	Action          string `json:"action,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	SrcAddress      string `json:"srcAddress,omitempty"`
	DstAddress      string `json:"dstAddress,omitempty"`
	SrcPort         string `json:"srcPort,omitempty"`
	DstPort         string `json:"dstPort,omitempty"`
	InInterface     string `json:"inInterface,omitempty"`
	OutInterface    string `json:"outInterface,omitempty"`
	ToAddresses     string `json:"toAddresses,omitempty"`
	ToPorts         string `json:"toPorts,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`
	Comment         string `json:"comment,omitempty"`
	Dynamic         bool   `json:"dynamic,omitempty"`
	Invalid         bool   `json:"invalid,omitempty"`
	Bytes           int64  `json:"bytes,omitempty"`
	Packets         int64  `json:"packets,omitempty"`
	ConnectionBytes int64  `json:"connectionBytes,omitempty"`
}
