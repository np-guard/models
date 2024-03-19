// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package connection

import (
	"log"
	"slices"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/hypercube"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/spec"
)

const (
	TCPCode           = 0
	UDPCode           = 1
	ICMPCode          = 2
	MinICMPType int64 = 0
	MaxICMPType int64 = netp.InformationReply
	MinICMPCode int64 = 0
	MaxICMPCode int64 = 5
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
	icmpType Dimension = "icmpType"
	icmpCode Dimension = "icmpCode"
)

const propertySeparator string = " "

// dimensionsList is the ordered list of dimensions in the Set object
// this should be the only place where the order is hard-coded
var dimensionsList = []Dimension{
	protocol,
	srcPort,
	dstPort,
	icmpType,
	icmpCode,
}

func entireDimension(dim Dimension) *interval.CanonicalSet {
	switch dim {
	case protocol:
		return interval.New(minProtocol, maxProtocol).ToSet()
	case srcPort:
		return interval.New(MinPort, MaxPort).ToSet()
	case dstPort:
		return interval.New(MinPort, MaxPort).ToSet()
	case icmpType:
		return interval.New(MinICMPType, MaxICMPType).ToSet()
	case icmpCode:
		return interval.New(MinICMPCode, MaxICMPCode).ToSet()
	}
	return nil
}

type Set struct {
	connectionProperties *hypercube.CanonicalSet
	IsStateful           StatefulState
}

func None() *Set {
	return &Set{connectionProperties: hypercube.NewCanonicalSet(len(dimensionsList))}
}

func All() *Set {
	all := make([]*interval.CanonicalSet, len(dimensionsList))
	for i := range dimensionsList {
		all[i] = entireDimension(dimensionsList[i])
	}
	return &Set{connectionProperties: hypercube.FromCube(all)}
}

var all = All()

func (c *Set) IsAll() bool {
	return c.Equal(all)
}

func (c *Set) Equal(other *Set) bool {
	return c.connectionProperties.Equal(other.connectionProperties)
}

func (c *Set) Copy() *Set {
	return &Set{
		connectionProperties: c.connectionProperties.Copy(),
		IsStateful:           c.IsStateful,
	}
}

func (c *Set) Intersect(other *Set) *Set {
	return &Set{connectionProperties: c.connectionProperties.Intersect(other.connectionProperties)}
}

func (c *Set) IsEmpty() bool {
	return c.connectionProperties.IsEmpty()
}

func (c *Set) Union(other *Set) *Set {
	if other.IsEmpty() {
		return c.Copy()
	}
	if c.IsEmpty() {
		return other.Copy()
	}
	return &Set{
		connectionProperties: c.connectionProperties.Union(other.connectionProperties),
	}
}

// Subtract
// ToDo: Subtract seems to ignore IsStateful (see https://github.com/np-guard/vpc-network-config-analyzer/issues/199):
//  1. is the delta connection stateful
//  2. connectionProperties is identical but c stateful while other is not
//     the 2nd item can be computed here, with enhancement to relevant structure
//     the 1st can not since we do not know where exactly the statefulness came from
func (c *Set) Subtract(other *Set) *Set {
	if c.IsEmpty() {
		return None()
	}
	if other.IsEmpty() {
		return c.Copy()
	}
	return &Set{connectionProperties: c.connectionProperties.Subtract(other.connectionProperties)}
}

// ContainedIn returns true if c is subset of other
func (c *Set) ContainedIn(other *Set) bool {
	res, err := c.connectionProperties.ContainedIn(other.connectionProperties)
	if err != nil {
		log.Panicf("invalid connection set. %e", err)
	}
	return res
}

func protocolStringToCode(protocol netp.ProtocolString) int64 {
	switch protocol {
	case netp.ProtocolStringTCP:
		return TCPCode
	case netp.ProtocolStringUDP:
		return UDPCode
	case netp.ProtocolStringICMP:
		return ICMPCode
	}
	log.Panicf("Impossible protocol code %v", protocol)
	return 0
}

func cube(protocolString netp.ProtocolString,
	srcMinP, srcMaxP, dstMinP, dstMaxP,
	icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax int64) *Set {
	protocol := protocolStringToCode(protocolString)
	return &Set{
		connectionProperties: hypercube.Cube(protocol, protocol,
			srcMinP, srcMaxP, dstMinP, dstMaxP,
			icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax)}
}

func TCPorUDPConnection(protocol netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *Set {
	return cube(protocol,
		srcMinP, srcMaxP, dstMinP, dstMaxP,
		MinICMPType, MaxICMPType, MinICMPCode, MaxICMPCode)
}

func ICMPConnection(icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax int64) *Set {
	return cube(netp.ProtocolStringICMP,
		MinPort, MaxPort, MinPort, MaxPort,
		icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax)
}

func protocolStringFromCode(protocolCode int64) netp.ProtocolString {
	switch protocolCode {
	case TCPCode:
		return netp.ProtocolStringTCP
	case UDPCode:
		return netp.ProtocolStringUDP
	case ICMPCode:
		return netp.ProtocolStringICMP
	}
	log.Panicf("impossible protocol code %v", protocolCode)
	return ""
}

func getDimensionString(cube []*interval.CanonicalSet, dim Dimension) string {
	dimValue := cubeAt(cube, dim)
	if dimValue.Equal(entireDimension(dim)) {
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
	case icmpType:
		return "icmp-type: " + dimValue.String()
	case icmpCode:
		return "icmp-code: " + dimValue.String()
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

func getConnsCubeStr(cube []*interval.CanonicalSet) string {
	protocols := cubeAt(cube, protocol)
	tcpOrUDP := protocols.Contains(TCPCode) || protocols.Contains(UDPCode)
	icmp := protocols.Contains(ICMPCode)
	switch {
	case tcpOrUDP && !icmp:
		return joinNonEmpty(
			getDimensionString(cube, protocol),
			getDimensionString(cube, srcPort),
			getDimensionString(cube, dstPort),
		)
	case icmp && !tcpOrUDP:
		return joinNonEmpty(
			getDimensionString(cube, protocol),
			getDimensionString(cube, icmpType),
			getDimensionString(cube, icmpCode),
		)
	default:
		// TODO: make sure other dimension values are full
		return getDimensionString(cube, protocol)
	}
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
	for _, cube := range c.connectionProperties.GetCubesList() {
		resStrings = append(resStrings, getConnsCubeStr(cube))
	}

	sort.Strings(resStrings)
	return strings.Join(resStrings, "; ")
}

func cubeAt(cube []*interval.CanonicalSet, dim Dimension) *interval.CanonicalSet {
	return cube[slices.Index(dimensionsList, dim)]
}

func getCubeAsTCPItems(cube []*interval.CanonicalSet, protocol spec.TcpUdpProtocol) []spec.TcpUdp {
	tcpItemsTemp := []spec.TcpUdp{}
	tcpItemsFinal := []spec.TcpUdp{}
	// consider src ports
	srcPorts := cubeAt(cube, srcPort)
	if !srcPorts.Equal(entireDimension(srcPort)) {
		// iterate the interval in the interval-set
		for _, interval := range srcPorts.Intervals() {
			tcpRes := spec.TcpUdp{Protocol: protocol, MinSourcePort: int(interval.Start()), MaxSourcePort: int(interval.End())}
			tcpItemsTemp = append(tcpItemsTemp, tcpRes)
		}
	} else {
		tcpItemsTemp = append(tcpItemsTemp, spec.TcpUdp{Protocol: protocol})
	}
	// consider dst ports
	dstPorts := cubeAt(cube, dstPort)
	if !dstPorts.Equal(entireDimension(dstPort)) {
		// iterate the interval in the interval-set
		for _, interval := range dstPorts.Intervals() {
			for _, tcpItemTemp := range tcpItemsTemp {
				tcpRes := spec.TcpUdp{
					Protocol:           protocol,
					MinSourcePort:      tcpItemTemp.MinSourcePort,
					MaxSourcePort:      tcpItemTemp.MaxSourcePort,
					MinDestinationPort: int(interval.Start()),
					MaxDestinationPort: int(interval.End()),
				}
				tcpItemsFinal = append(tcpItemsFinal, tcpRes)
			}
		}
	} else {
		tcpItemsFinal = tcpItemsTemp
	}
	return tcpItemsFinal
}

func getCubeAsICMPItems(cube []*interval.CanonicalSet) []spec.Icmp {
	icmpTypes := cubeAt(cube, icmpType)
	icmpCodes := cubeAt(cube, icmpCode)
	allTypes := icmpTypes.Equal(entireDimension(icmpType))
	allCodes := icmpCodes.Equal(entireDimension(icmpCode))
	switch {
	case allTypes && allCodes:
		return []spec.Icmp{{Protocol: spec.IcmpProtocolICMP}}
	case allTypes:
		// This does not really make sense: not all types can have all codes
		res := []spec.Icmp{}
		for _, code64 := range icmpCodes.Elements() {
			code := int(code64)
			res = append(res, spec.Icmp{Protocol: spec.IcmpProtocolICMP, Code: &code})
		}
		return res
	case allCodes:
		res := []spec.Icmp{}
		for _, type64 := range icmpTypes.Elements() {
			t := int(type64)
			res = append(res, spec.Icmp{Protocol: spec.IcmpProtocolICMP, Type: &t})
		}
		return res
	default:
		res := []spec.Icmp{}
		// iterate both codes and types
		for _, type64 := range icmpTypes.Elements() {
			t := int(type64)
			for _, code64 := range icmpCodes.Elements() {
				code := int(code64)
				res = append(res, spec.Icmp{Protocol: spec.IcmpProtocolICMP, Type: &t, Code: &code})
			}
		}
		return res
	}
}

type Details spec.ProtocolList

func ToJSON(c *Set) Details {
	if c == nil {
		return nil
	}
	if c.IsAll() {
		return Details(spec.ProtocolList{spec.AnyProtocol{Protocol: spec.AnyProtocolProtocolANY}})
	}
	res := spec.ProtocolList{}

	cubes := c.connectionProperties.GetCubesList()
	for _, cube := range cubes {
		protocols := cubeAt(cube, protocol)
		if protocols.Contains(TCPCode) {
			tcpItems := getCubeAsTCPItems(cube, spec.TcpUdpProtocolTCP)
			for _, item := range tcpItems {
				res = append(res, item)
			}
		}
		if protocols.Contains(UDPCode) {
			udpItems := getCubeAsTCPItems(cube, spec.TcpUdpProtocolUDP)
			for _, item := range udpItems {
				res = append(res, item)
			}
		}
		if protocols.Contains(ICMPCode) {
			icmpItems := getCubeAsICMPItems(cube)
			for _, item := range icmpItems {
				res = append(res, item)
			}
		}
	}

	return Details(res)
}
