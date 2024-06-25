/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package connection

import (
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
)

type Set = netset.TransportSet

func NewTCPorUDP(protocol netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *Set {
	return netset.NewTCPorUDPTransport(protocol, srcMinP, srcMaxP, dstMinP, dstMaxP)
}

func AllICMP() *Set {
	return netset.AllOrNothingTransport(false, true)
}

func NewTCPSet() *Set {
	return NewTCPorUDP(netp.ProtocolStringTCP, netp.MinPort, netp.MaxPort, netp.MinPort, netp.MaxPort)
}

func ICMPConnection(icmpType, icmpCode *int64) (*Set, error) {
	icmp, err := netp.ICMPFromTypeAndCode64(icmpType, icmpCode)
	if err != nil {
		return nil, err
	}
	return netset.NewICMPTransport(icmp), nil
}

func All() *Set {
	return netset.AllTransportSet()
}

func None() *Set {
	return netset.AllOrNothingTransport(false, false)
}
