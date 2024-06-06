/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package connection

import (
	"slices"

	"github.com/np-guard/models/pkg/netp"
)

// default is StatefulUnknown
type StatefulState int

const (
	// StatefulUnknown is the default value for a Set object,
	StatefulUnknown StatefulState = 0
	// StatefulTrue represents a connection object for which any allowed TCP (on all allowed src/dst ports)
	// has an allowed response connection
	StatefulTrue StatefulState = 1
	// StatefulFalse represents a connection object for which there exists some allowed TCP
	// (on any allowed subset from the allowed src/dst ports) that does not have an allowed response connection
	StatefulFalse StatefulState = 2
)

// EnhancedString returns a connection string with possibly added asterisk for stateless connection
func (c *Set) EnhancedString() string {
	if c.IsStateful == StatefulFalse {
		return c.String() + " *"
	}
	return c.String()
}

func newTCPSet() *Set {
	return TCPorUDPConnection(netp.ProtocolStringTCP, MinPort, MaxPort, MinPort, MaxPort)
}

// WithStatefulness updates `c` object with `IsStateful` property, based on input `secondDirectionConn`.
// `c` represents a src-to-dst connection, and `secondDirectionConn` represents dst-to-src connection.
// The property `IsStateful` of `c` is set as `StatefulFalse` if there is at least some subset within TCP from `c`
// which is not stateful (such that the response direction for this subset is not enabled).
// This function also returns a connection object with the exact subset of the stateful part (within TCP)
// from the entire connection `c`, and with the original connections on other protocols.
func (c *Set) WithStatefulness(secondDirectionConn *Set) *Set {
	connTCP := c.Intersect(newTCPSet())
	if connTCP.IsEmpty() {
		c.IsStateful = StatefulTrue
		return c
	}
	statefulCombinedConnTCP := connTCP.connTCPWithStatefulness(secondDirectionConn.Intersect(newTCPSet()))
	c.IsStateful = connTCP.IsStateful
	return c.Subtract(connTCP).Union(statefulCombinedConnTCP)
}

// connTCPWithStatefulness assumes that both `c` and `secondDirectionConn` are within TCP.
// it assigns IsStateful a value within `c`, and returns the subset from `c` which is stateful.
func (c *Set) connTCPWithStatefulness(secondDirectionConn *Set) *Set {
	// flip src/dst ports before intersection
	statefulCombinedConn := c.Intersect(secondDirectionConn.switchSrcDstPortsOnTCP())
	if c.Equal(statefulCombinedConn) {
		c.IsStateful = StatefulTrue
	} else {
		c.IsStateful = StatefulFalse
	}
	return statefulCombinedConn
}

// switchSrcDstPortsOnTCP returns a new Set object, built from the input Set object.
// It assumes the input connection object is only within TCP protocol.
// For TCP the src and dst ports on relevant cubes are being switched.
func (c *Set) switchSrcDstPortsOnTCP() *Set {
	if c.IsAll() {
		return c.Copy()
	}
	newConn := c.connectionProperties.SwapDimensions(slices.Index(dimensionsList, srcPort), slices.Index(dimensionsList, dstPort))
	return &Set{
		connectionProperties: newConn,
	}
}
