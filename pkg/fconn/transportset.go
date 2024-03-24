// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package fconn

import (
	"log"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
)

const (
	TCPCode           = 0
	UDPCode           = 1
	minProtocol int64 = 0
	maxProtocol int64 = 2
	MinPort           = 1
	MaxPort           = netp.MaxPort
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

type Set struct {
	props *ds.TripleSet[*ProtocolSet, *PortSet, *PortSet]
}

func (c *Set) Equal(other *Set) bool {
	return c.props.Equal(other.props)
}

func (c *Set) Hash() int {
	return c.props.Hash()
}

func (c *Set) Copy() *Set {
	return &Set{
		props: c.props.Copy(),
	}
}

func (c *Set) Intersect(other *Set) *Set {
	return &Set{props: c.props.Intersect(other.props)}
}

func (c *Set) IsEmpty() bool {
	return c.props.IsEmpty()
}

func (c *Set) Union(other *Set) *Set {
	if other.IsEmpty() {
		return c.Copy()
	}
	if c.IsEmpty() {
		return other.Copy()
	}
	return &Set{
		props: c.props.Union(other.props),
	}
}

// Subtract
// ToDo: Subtract seems to ignore IsStateful (see https://github.com/np-guard/vpc-network-config-analyzer/issues/199):
//  1. is the delta connection stateful
//  2. props is identical but c stateful while other is not
//     the 2nd item can be computed here, with enhancement to relevant structure
//     the 1st can not since we do not know where exactly the statefulness came from
func (c *Set) Subtract(other *Set) *Set {
	if c.IsEmpty() {
		return None()
	}
	if other.IsEmpty() {
		return c.Copy()
	}
	return &Set{props: c.props.Subtract(other.props)}
}

// ContainedIn returns true if c is subset of other
func (c *Set) ContainedIn(other *Set) bool {
	return c.props.ContainedIn(other.props)
}

func None() *Set {
	return &Set{props: ds.NewTripleSet[*PortSet, *PortSet, *PortSet]()}
}

func Path(protocol *ProtocolSet, srcPort, dstPort *PortSet) *Set {
	return &Set{props: ds.Path(protocol, srcPort, dstPort)}
}

// dimensionsList is the ordered list of dimensions in the Set object
// this should be the only place where the order is hard-coded
func entireDimension(dim Dimension) interval.Interval {
	switch dim {
	case protocol:
		return interval.New(minProtocol, maxProtocol)
	case srcPort:
		return interval.New(MinPort, MaxPort)
	case dstPort:
		return interval.New(MinPort, MaxPort)
	}
	return interval.New(0, -1)
}

func All() *Set {
	return Path(
		entireDimension(protocol).ToSet(),
		entireDimension(dstPort).ToSet(),
		entireDimension(srcPort).ToSet(),
	)
}

var all = All()

func (c *Set) IsAll() bool {
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

func TCPorUDPConnection(protocolString netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *Set {
	protocol := protocolStringToCode(protocolString)
	return Path(
		interval.New(protocol, protocol).ToSet(),
		interval.New(srcMinP, srcMaxP).ToSet(),
		interval.New(dstMinP, dstMaxP).ToSet(),
	)
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
		pList := []string{}
		for code := minProtocol; code <= maxProtocol; code++ {
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
	res := []string{}
	for _, propertyStr := range inputList {
		if propertyStr != "" {
			res = append(res, propertyStr)
		}
	}
	return strings.Join(res, propertySeparator)
}

// String returns a string representation of a Set object
func (c *Set) String() string {
	if c.IsEmpty() {
		return NoConnections
	} else if c.IsAll() {
		return AllConnections
	}
	// get cubes and cube str per each cube
	resStrings := []string{}
	for _, triple := range c.props.Triples() {
		resStrings = append(resStrings, joinNonEmpty(
			getDimensionString(triple.S1, protocol),
			getDimensionString(triple.S2, srcPort),
			getDimensionString(triple.S3, dstPort),
		))
	}

	sort.Strings(resStrings)
	return strings.Join(resStrings, "; ")
}
