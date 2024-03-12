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
	MinICMPtype int64 = 0
	MaxICMPtype int64 = netp.InformationReply
	MinICMPcode int64 = 0
	MaxICMPcode int64 = 5
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
		return interval.CreateSetFromInterval(minProtocol, maxProtocol)
	case srcPort:
		return interval.CreateSetFromInterval(MinPort, MaxPort)
	case dstPort:
		return interval.CreateSetFromInterval(MinPort, MaxPort)
	case icmpType:
		return interval.CreateSetFromInterval(MinICMPtype, MaxICMPtype)
	case icmpCode:
		return interval.CreateSetFromInterval(MinICMPcode, MaxICMPcode)
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
	AllowAll             bool
	connectionProperties *hypercube.CanonicalSet
	IsStateful           StatefulState
}

func newSet(all bool) *Set {
	return &Set{AllowAll: all, connectionProperties: hypercube.NewCanonicalSet(numDimensions)}
}

func All() *Set {
	return newSet(true)
}

func None() *Set {
	return newSet(false)
}

func (conn *Set) Copy() *Set {
	return &Set{
		AllowAll:             conn.AllowAll,
		connectionProperties: conn.connectionProperties.Copy(),
		IsStateful:           conn.IsStateful,
	}
}

func (conn *Set) Intersect(other *Set) *Set {
	if other.AllowAll {
		return conn.Copy()
	}
	if conn.AllowAll {
		return other.Copy()
	}
	return &Set{AllowAll: false, connectionProperties: conn.connectionProperties.Intersect(other.connectionProperties)}
}

func (conn *Set) IsEmpty() bool {
	if conn.AllowAll {
		return false
	}
	return conn.connectionProperties.IsEmpty()
}

func (conn *Set) Union(other *Set) *Set {
	if conn.AllowAll || other.AllowAll {
		return All()
	}
	if other.IsEmpty() {
		return conn.Copy()
	}
	if conn.IsEmpty() {
		return other.Copy()
	}
	res := &Set{
		AllowAll:             false,
		connectionProperties: conn.connectionProperties.Union(other.connectionProperties),
	}
	res.canonicalize()
	return res
}

func getAllPropertiesObject() *hypercube.CanonicalSet {
	return hypercube.FromCube(getDimensionDomainsList())
}

// Subtract
// ToDo: Subtract seems to ignore IsStateful (see https://github.com/np-guard/vpc-network-config-analyzer/issues/199):
//  1. is the delta connection stateful
//  2. connectionProperties is identical but conn stateful while other is not
//     the 2nd item can be computed here, with enhancement to relevant structure
//     the 1st can not since we do not know where exactly the statefulness came from
func (conn *Set) Subtract(other *Set) *Set {
	if conn.IsEmpty() || other.AllowAll {
		return None()
	}
	if other.IsEmpty() {
		return conn.Copy()
	}
	var connProperties *hypercube.CanonicalSet
	if conn.AllowAll {
		connProperties = getAllPropertiesObject()
	} else {
		connProperties = conn.connectionProperties
	}
	return &Set{AllowAll: false, connectionProperties: connProperties.Subtract(other.connectionProperties)}
}

// ContainedIn returns true if conn is subset of other
func (conn *Set) ContainedIn(other *Set) bool {
	if other.AllowAll {
		return true
	}
	if conn.AllowAll {
		return false
	}
	res, err := conn.connectionProperties.ContainedIn(other.connectionProperties)
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

func (conn *Set) addConnection(protocol netp.ProtocolString,
	srcMinP, srcMaxP, dstMinP, dstMaxP,
	icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax int64) {
	code := ProtocolStringToCode(protocol)
	cube := hypercube.FromCubeShort(code, code,
		srcMinP, srcMaxP, dstMinP, dstMaxP,
		icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax)
	conn.connectionProperties = conn.connectionProperties.Union(cube)
	conn.canonicalize()
}

func (conn *Set) canonicalize() {
	if !conn.AllowAll && conn.connectionProperties.Equal(getAllPropertiesObject()) {
		conn.AllowAll = true
		conn.connectionProperties = hypercube.NewCanonicalSet(numDimensions)
	}
}

func TCPorUDPConnection(protocol netp.ProtocolString, srcMinP, srcMaxP, dstMinP, dstMaxP int64) *Set {
	conn := None()
	conn.addConnection(protocol,
		srcMinP, srcMaxP, dstMinP, dstMaxP,
		MinICMPtype, MaxICMPtype, MinICMPcode, MaxICMPcode)
	return conn
}

func ICMPConnection(icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax int64) *Set {
	conn := None()
	conn.addConnection(netp.ProtocolStringICMP,
		MinPort, MaxPort, MinPort, MaxPort,
		icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax)
	return conn
}

func (conn *Set) Equal(other *Set) bool {
	if conn.AllowAll != other.AllowAll {
		return false
	}
	if conn.AllowAll {
		return true
	}
	return conn.connectionProperties.Equal(other.connectionProperties)
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
func (conn *Set) String() string {
	if conn.AllowAll {
		return AllConnections
	} else if conn.IsEmpty() {
		return NoConnections
	}
	resStrings := []string{}
	// get cubes and cube str per each cube
	cubes := conn.connectionProperties.GetCubesList()
	for _, cube := range cubes {
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

func ConnToJSONRep(c *Set) Details {
	if c == nil {
		return nil // one of the connections in connectionDiff can be empty
	}
	if c.AllowAll {
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
