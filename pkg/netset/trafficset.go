/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset

import (
	"fmt"
	"sort"
	"strings"

	"github.com/np-guard/models/pkg/ds"
)

// EndpointsTrafficSet captures a set of traffic attributes for tuples of (source IP range, desination IP range, TransportSet),
// where TransportSet is a set of TCP/UPD/ICMP with their properties (src,dst ports / icmp type,code)
type EndpointsTrafficSet struct {
	props ds.TripleSet[*IPBlock, *IPBlock, *TransportSet]
}

// NewEndpointsTrafficSet returns an empty EndpointsTrafficSet
func NewEndpointsTrafficSet() *EndpointsTrafficSet {
	return &EndpointsTrafficSet{props: ds.NewLeftTripleSet[*IPBlock, *IPBlock, *TransportSet]()}
}

// Equal returns true is this EndpointsTrafficSet captures the exact same set of connections as `other` does.
func (c *EndpointsTrafficSet) Equal(other *EndpointsTrafficSet) bool {
	return c.props.Equal(other.props)
}

// Copy returns new EndpointsTrafficSet object with same set of connections as current one
func (c *EndpointsTrafficSet) Copy() *EndpointsTrafficSet {
	return &EndpointsTrafficSet{
		props: c.props.Copy(),
	}
}

// Intersect returns a EndpointsTrafficSet object with connection tuples that result from intersection of
// this and `other` sets
func (c *EndpointsTrafficSet) Intersect(other *EndpointsTrafficSet) *EndpointsTrafficSet {
	return &EndpointsTrafficSet{props: c.props.Intersect(other.props)}
}

// IsEmpty returns true of the EndpointsTrafficSet is empty
func (c *EndpointsTrafficSet) IsEmpty() bool {
	return c.props.IsEmpty()
}

// Union returns a EndpointsTrafficSet object with connection tuples that result from union of
// this and `other` sets
func (c *EndpointsTrafficSet) Union(other *EndpointsTrafficSet) *EndpointsTrafficSet {
	if other.IsEmpty() {
		return c.Copy()
	}
	if c.IsEmpty() {
		return other.Copy()
	}
	return &EndpointsTrafficSet{
		props: c.props.Union(other.props),
	}
}

// Subtract returns a EndpointsTrafficSet object with connection tuples that result from subtraction of
// `other` from this set
func (c *EndpointsTrafficSet) Subtract(other *EndpointsTrafficSet) *EndpointsTrafficSet {
	if other.IsEmpty() {
		return c.Copy()
	}
	return &EndpointsTrafficSet{props: c.props.Subtract(other.props)}
}

// IsSubset returns true if c is subset of other
func (c *EndpointsTrafficSet) IsSubset(other *EndpointsTrafficSet) bool {
	return c.props.IsSubset(other.props)
}

// EndpointsTrafficSetFrom returns a new EndpointsTrafficSet object from input src, dst IP-ranges sets ands
// TransportSet connections
func EndpointsTrafficSetFrom(src, dst *IPBlock, conn *TransportSet) *EndpointsTrafficSet {
	return &EndpointsTrafficSet{props: ds.CartesianLeftTriple(src, dst, conn)}
}

func (c *EndpointsTrafficSet) Partitions() []ds.Triple[*IPBlock, *IPBlock, *TransportSet] {
	return c.props.Partitions()
}

func cubeStr(c ds.Triple[*IPBlock, *IPBlock, *TransportSet]) string {
	return fmt.Sprintf("src: %s, dst: %s, conns: %s", c.S1.String(), c.S2.String(), c.S3.String())
}

func (c *EndpointsTrafficSet) String() string {
	cubes := c.Partitions()
	var resStrings = make([]string, len(cubes))
	for i, cube := range cubes {
		resStrings[i] = cubeStr(cube)
	}
	sort.Strings(resStrings)
	return strings.Join(resStrings, comma)
}
