// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package netset

import (
	"log"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
)

const (
	TCPCode = 0
	UDPCode = 1
)

const (
	AllConnections = "All Connections"
	NoConnections  = "No Connections"
)

type Dimension string

const (
	protocol Dimension = "protocol"
	srcPort  Dimension = "srcPort"
	dstPort  Dimension = "dstPort"
)

const propertySeparator string = " "

type PortSet = interval.CanonicalSet
type ProtocolSet = interval.CanonicalSet

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

func (c *TCPUDPSet) IsEmpty() bool {
	return c.props.IsEmpty()
}

func (c *TCPUDPSet) Union(other *TCPUDPSet) *TCPUDPSet {
	return &TCPUDPSet{props: c.props.Union(other.props)}
}

func (c *TCPUDPSet) Size() int {
	return c.props.Size()
}

// SwapPorts returns a new NProduct object, built from the input NProduct object,
// with src ports and dst ports swapped
func (c *TCPUDPSet) SwapPorts() *TCPUDPSet {
	return &TCPUDPSet{props: c.props.Swap23()}
}

// Subtract
// ToDo: Subtract seems to ignore IsStateful (see https://github.com/np-guard/vpc-network-config-analyzer/issues/199):
//  1. is the delta connection stateful
//  2. props is identical but c stateful while other is not
//     the 2nd item can be computed here, with enhancement to relevant structure
//     the 1st can not since we do not know where exactly the statefulness came from
func (c *TCPUDPSet) Subtract(other *TCPUDPSet) *TCPUDPSet {
	return &TCPUDPSet{props: c.props.Subtract(other.props)}
}

// IsSubset returns true if c is subset of other
func (c *TCPUDPSet) IsSubset(other *TCPUDPSet) bool {
	return c.props.IsSubset(other.props)
}

// String returns a string representation of a TCPUDPSet object
func (c *TCPUDPSet) String() string {
	if c.IsEmpty() {
		return NoConnections
	} else if c.IsAll() {
		return AllConnections
	}
	// get cubes and cube str per each cube
	partitions := c.props.Partitions()
	resStrings := make([]string, len(partitions))
	for i, triple := range partitions {
		resStrings[i] = joinNonEmpty(
			getDimensionString(triple.S1, protocol),
			getDimensionString(triple.S2, srcPort),
			getDimensionString(triple.S3, dstPort),
		)
	}

	sort.Strings(resStrings)
	return strings.Join(resStrings, "; ")
}

func path(protocol *ProtocolSet, srcPort, dstPort *PortSet) *TCPUDPSet {
	return &TCPUDPSet{props: ds.CartesianRightTriple(protocol, srcPort, dstPort)}
}

// dimensionsList is the ordered list of dimensions in the TCPUDPSet object
// this should be the only place where the order is hard-coded
func entireDimension(dim Dimension) interval.Interval {
	switch dim {
	case protocol:
		return interval.New(TCPCode, UDPCode)
	case srcPort:
		return interval.New(netp.MinPort, netp.MaxPort)
	case dstPort:
		return interval.New(netp.MinPort, netp.MaxPort)
	}
	return interval.New(0, -1)
}

func All() *TCPUDPSet {
	return path(
		entireDimension(protocol).ToSet(),
		entireDimension(dstPort).ToSet(),
		entireDimension(srcPort).ToSet(),
	)
}

var all = All()

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

func NewProtocolSet(protocolString netp.ProtocolString) *ProtocolSet {
	p := protocolStringToCode(protocolString)
	return interval.New(p, p).ToSet()
}

func NewTCPorUDPSet(protocolString netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *TCPUDPSet {
	protocol := protocolStringToCode(protocolString)
	return path(
		interval.New(protocol, protocol).ToSet(),
		interval.New(srcMinP, srcMaxP).ToSet(),
		interval.New(dstMinP, dstMaxP).ToSet(),
	)
}

func EmptyTCPorUDPSet() *TCPUDPSet {
	return &TCPUDPSet{props: ds.NewRightTripleSet[*PortSet, *PortSet, *PortSet]()}
}

func protocolStringFromCode(protocolCode int64) netp.ProtocolString {
	switch protocolCode {
	case TCPCode:
		return netp.ProtocolStringTCP
	case UDPCode:
		return netp.ProtocolStringUDP
	}
	log.Panicf("impossible protocol code %v", protocolCode)
	return ""
}

func getDimensionString(dimValue *interval.CanonicalSet, dim Dimension) string {
	if dimValue.Equal(entireDimension(dim).ToSet()) {
		// avoid adding dimension str on full dimension values
		return ""
	}
	switch dim {
	case protocol:
		var pList []string
		for _, code := range []int64{TCPCode, UDPCode} {
			if dimValue.Contains(code) {
				pList = append(pList, string(protocolStringFromCode(code)))
			}
		}
		// sort by string values to avoid dependence on internal encoding
		sort.Strings(pList)
		return "protocol: " + strings.Join(pList, ",")
	case srcPort:
		return "src-ports: " + dimValue.String()
	case dstPort:
		return "dst-ports: " + dimValue.String()
	}
	return ""
}

func joinNonEmpty(inputList ...string) string {
	var res []string
	for _, propertyStr := range inputList {
		if propertyStr != "" {
			res = append(res, propertyStr)
		}
	}
	return strings.Join(res, propertySeparator)
}
