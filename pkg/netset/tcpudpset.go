/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset

import (
	"log"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
)

const (
	TCPCode = 0
	UDPCode = 1
)

type PortSet = interval.CanonicalSet
type ProtocolSet = interval.CanonicalSet

func AllPorts() *PortSet {
	return netp.AllPorts().ToSet()
}

type TCPUDPSet struct {
	props ds.TripleSet[*ProtocolSet, *PortSet, *PortSet]
}

func (c *TCPUDPSet) Equal(other *TCPUDPSet) bool {
	return c.props.Equal(other.props)
}

func (c *TCPUDPSet) Hash() int {
	return c.props.Hash()
}

func (c *TCPUDPSet) Copy() *TCPUDPSet {
	return &TCPUDPSet{props: c.props.Copy()}
}

func (c *TCPUDPSet) Intersect(other *TCPUDPSet) *TCPUDPSet {
	return &TCPUDPSet{props: c.props.Intersect(other.props)}
}

func (c *TCPUDPSet) Partitions() []ds.Triple[*ProtocolSet, *PortSet, *PortSet] {
	return c.props.Partitions()
}

func (c *TCPUDPSet) IsEmpty() bool {
	return c.props.IsEmpty()
}

func (c *TCPUDPSet) Union(other *TCPUDPSet) *TCPUDPSet {
	return &TCPUDPSet{props: c.props.Union(other.props)}
}

func (c *TCPUDPSet) Size() int {
	return c.props.Size()
}

// SwapPorts returns a new TCPUDPSet object, built from the input TCPUDPSet object,
// with src ports and dst ports swapped
func (c *TCPUDPSet) SwapPorts() *TCPUDPSet {
	return &TCPUDPSet{props: ds.MapTripleSet(c.props, ds.Triple[*ProtocolSet, *PortSet, *PortSet].Swap23)}
}

// Subtract returns the subtraction of the other from c
func (c *TCPUDPSet) Subtract(other *TCPUDPSet) *TCPUDPSet {
	return &TCPUDPSet{props: c.props.Subtract(other.props)}
}

// IsSubset returns true if c is subset of other
func (c *TCPUDPSet) IsSubset(other *TCPUDPSet) bool {
	return c.props.IsSubset(other.props)
}

// pathLeft creates a new TCPUDPSet, implemented using LeftTriple.
func pathLeft(protocol *ProtocolSet, srcPort, dstPort *PortSet) *TCPUDPSet {
	return &TCPUDPSet{props: ds.CartesianLeftTriple(protocol, srcPort, dstPort)}
}

func EmptyTCPorUDPSet() *TCPUDPSet {
	return &TCPUDPSet{props: ds.NewLeftTripleSet[*ProtocolSet, *PortSet, *PortSet]()}
}

func AllTCPUDPSet() *TCPUDPSet {
	return pathLeft(
		interval.New(TCPCode, UDPCode).ToSet(),
		AllPorts(),
		AllPorts(),
	)
}

func NewTCPorUDPSet(protocolString netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *TCPUDPSet {
	protocol := protocolStringToCode(protocolString)
	return pathLeft(
		interval.New(protocol, protocol).ToSet(),
		interval.New(srcMinP, srcMaxP).ToSet(),
		interval.New(dstMinP, dstMaxP).ToSet(),
	)
}

var all = AllTCPUDPSet()

func (c *TCPUDPSet) IsAll() bool {
	return c.Equal(all)
}

func protocolStringToCode(protocol netp.ProtocolString) int64 {
	switch protocol {
	case netp.ProtocolStringTCP:
		return TCPCode
	case netp.ProtocolStringUDP:
		return UDPCode
	}
	log.Panicf("Impossible protocol code %v", protocol)
	return 0
}
