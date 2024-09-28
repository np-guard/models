/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset

import (
	"log"
	"slices"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
)

// thid file defines type TCPUDPSet as TripleSet[*ProtocolSet, *PortSet, *PortSet]

// encoding TCP/UDP protocols as integers for TCPUDPSet
const (
	TCPCode = 0
	UDPCode = 1
)

// TODO: currently assuming input values are always within valid ranges.
// should add validation / error handling for this?

type ProtocolSet = interval.CanonicalSet // valid range: [0,1] (see TCPCode , UDPCode)
type PortSet = interval.CanonicalSet     // valid range: [1,65535]  (see netp.MinPort , netp.MaxPort)

func AllPorts() *PortSet {
	return netp.AllPorts().ToSet()
}

func AllTCPUDPProtocolSet() *ProtocolSet {
	return interval.New(TCPCode, UDPCode).ToSet()
}

// TCPUDPSet captures sets of protocols (within TCP,UDP only) and ports (source and destinaion)
type TCPUDPSet struct {
	// S1: protocols, S2: src ports, S3: dst ports
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
		AllTCPUDPProtocolSet(),
		AllPorts(),
		AllPorts(),
	)
}

func NewAllTCPOnlySet() *TCPUDPSet {
	return pathLeft(
		interval.New(TCPCode, TCPCode).ToSet(),
		AllPorts(),
		AllPorts(),
	)
}

func NewAllUDPOnlySet() *TCPUDPSet {
	return pathLeft(
		interval.New(UDPCode, UDPCode).ToSet(),
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

func protocolCodeToString(pSet *ProtocolSet) string {
	switch {
	case pSet.Equal(AllTCPUDPProtocolSet()):
		return string(netp.ProtocolStringTCP) + comma + string(netp.ProtocolStringUDP)
	case pSet.Contains(UDPCode):
		return string(netp.ProtocolStringUDP)
	case pSet.Contains(TCPCode):
		return string(netp.ProtocolStringTCP)
	}
	return ""
}

func getTCPUDPCubeStr(cube ds.Triple[*ProtocolSet, *PortSet, *PortSet]) string {
	var ports []string
	if !cube.S2.Equal(AllPorts()) {
		ports = append(ports, "src-ports: "+cube.S2.String())
	}
	if !cube.S3.Equal(AllPorts()) {
		ports = append(ports, "dst-ports: "+cube.S3.String())
	}
	protocolsStr := protocolCodeToString(cube.S1)
	allComponentsStrList := slices.Concat([]string{protocolsStr}, ports)
	return strings.Join(allComponentsStrList, " ")
}

func (c *TCPUDPSet) String() string {
	cubes := c.Partitions()
	var resStrings = make([]string, len(cubes))
	for i, cube := range cubes {
		resStrings[i] = getTCPUDPCubeStr(cube)
	}
	sort.Strings(resStrings)
	return strings.Join(resStrings, " | ")
}
