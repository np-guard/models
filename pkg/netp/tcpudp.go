// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package netp

import "github.com/np-guard/models/pkg/interval"

const MinPort = 1
const MaxPort = 65535

type PortRangePair struct {
	SrcPort interval.Interval
	DstPort interval.Interval
}

type TCPUDP struct {
	IsTCP         bool
	PortRangePair PortRangePair
}

func (t TCPUDP) InverseDirection() Protocol {
	if !t.IsTCP {
		return nil
	}
	return TCPUDP{
		IsTCP:         true,
		PortRangePair: PortRangePair{SrcPort: t.PortRangePair.DstPort, DstPort: t.PortRangePair.SrcPort},
	}
}

func (t TCPUDP) ProtocolString() ProtocolString {
	if t.IsTCP {
		return ProtocolStringTCP
	}
	return ProtocolStringUDP
}
