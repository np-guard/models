/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset

import (
	"github.com/np-guard/models/pkg/ds"
)

// ConnectionSet captures a set of connections for tuples of (src IP range, dst IP range, connection.Set),
// where connection.Set is a set of TCP/UPD/ICMP with their properties (ports/icmp type&code)
type ConnectionSet struct {
	props ds.TripleSet[*IPBlock, *IPBlock, *TransportSet]
}

// NewConnectionSet returns an empty ConnectionSet
func NewConnectionSet() *ConnectionSet {
	return &ConnectionSet{props: ds.NewLeftTripleSet[*IPBlock, *IPBlock, *TransportSet]()}
}

// Equal returns true is this ConnectionSet captures the exact same set of connections as `other` does.
func (c *ConnectionSet) Equal(other *ConnectionSet) bool {
	return c.props.Equal(other.props)
}

// Copy returns new ConnectionSet object with same set of connections as current one
func (c *ConnectionSet) Copy() *ConnectionSet {
	return &ConnectionSet{
		props: c.props.Copy(),
	}
}

// Intersect returns a ConnectionSet object with connection tuples that result from intersecion of
// this and `other` sets
func (c *ConnectionSet) Intersect(other *ConnectionSet) *ConnectionSet {
	return &ConnectionSet{props: c.props.Intersect(other.props)}
}

// IsEmpty returns true of the ConnectionSet is empty
func (c *ConnectionSet) IsEmpty() bool {
	return c.props.IsEmpty()
}

// Union returns a ConnectionSet object with connection tuples that result from union of
// this and `other` sets
func (c *ConnectionSet) Union(other *ConnectionSet) *ConnectionSet {
	if other.IsEmpty() {
		return c.Copy()
	}
	if c.IsEmpty() {
		return other.Copy()
	}
	return &ConnectionSet{
		props: c.props.Union(other.props),
	}
}

// Subtract returns a ConnectionSet object with connection tuples that result from subtraction of
// `other` from this set
func (c *ConnectionSet) Subtract(other *ConnectionSet) *ConnectionSet {
	if other.IsEmpty() {
		return c.Copy()
	}
	return &ConnectionSet{props: c.props.Subtract(other.props)}
}

// IsSubset returns true if c is subset of other
func (c *ConnectionSet) IsSubset(other *ConnectionSet) bool {
	return c.props.IsSubset(other.props)
}

// ConnectionSetFrom returns a new ConnectionSet object from input src, dst IP-ranges sets ands
// TransportSet connections
func ConnectionSetFrom(src, dst *IPBlock, conn *TransportSet) *ConnectionSet {
	return &ConnectionSet{props: ds.CartesianLeftTriple(src, dst, conn)}
}
