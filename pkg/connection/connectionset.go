// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package connection

import (
	"log"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/hypercube"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
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

type Dimension int

const (
	protocol      Dimension = 0
	srcPort       Dimension = 1
	dstPort       Dimension = 2
	icmpType      Dimension = 3
	icmpCode      Dimension = 4
	numDimensions           = 5
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

func getDimensionDomainsList() []*interval.CanonicalSet {
	res := make([]*interval.CanonicalSet, len(dimensionsList))
	for i := range dimensionsList {
		res[i] = entireDimension(dimensionsList[i])
	}
	return res
}

type Set struct {
	connectionProperties *hypercube.CanonicalSet
	IsStateful           StatefulState
}

func None() *Set {
	return &Set{connectionProperties: hypercube.NewCanonicalSet(numDimensions)}
}

func All() *Set {
	return &Set{connectionProperties: hypercube.FromCube(getDimensionDomainsList())}
}

var all = All()

func (c *Set) IsAll() bool {
	return c.Equal(all)
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
		log.Fatalf("invalid connection set. %e", err)
	}
	return res
}

func ProtocolStringToCode(protocol netp.ProtocolString) int64 {
	switch protocol {
	case netp.ProtocolStringTCP:
		return TCPCode
	case netp.ProtocolStringUDP:
		return UDPCode
	case netp.ProtocolStringICMP:
		return ICMPCode
	}
	log.Fatalf("Impossible protocol code %v", protocol)
	return 0
}

func (c *Set) addConnection(protocol netp.ProtocolString,
	srcMinP, srcMaxP, dstMinP, dstMaxP,
	icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax int64) {
	code := ProtocolStringToCode(protocol)
	cube := hypercube.Cube(code, code,
		srcMinP, srcMaxP, dstMinP, dstMaxP,
		icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax)
	c.connectionProperties = c.connectionProperties.Union(cube)
}

func TCPorUDPConnection(protocol netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *Set {
	c := None()
	c.addConnection(protocol,
		srcMinP, srcMaxP, dstMinP, dstMaxP,
		MinICMPType, MaxICMPType, MinICMPCode, MaxICMPCode)
	return c
}

func ICMPConnection(icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax int64) *Set {
	c := None()
	c.addConnection(netp.ProtocolStringICMP,
		MinPort, MaxPort, MinPort, MaxPort,
		icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax)
	return c
}

func (c *Set) Equal(other *Set) bool {
	return c.connectionProperties.Equal(other.connectionProperties)
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
	log.Fatalf("impossible protocol code %v", protocolCode)
	return ""
}

func getDimensionString(dimValue *interval.CanonicalSet, dim Dimension) string {
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

func filterEmptyPropertiesStr(inputList []string) []string {
	res := []string{}
	for _, propertyStr := range inputList {
		if propertyStr != "" {
			res = append(res, propertyStr)
		}
	}
	return res
}

func getICMPbasedCubeStr(protocolsValues, icmpTypeValues, icmpCodeValues *interval.CanonicalSet) string {
	strList := []string{
		getDimensionString(protocolsValues, protocol),
		getDimensionString(icmpTypeValues, icmpType),
		getDimensionString(icmpCodeValues, icmpCode),
	}
	return strings.Join(filterEmptyPropertiesStr(strList), propertySeparator)
}

func getPortBasedCubeStr(protocolsValues, srcPortsValues, dstPortsValues *interval.CanonicalSet) string {
	strList := []string{
		getDimensionString(protocolsValues, protocol),
		getDimensionString(srcPortsValues, srcPort),
		getDimensionString(dstPortsValues, dstPort),
	}
	return strings.Join(filterEmptyPropertiesStr(strList), propertySeparator)
}

func getMixedProtocolsCubeStr(protocols *interval.CanonicalSet) string {
	// TODO: make sure other dimension values are full
	return getDimensionString(protocols, protocol)
}

func getConnsCubeStr(cube []*interval.CanonicalSet) string {
	protocols := cube[protocol]
	if (protocols.Contains(TCPCode) || protocols.Contains(UDPCode)) && !protocols.Contains(ICMPCode) {
		return getPortBasedCubeStr(protocols, cube[srcPort], cube[dstPort])
	}
	if protocols.Contains(ICMPCode) && !(protocols.Contains(TCPCode) || protocols.Contains(UDPCode)) {
		return getICMPbasedCubeStr(protocols, cube[icmpType], cube[icmpCode])
	}
	return getMixedProtocolsCubeStr(protocols)
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

func getCubeAsTCPorUDPItems(cube []*interval.CanonicalSet, isTCP bool) []netp.Protocol {
	tcpItemsTemp := []netp.Protocol{}
	// consider src ports
	srcPorts := cube[srcPort]
	if srcPorts.Equal(entireDimension(srcPort)) {
		tcpItemsTemp = append(tcpItemsTemp, netp.TCPUDP{IsTCP: isTCP})
	} else {
		// iterate the intervals in the interval-set
		for _, portRange := range srcPorts.Intervals() {
			tcpRes := netp.TCPUDP{
				IsTCP: isTCP,
				PortRangePair: netp.PortRangePair{
					SrcPort: portRange,
					DstPort: interval.Interval{Start: netp.MinPort, End: netp.MaxPort},
				},
			}
			tcpItemsTemp = append(tcpItemsTemp, tcpRes)
		}
	}
	// consider dst ports
	dstPorts := cube[dstPort]
	if dstPorts.Equal(entireDimension(dstPort)) {
		return tcpItemsTemp
	}
	tcpItemsFinal := []netp.Protocol{}
	for _, portRange := range dstPorts.Intervals() {
		for _, tcpItemTemp := range tcpItemsTemp {
			item, _ := tcpItemTemp.(netp.TCPUDP)
			tcpItemsFinal = append(tcpItemsFinal, netp.TCPUDP{
				IsTCP: isTCP,
				PortRangePair: netp.PortRangePair{
					SrcPort: item.PortRangePair.SrcPort,
					DstPort: portRange,
				},
			})
		}
	}
	return tcpItemsFinal
}

func getCubeAsICMPItems(cube []*interval.CanonicalSet) []netp.Protocol {
	icmpTypes := cube[icmpType]
	icmpCodes := cube[icmpCode]
	if icmpCodes.Equal(entireDimension(icmpCode)) {
		if icmpTypes.Equal(entireDimension(icmpType)) {
			return []netp.Protocol{netp.ICMP{}}
		}
		res := []netp.Protocol{}
		for _, t := range icmpTypes.Elements() {
			icmp, err := netp.NewICMP(&netp.ICMPTypeCode{Type: int(t)})
			if err != nil {
				log.Panic(err)
			}
			res = append(res, icmp)
		}
		return res
	}

	// iterate both codes and types
	res := []netp.Protocol{}
	for _, t := range icmpTypes.Elements() {
		codes := icmpCodes.Elements()
		for i := range codes {
			// TODO: merge when all codes for certain type exist
			c := int(codes[i])
			icmp, err := netp.NewICMP(&netp.ICMPTypeCode{Type: int(t), Code: &c})
			if err != nil {
				log.Panic(err)
			}
			res = append(res, icmp)
		}
	}
	return res
}

type Details []netp.Protocol

func ToJSON(c *Set) Details {
	if c == nil {
		return nil // one of the connections in connectionDiff can be empty
	}
	if c.IsAll() {
		return []netp.Protocol{netp.AnyProtocol{}}
	}
	var res []netp.Protocol

	cubes := c.connectionProperties.GetCubesList()
	for _, cube := range cubes {
		protocols := cube[protocol]
		if protocols.Contains(TCPCode) {
			res = append(res, getCubeAsTCPorUDPItems(cube, true)...)
		}
		if protocols.Contains(UDPCode) {
			res = append(res, getCubeAsTCPorUDPItems(cube, false)...)
		}
		if protocols.Contains(ICMPCode) {
			res = append(res, getCubeAsICMPItems(cube)...)
		}
	}

	return res
}
