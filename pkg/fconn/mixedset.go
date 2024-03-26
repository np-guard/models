// Copyright 2020- IBM Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package fconn

import (
	"fmt"

	"github.com/np-guard/models/pkg/interval"
)

type ICMPSet = interval.CanonicalSet

type MixedSet struct {
	transport *TCPUDPSet
	icmp      *ICMPSet
}

func (c *MixedSet) Equal(other *MixedSet) bool {
	return c.transport.Equal(other.transport) && c.icmp.Equal(other.icmp)
}

func (c *MixedSet) Hash() int {
	return c.transport.Hash() ^ c.icmp.Hash()
}
func (c *MixedSet) Copy() *MixedSet {
	return &MixedSet{
		transport: c.transport.Copy(),
		icmp:      c.icmp.Copy(),
	}
}

func (c *MixedSet) Intersect(other *MixedSet) *MixedSet {
	return &MixedSet{
		transport: c.transport.Intersect(other.transport),
		icmp:      c.icmp.Intersect(other.icmp),
	}
}

func (c *MixedSet) IsEmpty() bool {
	return c.transport.IsEmpty() && c.icmp.IsEmpty()
}

func (c *MixedSet) Union(other *MixedSet) *MixedSet {
	return &MixedSet{
		transport: c.transport.Union(other.transport),
		icmp:      c.icmp.Union(other.icmp),
	}
}

func (c *MixedSet) Subtract(other *MixedSet) *MixedSet {
	return &MixedSet{
		transport: c.transport.Subtract(other.transport),
		icmp:      c.icmp.Subtract(other.icmp),
	}
}

// ContainedIn returns true if c is subset of other
func (c *MixedSet) ContainedIn(other *MixedSet) bool {
	return c.transport.ContainedIn(other.transport) && c.icmp.ContainedIn(other.icmp)
}

// String returns a string representation of a MixedSet object
func (c *MixedSet) String() string {
	return fmt.Sprintf("{%s|%s}", c.transport, c.icmp)
}
