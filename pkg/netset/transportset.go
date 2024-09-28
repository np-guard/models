/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset

import (
	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/netp"
)

// type connection.Set is an alias for netset.TransportSet

// TransportSet captures connection-sets for protocols from {TCP, UDP, ICMP}
type TransportSet struct {
	set *ds.Disjoint[*TCPUDPSet, *ICMPSet]
}

func NewTCPorUDPTransport(protocol netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *TransportSet {
	return &TransportSet{ds.NewDisjoint(
		NewTCPorUDPSet(protocol, srcMinP, srcMaxP, dstMinP, dstMaxP),
		EmptyICMPSet(),
	)}
}

func NewICMPTransport(tc netp.ICMP) *TransportSet {
	return &TransportSet{ds.NewDisjoint(
		EmptyTCPorUDPSet(),
		NewICMPSet(tc),
	)}
}

func AllOrNothingTransport(allTcpubp, allIcmp bool) *TransportSet {
	var tcpudp *TCPUDPSet
	var icmp *ICMPSet
	if allTcpubp {
		tcpudp = AllTCPUDPSet()
	} else {
		tcpudp = EmptyTCPorUDPSet()
	}
	if allIcmp {
		icmp = AllICMPSet()
	} else {
		icmp = EmptyICMPSet()
	}
	return &TransportSet{ds.NewDisjoint(tcpudp, icmp)}
}

func AllTransportSet() *TransportSet {
	return AllOrNothingTransport(true, true)
}

func (t *TransportSet) SwapPorts() *TransportSet {
	return &TransportSet{ds.NewDisjoint(t.TCPUDPSet().SwapPorts(), t.ICMPSet())}
}

func (t *TransportSet) TCPUDPSet() *TCPUDPSet {
	return t.set.Left()
}

func (t *TransportSet) ICMPSet() *ICMPSet {
	return t.set.Right()
}

func (t *TransportSet) Equal(other *TransportSet) bool {
	return t.set.Equal(other.set)
}

func (t *TransportSet) Copy() *TransportSet {
	return &TransportSet{t.set.Copy()}
}

func (t *TransportSet) Hash() int {
	return t.set.Hash()
}

func (t *TransportSet) IsEmpty() bool {
	return t.set.IsEmpty()
}

func (t *TransportSet) IsAll() bool {
	return t.Equal(AllTransportSet())
}

func (t *TransportSet) Size() int {
	return t.set.Size()
}

// IsSubset returns true if c is subset of other
func (t *TransportSet) IsSubset(other *TransportSet) bool {
	return t.set.IsSubset(other.set)
}

func (t *TransportSet) Union(other *TransportSet) *TransportSet {
	return &TransportSet{t.set.Union(other.set)}
}

func (t *TransportSet) Intersect(other *TransportSet) *TransportSet {
	return &TransportSet{t.set.Intersect(other.set)}
}

func (t *TransportSet) Subtract(other *TransportSet) *TransportSet {
	return &TransportSet{t.set.Subtract(other.set)}
}

func (t *TransportSet) String() string {
	return ""
}
