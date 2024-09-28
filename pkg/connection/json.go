/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package connection

import (
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
	"github.com/np-guard/models/pkg/spec"
)

func getCubeAsTCPItems(srcPorts, dstPorts *netset.PortSet, p int64) []spec.TcpUdp {
	protocol := spec.TcpUdpProtocol(netp.ProtocolStringTCP)
	if p == netset.UDPCode {
		protocol = spec.TcpUdpProtocol(netp.ProtocolStringUDP)
	}
	var tcpItemsTemp []spec.TcpUdp
	var tcpItemsFinal []spec.TcpUdp
	// consider src ports
	if !srcPorts.Equal(netset.AllPorts()) {
		// iterate the interval in the interval-set
		for _, span := range srcPorts.Intervals() {
			tcpRes := spec.TcpUdp{Protocol: protocol, MinSourcePort: int(span.Start()), MaxSourcePort: int(span.End())}
			tcpItemsTemp = append(tcpItemsTemp, tcpRes)
		}
	} else {
		tcpItemsTemp = append(tcpItemsTemp, spec.TcpUdp{Protocol: protocol})
	}
	// consider dst ports
	if !dstPorts.Equal(netset.AllPorts()) {
		// iterate the interval in the interval-set
		for _, span := range dstPorts.Intervals() {
			for _, tcpItemTemp := range tcpItemsTemp {
				tcpRes := spec.TcpUdp{
					Protocol:           protocol,
					MinSourcePort:      tcpItemTemp.MinSourcePort,
					MaxSourcePort:      tcpItemTemp.MaxSourcePort,
					MinDestinationPort: int(span.Start()),
					MaxDestinationPort: int(span.End()),
				}
				tcpItemsFinal = append(tcpItemsFinal, tcpRes)
			}
		}
	} else {
		tcpItemsFinal = tcpItemsTemp
	}
	return tcpItemsFinal
}

type Details spec.ProtocolList

// ToJSON returns a `Details` object for JSON representation of the input connection Set.
func ToJSON(c *Set) Details {
	if c == nil {
		return Details{}
	}
	if c.IsAll() {
		return Details(spec.ProtocolList{spec.AnyProtocol{Protocol: spec.AnyProtocolProtocolANY}})
	}
	res := spec.ProtocolList{}

	for _, cube := range c.TCPUDPSet().Partitions() {
		protocols := cube.S1
		for _, p := range protocols.Elements() {
			tcpItems := getCubeAsTCPItems(cube.S2, cube.S3, p)
			for _, item := range tcpItems {
				res = append(res, item)
			}
		}
	}
	for _, item := range c.ICMPSet().Partitions() {
		if item.TypeCode != nil {
			t := item.TypeCode.Type
			res = append(res, spec.Icmp{
				Protocol: spec.IcmpProtocolICMP,
				Type:     &t,
				Code:     item.TypeCode.Code,
			})
		} else {
			res = append(res, spec.Icmp{
				Protocol: spec.IcmpProtocolICMP,
			})
		}
	}

	return Details(res)
}
