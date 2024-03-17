// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package connection

import (
	"slices"

	"github.com/np-guard/models/pkg/hypercube"
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
func (conn *Set) EnhancedString() string {
	if conn.IsStateful == StatefulFalse {
		return conn.String() + " *"
	}
	return conn.String()
}

func newTCPSet() *Set {
	return TCPorUDPConnection(netp.ProtocolStringTCP, MinPort, MaxPort, MinPort, MaxPort)
}

// ConnectionWithStatefulness updates `conn` object with `IsStateful` property, based on input `secondDirectionConn`.
// `conn` represents a src-to-dst connection, and `secondDirectionConn` represents dst-to-src connection.
// The property `IsStateful` of `conn` is set as `StatefulFalse` if there is at least some subset within TCP from `conn`
// which is not stateful (such that the response direction for this subset is not enabled).
// This function also returns a connection object with the exact subset of the stateful part (within TCP)
// from the entire connection `conn`, and with the original connections on other protocols.
func (conn *Set) ConnectionWithStatefulness(secondDirectionConn *Set) *Set {
	connTCP := conn.Intersect(newTCPSet())
	if connTCP.IsEmpty() {
		conn.IsStateful = StatefulTrue
		return conn
	}
	statefulCombinedConnTCP := connTCP.connTCPWithStatefulness(secondDirectionConn.Intersect(newTCPSet()))
	conn.IsStateful = connTCP.IsStateful
	return conn.Subtract(connTCP).Union(statefulCombinedConnTCP)
}

// connTCPWithStatefulness assumes that both `conn` and `secondDirectionConn` are within TCP.
// it assigns IsStateful a value within `conn`, and returns the subset from `conn` which is stateful.
func (conn *Set) connTCPWithStatefulness(secondDirectionConn *Set) *Set {
	// flip src/dst ports before intersection
	statefulCombinedConn := conn.Intersect(secondDirectionConn.switchSrcDstPortsOnTCP())
	if conn.Equal(statefulCombinedConn) {
		conn.IsStateful = StatefulTrue
	} else {
		conn.IsStateful = StatefulFalse
	}
	return statefulCombinedConn
}

// switchSrcDstPortsOnTCP returns a new Set object, built from the input Set object.
// It assumes the input connection object is only within TCP protocol.
// For TCP the src and dst ports on relevant cubes are being switched.
func (conn *Set) switchSrcDstPortsOnTCP() *Set {
	if conn.AllowAll || conn.IsEmpty() {
		return conn.Copy()
	}
	res := None()
	for _, cube := range conn.connectionProperties.GetCubesList() {
		// assuming cube[protocol] contains TCP only
		// no need to switch if src equals dst
		if !cube[srcPort].Equal(cube[dstPort]) {
			// Shallow clone should be enough, since we do shallow swap:
			cube = slices.Clone(cube)
			cube[srcPort], cube[dstPort] = cube[dstPort], cube[srcPort]
		}
		res.connectionProperties = res.connectionProperties.Union(hypercube.FromCube(cube))
	}
	return res
}
