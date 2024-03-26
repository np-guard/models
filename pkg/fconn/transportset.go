// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package fconn

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netp"
)

type TransportSet = ds.Disjoint[*TCPUDPSet, *ICMPSet]

func NewTCPorUDPTransport(protocol netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *TransportSet {
	return ds.NewDisjoint(
		NewTCPorUDPSet(protocol, srcMinP, srcMaxP, dstMinP, dstMaxP),
		EmptyICMPSet(),
	)
}

func NewICMPTransport(tc netp.ICMP) *TransportSet {
	return ds.NewDisjoint(
		EmptyTCPorUDPSet(),
		NewICMPSet(tc),
	)
}
