/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package connection

import (
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
)

// Set captures a set of connections for protocols TCP/UPD/ICMP with their properties (ports/icmp type&code)
type Set = netset.TransportSet

// NewTCPorUDP returns a set of connections containing the specified protocol (TCP/UDP) and ports
func NewTCPorUDP(protocol netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *Set {
	return netset.NewTCPorUDPTransport(protocol, srcMinP, srcMaxP, dstMinP, dstMaxP)
}

// NewTCP returns a set of TCP connections containing the specified ports
func NewTCP(srcMinP, srcMaxP, dstMinP, dstMaxP int64) *Set {
	return NewTCPorUDP(netp.ProtocolStringTCP, srcMinP, srcMaxP, dstMinP, dstMaxP)
}

// NewUDP returns a set of UDP connections containing the specified ports
func NewUDP(srcMinP, srcMaxP, dstMinP, dstMaxP int64) *Set {
	return NewTCPorUDP(netp.ProtocolStringUDP, srcMinP, srcMaxP, dstMinP, dstMaxP)
}

// AllTCPorUDP returns a set of connections containing the specified protocol (TCP/UDP) with all possible ports
func AllTCPorUDP(protocol netp.ProtocolString) *Set {
	return NewTCPorUDP(protocol, netp.MinPort, netp.MaxPort, netp.MinPort, netp.MaxPort)
}

// AllICMP returns a set of connections containing the ICMP protocol with all its possible types,codes
func AllICMP() *Set {
	return netset.AllOrNothingTransport(false, true)
}

// NewTCPSet returns a set of connections containing the TCP protocol with all its possible ports
func NewTCPSet() *Set {
	return AllTCPorUDP(netp.ProtocolStringTCP)
}

// NewUDPSet returns a set of connections containing the UDP protocol with all its possible ports
func NewUDPSet() *Set {
	return AllTCPorUDP(netp.ProtocolStringUDP)
}

// ICMPConnection returns a set of connections containing the ICMP protocol with specified type,code values
func ICMPConnection(icmpType, icmpCode *int64) (*Set, error) {
	icmp, err := netp.ICMPFromTypeAndCode64(icmpType, icmpCode)
	if err != nil {
		return nil, err
	}
	return netset.NewICMPTransport(icmp), nil
}

// All returns a set of all protocols (TCP,UPD,ICMP) in the set (with all possible properties values)
func All() *Set {
	return netset.AllTransportSet()
}

// None returns an empty set of protocols connections
func None() *Set {
	return netset.AllOrNothingTransport(false, false)
}
