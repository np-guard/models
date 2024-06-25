/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package netset

import (
	"github.com/np-guard/models/pkg/ds"
)

type ConnectionSet struct {
	props ds.TripleSet[*IPBlock, *IPBlock, *TransportSet]
}

func NewConnectionSet() *ConnectionSet {
	return &ConnectionSet{props: ds.NewLeftTripleSet[*IPBlock, *IPBlock, *TransportSet]()}
}

func (c *ConnectionSet) Equal(other *ConnectionSet) bool {
	return c.props.Equal(other.props)
}

func (c *ConnectionSet) Copy() *ConnectionSet {
	return &ConnectionSet{
		props: c.props.Copy(),
	}
}

func (c *ConnectionSet) Intersect(other *ConnectionSet) *ConnectionSet {
	return &ConnectionSet{props: c.props.Intersect(other.props)}
}

func (c *ConnectionSet) IsEmpty() bool {
	return c.props.IsEmpty()
}

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
