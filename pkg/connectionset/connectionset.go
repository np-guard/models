package connectionset

import (
	"log"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/hypercubes"
	"github.com/np-guard/models/pkg/intervals"
)

const (
	ICMPCode          = -1
	TCPCode           = 0
	UDPCode           = 1
	MinICMPtype int64 = 0
	MaxICMPtype int64 = informationReply
	MinICMPcode int64 = 0
	MaxICMPcode int64 = 5
	minProtocol int64 = ICMPCode
	maxProtocol int64 = UDPCode
	MinPort     int64 = 1
	MaxPort     int64 = 65535
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

// dimensionsList is the ordered list of dimensions in the ConnectionSet object
// this should be the only place where the order is hard-coded
var dimensionsList = []Dimension{
	protocol,
	srcPort,
	dstPort,
	icmpType,
	icmpCode,
}

func entireDimension(dim Dimension) *intervals.CanonicalIntervalSet {
	switch dim {
	case protocol:
		return intervals.CreateFromInterval(minProtocol, maxProtocol)
	case srcPort:
		return intervals.CreateFromInterval(MinPort, MaxPort)
	case dstPort:
		return intervals.CreateFromInterval(MinPort, MaxPort)
	case icmpType:
		return intervals.CreateFromInterval(MinICMPtype, MaxICMPtype)
	case icmpCode:
		return intervals.CreateFromInterval(MinICMPcode, MaxICMPcode)
	}
	return nil
}

func getDimensionDomainsList() []*intervals.CanonicalIntervalSet {
	res := make([]*intervals.CanonicalIntervalSet, len(dimensionsList))
	for i := range dimensionsList {
		res[i] = entireDimension(dimensionsList[i])
	}
	return res
}

type ConnectionSet struct {
	AllowAll             bool
	connectionProperties *hypercubes.CanonicalHypercubeSet
	IsStateful           int // default is StatefulUnknown
}

func NewConnectionSet(all bool) *ConnectionSet {
	return &ConnectionSet{AllowAll: all, connectionProperties: hypercubes.NewCanonicalHypercubeSet(numDimensions)}
}

func NewConnectionSetWithCube(cube *hypercubes.CanonicalHypercubeSet) *ConnectionSet {
	res := NewConnectionSet(false)
	res.connectionProperties.Union(cube)
	if res.isAllConnectionsWithoutAllowAll() {
		return NewConnectionSet(true)
	}
	return res
}

func (conn *ConnectionSet) Copy() *ConnectionSet {
	return &ConnectionSet{
		AllowAll:             conn.AllowAll,
		connectionProperties: conn.connectionProperties.Copy(),
		IsStateful:           conn.IsStateful,
	}
}

func (conn *ConnectionSet) Intersection(other *ConnectionSet) *ConnectionSet {
	if other.AllowAll {
		return conn.Copy()
	}
	if conn.AllowAll {
		return other.Copy()
	}
	return &ConnectionSet{AllowAll: false, connectionProperties: conn.connectionProperties.Intersection(other.connectionProperties)}
}

func (conn *ConnectionSet) IsEmpty() bool {
	if conn.AllowAll {
		return false
	}
	return conn.connectionProperties.IsEmpty()
}

func (conn *ConnectionSet) Union(other *ConnectionSet) *ConnectionSet {
	if conn.AllowAll || other.AllowAll {
		return NewConnectionSet(true)
	}
	if other.IsEmpty() {
		return conn.Copy()
	}
	if conn.IsEmpty() {
		return other.Copy()
	}
	res := &ConnectionSet{
		AllowAll:             false,
		connectionProperties: conn.connectionProperties.Union(other.connectionProperties),
	}
	if res.isAllConnectionsWithoutAllowAll() {
		return NewConnectionSet(true)
	}
	return res
}

func getAllPropertiesObject() *hypercubes.CanonicalHypercubeSet {
	return hypercubes.CreateFromCube(getDimensionDomainsList())
}

func (conn *ConnectionSet) isAllConnectionsWithoutAllowAll() bool {
	if conn.AllowAll {
		return false
	}
	return conn.connectionProperties.Equals(getAllPropertiesObject())
}

// Subtract
// ToDo: Subtract seems to ignore IsStateful (see https://github.com/np-guard/vpc-network-config-analyzer/issues/199):
//  1. is the delta connection stateful
//  2. connectionProperties is identical but conn stateful while other is not
//     the 2nd item can be computed here, with enhancement to relevant structure
//     the 1st can not since we do not know where exactly the statefulness came from
func (conn *ConnectionSet) Subtract(other *ConnectionSet) *ConnectionSet {
	if conn.IsEmpty() || other.IsEmpty() {
		return conn
	}
	if other.AllowAll {
		return NewConnectionSet(false)
	}
	var connProperties *hypercubes.CanonicalHypercubeSet
	if conn.AllowAll {
		connProperties = getAllPropertiesObject()
	} else {
		connProperties = conn.connectionProperties
	}
	return &ConnectionSet{AllowAll: false, connectionProperties: connProperties.Subtraction(other.connectionProperties)}
}

func (conn *ConnectionSet) ContainedIn(other *ConnectionSet) (bool, error) {
	if other.AllowAll {
		return true, nil
	}
	if conn.AllowAll {
		return false, nil
	}
	res, err := conn.connectionProperties.ContainedIn(other.connectionProperties)
	return res, err
}

func ProtocolStringToCode(protocol ProtocolStr) int64 {
	switch protocol {
	case ProtocolStringTCP:
		return TCPCode
	case ProtocolStringUDP:
		return UDPCode
	case ProtocolStringICMP:
		return ICMPCode
	}
	log.Fatalf("Impossible protocol code %v", protocol)
	return 0
}

func (conn *ConnectionSet) addConnection(protocol ProtocolStr,
	srcMinP, srcMaxP, dstMinP, dstMaxP,
	icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax int64) {
	code := ProtocolStringToCode(protocol)
	cube := hypercubes.CreateFromCubeShort(code, code,
		srcMinP, srcMaxP, dstMinP, dstMaxP,
		icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax)
	conn.connectionProperties = conn.connectionProperties.Union(cube)
	// check if all connections allowed after this union
	if conn.isAllConnectionsWithoutAllowAll() {
		conn.AllowAll = true
		conn.connectionProperties = hypercubes.NewCanonicalHypercubeSet(numDimensions)
	}
}

func (conn *ConnectionSet) AddTCPorUDPConn(protocol ProtocolStr, srcMinP, srcMaxP, dstMinP, dstMaxP int64) {
	conn.addConnection(protocol,
		srcMinP, srcMaxP, dstMinP, dstMaxP,
		MinICMPtype, MaxICMPtype, MinICMPcode, MaxICMPcode)
}

func (conn *ConnectionSet) AddICMPConnection(icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax int64) {
	conn.addConnection(ProtocolStringICMP,
		MinPort, MaxPort, MinPort, MaxPort,
		icmpTypeMin, icmpTypeMax, icmpCodeMin, icmpCodeMax)
}

func (conn *ConnectionSet) Equal(other *ConnectionSet) bool {
	if conn.AllowAll != other.AllowAll {
		return false
	}
	if conn.AllowAll {
		return true
	}
	return conn.connectionProperties.Equals(other.connectionProperties)
}

func getProtocolStr(p int64) ProtocolStr {
	switch p {
	case TCPCode:
		return ProtocolStringTCP
	case UDPCode:
		return ProtocolStringUDP
	case ICMPCode:
		return ProtocolStringICMP
	}
	log.Fatalf("Impossible protocol value %v", p)
	return ""
}

func getDimensionStr(dimValue *intervals.CanonicalIntervalSet, dim Dimension) string {
	domainValues := entireDimension(dim)
	if dimValue.Equal(*domainValues) {
		// avoid adding dimension str on full dimension values
		return ""
	}
	switch dim {
	case protocol:
		pList := []string{}
		for p := minProtocol; p <= maxProtocol; p++ {
			pList = append(pList, string(getProtocolStr(p)))
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

func getICMPbasedCubeStr(protocolsValues, icmpTypeValues, icmpCodeValues *intervals.CanonicalIntervalSet) string {
	strList := []string{
		getDimensionStr(protocolsValues, protocol),
		getDimensionStr(icmpTypeValues, icmpType),
		getDimensionStr(icmpCodeValues, icmpCode),
	}
	return strings.Join(filterEmptyPropertiesStr(strList), propertySeparator)
}

func getPortBasedCubeStr(protocolsValues, srcPortsValues, dstPortsValues *intervals.CanonicalIntervalSet) string {
	strList := []string{
		getDimensionStr(protocolsValues, protocol),
		getDimensionStr(srcPortsValues, srcPort),
		getDimensionStr(dstPortsValues, dstPort),
	}
	return strings.Join(filterEmptyPropertiesStr(strList), propertySeparator)
}

func getMixedProtocolsCubeStr(protocols *intervals.CanonicalIntervalSet) string {
	// TODO: make sure other dimension values are full
	return getDimensionStr(protocols, protocol)
}

func getConnsCubeStr(cube []*intervals.CanonicalIntervalSet) string {
	protocols := cube[protocol]
	if (protocols.Contains(TCPCode) || protocols.Contains(UDPCode)) && !protocols.Contains(ICMPCode) {
		return getPortBasedCubeStr(cube[protocol], cube[srcPort], cube[dstPort])
	}
	if protocols.Contains(ICMPCode) && !(protocols.Contains(TCPCode) || protocols.Contains(UDPCode)) {
		return getICMPbasedCubeStr(cube[protocol], cube[icmpType], cube[icmpCode])
	}
	return getMixedProtocolsCubeStr(protocols)
}

// String returns a string representation of a ConnectionSet object
func (conn *ConnectionSet) String() string {
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

func getCubeAsTCPItems(cube []*intervals.CanonicalIntervalSet, protocol TransportLayerProtocolName) []Protocol {
	tcpItemsTemp := []Protocol{}
	// consider src ports
	srcPorts := cube[srcPort]
	if srcPorts.Equal(*entireDimension(srcPort)) {
		tcpItemsTemp = append(tcpItemsTemp, TCPUDP{Protocol: protocol})
	} else {
		// iterate the intervals in the interval-set
		for _, interval := range srcPorts.IntervalSet {
			tcpRes := TCPUDP{
				Protocol: protocol,
				PortRangePair: PortRangePair{
					SrcPort: PortRange{Min: int(interval.Start), Max: int(interval.End)},
				},
			}
			tcpItemsTemp = append(tcpItemsTemp, tcpRes)
		}
	}
	// consider dst ports
	dstPorts := cube[dstPort]
	if dstPorts.Equal(*entireDimension(dstPort)) {
		return tcpItemsTemp
	}
	tcpItemsFinal := []Protocol{}
	for _, interval := range dstPorts.IntervalSet {
		for _, tcpItemTemp := range tcpItemsTemp {
			item, _ := tcpItemTemp.(TCPUDP)
			tcpItemsFinal = append(tcpItemsFinal, TCPUDP{
				Protocol: protocol,
				PortRangePair: PortRangePair{
					SrcPort: item.PortRangePair.SrcPort,
					DstPort: PortRange{int(interval.Start), int(interval.End)},
				},
			})
		}
	}
	return tcpItemsFinal
}

func getCubeAsICMPItems(cube []*intervals.CanonicalIntervalSet) []Protocol {
	icmpTypes := cube[icmpType]
	icmpCodes := cube[icmpCode]
	if icmpCodes.Equal(*entireDimension(icmpCode)) {
		if icmpTypes.Equal(*entireDimension(icmpType)) {
			return []Protocol{ICMP{}}
		}
		res := []Protocol{}
		for _, t := range icmpTypes.Elements() {
			res = append(res, ICMP{ICMPCodeType: &ICMPCodeType{Type: t}})
		}
		return res
	}

	// iterate both codes and types
	res := []Protocol{}
	for _, t := range icmpTypes.Elements() {
		codes := icmpCodes.Elements()
		for i := range codes {
			c := codes[i]
			if ValidateICMP(t, c) == nil {
				res = append(res, ICMP{ICMPCodeType: &ICMPCodeType{Type: t, Code: &c}})
			}
		}
	}
	return res
}

type ConnDetails []Protocol

func ConnToJSONRep(c *ConnectionSet) ConnDetails {
	if c == nil {
		return nil // one of the connections in connectionDiff can be empty
	}
	if c.AllowAll {
		return []Protocol{}
	}
	var res []Protocol

	cubes := c.connectionProperties.GetCubesList()
	for _, cube := range cubes {
		protocols := cube[protocol]
		if protocols.Contains(TCPCode) {
			res = append(res, getCubeAsTCPItems(cube, TCP)...)
		}
		if protocols.Contains(UDPCode) {
			res = append(res, getCubeAsTCPItems(cube, UDP)...)
		}
		if protocols.Contains(ICMPCode) {
			res = append(res, getCubeAsICMPItems(cube)...)
		}
	}

	return res
}

// NewTCPConnectionSet returns a ConnectionSet object with TCPCode protocol (all ports)
func NewTCPConnectionSet() *ConnectionSet {
	res := NewConnectionSet(false)
	res.AddTCPorUDPConn(ProtocolStringTCP, MinPort, MaxPort, MinPort, MaxPort)
	return res
}
