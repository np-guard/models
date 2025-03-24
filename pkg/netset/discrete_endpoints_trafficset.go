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
	"github.com/np-guard/models/pkg/interval"
)

// DiscreteEndpointsTrafficSet captures a set of traffic attributes for tuples of (source endpoints, desination endpoints, TransportSet),
// where TransportSet is a set of TCP/UPD/ICMP with their properties (src,dst ports / icmp type,code)
// and source/destination endpoints are from a discrete set represented by integer IDs (could be mapped to VMs UIDs / Pod UIDs, etc.. )
type DiscreteEndpointsTrafficSet struct {
	props ds.TripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *TransportSet]
}

// EmptyDiscreteEndpointsTrafficSet returns an empty DiscreteEndpointsTrafficSet
func EmptyDiscreteEndpointsTrafficSet() *DiscreteEndpointsTrafficSet {
	return &DiscreteEndpointsTrafficSet{props: ds.NewLeftTripleSet[*interval.CanonicalSet, *interval.CanonicalSet, *TransportSet]()}
}

// Equal returns true is this DiscreteEndpointsTrafficSet captures the exact same set of connections as `other` does.
func (c *DiscreteEndpointsTrafficSet) Equal(other *DiscreteEndpointsTrafficSet) bool {
	return c.props.Equal(other.props)
}

// Copy returns new DiscreteEndpointsTrafficSet object with same set of connections as current one
func (c *DiscreteEndpointsTrafficSet) Copy() *DiscreteEndpointsTrafficSet {
	return &DiscreteEndpointsTrafficSet{
		props: c.props.Copy(),
	}
}

// Intersect returns a DiscreteEndpointsTrafficSet object with connection tuples that result from intersection of
// this and `other` sets
func (c *DiscreteEndpointsTrafficSet) Intersect(other *DiscreteEndpointsTrafficSet) *DiscreteEndpointsTrafficSet {
	return &DiscreteEndpointsTrafficSet{props: c.props.Intersect(other.props)}
}

// IsEmpty returns true of the DiscreteEndpointsTrafficSet is empty
func (c *DiscreteEndpointsTrafficSet) IsEmpty() bool {
	return c.props.IsEmpty()
}

// Union returns a DiscreteEndpointsTrafficSet object with connection tuples that result from union of
// this and `other` sets
func (c *DiscreteEndpointsTrafficSet) Union(other *DiscreteEndpointsTrafficSet) *DiscreteEndpointsTrafficSet {
	if other.IsEmpty() {
		return c.Copy()
	}
	if c.IsEmpty() {
		return other.Copy()
	}
	return &DiscreteEndpointsTrafficSet{
		props: c.props.Union(other.props),
	}
}

// Subtract returns a DiscreteEndpointsTrafficSet object with connection tuples that result from subtraction of
// `other` from this set
func (c *DiscreteEndpointsTrafficSet) Subtract(other *DiscreteEndpointsTrafficSet) *DiscreteEndpointsTrafficSet {
	if other.IsEmpty() {
		return c.Copy()
	}
	return &DiscreteEndpointsTrafficSet{props: c.props.Subtract(other.props)}
}

// IsSubset returns true if c is subset of other
func (c *DiscreteEndpointsTrafficSet) IsSubset(other *DiscreteEndpointsTrafficSet) bool {
	return c.props.IsSubset(other.props)
}

// NewDiscreteEndpointsTrafficSet returns a new DiscreteEndpointsTrafficSet object from input src, dst endpoint sets ands
// TransportSet connections
func NewDiscreteEndpointsTrafficSet(src, dst *interval.CanonicalSet, conn *TransportSet) *DiscreteEndpointsTrafficSet {
	return &DiscreteEndpointsTrafficSet{props: ds.CartesianLeftTriple(src, dst, conn)}
}

func (c *DiscreteEndpointsTrafficSet) Partitions() []ds.Triple[*interval.CanonicalSet, *interval.CanonicalSet, *TransportSet] {
	return c.props.Partitions()
}

func (c *DiscreteEndpointsTrafficSet) String() string {
	if c.IsEmpty() {
		return "<empty>"
	}
	cubes := c.Partitions()
	var resStrings = make([]string, len(cubes))
	for i, c := range cubes {
		resStrings[i] = fmt.Sprintf("src: %s, dst: %s, conns: %s", c.S1.String(), c.S2.String(), c.S3.String())
	}
	sort.Strings(resStrings)
	return strings.Join(resStrings, comma)
}
