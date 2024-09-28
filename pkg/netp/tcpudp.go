/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netp

import (
	"fmt"

	"github.com/np-guard/models/pkg/interval"
)

const MinPort = 1
const MaxPort = 65535

// AllPorts returns an interval representing all possible valid ports.
func AllPorts() interval.Interval {
	return interval.New(MinPort, MaxPort)
}

// TODO: code below can be removed?

// TCPUDP represents a TCP or UDP protocol with contiguous source and destination port ranges.
type TCPUDP struct {
	isTCP    bool
	srcPorts interval.Interval
	dstPorts interval.Interval
}

// SrcPorts returns the source port range.
func (t TCPUDP) SrcPorts() interval.Interval {
	return t.srcPorts
}

// DstPorts returns the destination port range.
func (t TCPUDP) DstPorts() interval.Interval {
	return t.dstPorts
}

// InverseDirection returns a new TCPUDP representing a TCP response with source and destination ports swapped.
// If the current TCPUDP is a UDP protocol, InverseDirection returns nil.
func (t TCPUDP) InverseDirection() Protocol {
	if !t.isTCP {
		return nil
	}
	return TCPUDP{
		isTCP:    true,
		srcPorts: t.dstPorts,
		dstPorts: t.srcPorts,
	}
}

func (t TCPUDP) ProtocolString() ProtocolString {
	if t.isTCP {
		return ProtocolStringTCP
	}
	return ProtocolStringUDP
}

// IsAllPorts returns true if the input port range covers all possible valid ports.
func IsAllPorts(portRange interval.Interval) bool {
	return portRange.Equal(AllPorts())
}

// AllTCPUDP returns a new TCPUDP object representing all possible TCP or UDP connections.
func AllTCPUDP(isTCP bool) *TCPUDP {
	return &TCPUDP{
		isTCP:    isTCP,
		srcPorts: AllPorts(),
		dstPorts: AllPorts(),
	}
}

// NewTCPUDP returns a new TCPUDP object with the specified protocol (TCP/UDP) and port ranges.
func NewTCPUDP(isTCP bool, minSrcPort, maxSrcPort, minDstPort, maxDstPort int) (TCPUDP, error) {
	allPorts := AllPorts()
	srcPorts := interval.New(int64(minSrcPort), int64(maxSrcPort))
	dstPorts := interval.New(int64(minDstPort), int64(maxDstPort))
	if srcPorts.IsEmpty() || dstPorts.IsEmpty() || !srcPorts.IsSubset(allPorts) || !dstPorts.IsSubset(allPorts) {
		return TCPUDP{}, fmt.Errorf("ports must be in the range [%d-%d]; got src=[%d-%d] dst=[%d-%d]",
			MinPort, MaxPort, minSrcPort, maxSrcPort, minDstPort, maxDstPort)
	}
	return TCPUDP{isTCP: isTCP, srcPorts: srcPorts, dstPorts: dstPorts}, nil
}
