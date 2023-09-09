package models

import (
	"fmt"
	"strings"
)

// A PortMapping represents SORACOM Napter port mapping
type PortMapping struct {
	Duration    int    `json:"duration"`    // duration in seconds
	Endpoint    string `json:"endpoint"`    // SORACOM Napter endpoint
	Hostname    string `json:"hostname"`    // SORACOM Napter hostname
	IPAddress   string `json:"ipAddress"`   // SORACOM Napter IP address
	Port        int    `json:"port"`        // SORACOM Napter port number
	TLSRequired bool   `json:"tlsRequired"` // is TLS required
	Destination struct {
		ID   string `json:"simId"` // target SIM ID
		Port int    `json:"port"`  // target port
	} `json:"destination"`
	Source struct {
		IPRanges []string `json:"ipRanges"` // permitted source CIDRs
	} `json:"source"`
}

func (pm PortMapping) String() string {
	return fmt.Sprintf("- Endpoint: %v:%v\n"+
		"- Destination: %v:%v\n"+
		"- Duration: %v hours\n"+
		"- Source: %v\n"+
		"- TLS required: %v",
		pm.Hostname, pm.Port, pm.Destination.ID, pm.Destination.Port, float32(pm.Duration)/60/60, strings.Join(pm.Source.IPRanges, ","), pm.TLSRequired)
}
