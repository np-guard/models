// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package netset

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

type ConnectionSet struct {
	props ds.TripleSet[*IPBlock, *IPBlock, *TransportSet]
}

func (c *ConnectionSet) Equal(other *ConnectionSet) bool {
	return c.props.Equal(other.props)
}

func (c *ConnectionSet) Copy() *ConnectionSet {
	return &ConnectionSet{
		props: c.props.Copy(),
	}
}

func (c *ConnectionSet) Intersect(other *ConnectionSet) *ConnectionSet {
	return &ConnectionSet{props: c.props.Intersect(other.props)}
}

func (c *ConnectionSet) IsEmpty() bool {
	return c.props.IsEmpty()
}

func (c *ConnectionSet) Union(other *ConnectionSet) *ConnectionSet {
	if other.IsEmpty() {
		return c.Copy()
	}
	if c.IsEmpty() {
		return other.Copy()
	}
	return &ConnectionSet{
		props: c.props.Union(other.props),
	}
}

// Subtract
// ToDo: Subtract seems to ignore IsStateful (see https://github.com/np-guard/vpc-network-config-analyzer/issues/199):
//  1. is the delta connection stateful
//  2. props is identical but c stateful while other is not
//     the 2nd item can be computed here, with enhancement to relevant structure
//     the 1st can not since we do not know where exactly the statefulness came from
func (c *ConnectionSet) Subtract(other *ConnectionSet) *ConnectionSet {
	if c.IsEmpty() {
		return &ConnectionSet{props: ds.NewRightTripleSet[*IPBlock, *IPBlock, *ds.Disjoint[*TCPUDPSet, *ICMPSet]]()}
	}
	if other.IsEmpty() {
		return c.Copy()
	}
	return &ConnectionSet{props: c.props.Subtract(other.props)}
}

// IsSubset returns true if c is subset of other
func (c *ConnectionSet) IsSubset(other *ConnectionSet) bool {
	return c.props.IsSubset(other.props)
}
