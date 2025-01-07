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
	"github.com/np-guard/models/pkg/netp"
)

type TypeSet = interval.CanonicalSet
type CodeSet = interval.CanonicalSet

type ICMPSet struct {
	props ds.Product[*TypeSet, *CodeSet]
}

func (c *ICMPSet) Equal(other *ICMPSet) bool {
	return c.props.Equal(other.props)
}

func (c *ICMPSet) Hash() int {
	return c.props.Hash()
}

func (c *ICMPSet) Copy() *ICMPSet {
	return &ICMPSet{props: c.props.Copy()}
}

func (c *ICMPSet) Intersect(other *ICMPSet) *ICMPSet {
	return &ICMPSet{props: c.props.Intersect(other.props)}
}

func (c *ICMPSet) Partitions() []ds.Pair[*TypeSet, *CodeSet] {
	return c.props.Partitions()
}

func (c *ICMPSet) IsEmpty() bool {
	return c.props.IsEmpty()
}

func (c *ICMPSet) Union(other *ICMPSet) *ICMPSet {
	return &ICMPSet{props: c.props.Union(other.props)}
}

func (c *ICMPSet) Size() int {
	return c.props.Size()
}

// Subtract returns the subtraction of the other from c
func (c *ICMPSet) Subtract(other *ICMPSet) *ICMPSet {
	return &ICMPSet{props: c.props.Subtract(other.props)}
}

// IsSubset returns true if c is subset of other
func (c *ICMPSet) IsSubset(other *ICMPSet) bool {
	return c.props.IsSubset(other.props)
}

// icmpPropsPathLeft creates a new ICMPSet, implemented using CartesianPairLeft.
func icmpPropsPathLeft(typesSet *TypeSet, codeSet *CodeSet) *ICMPSet {
	return &ICMPSet{props: ds.CartesianPairLeft(typesSet, codeSet)}
}

func NewICMPSet(minType, maxType, minCode, maxCode int64) *ICMPSet {
	return icmpPropsPathLeft(
		interval.New(minType, maxType).ToSet(),
		interval.New(minCode, maxCode).ToSet(),
	)
}

func ICMPSetFromICMP(icmp netp.ICMP) *ICMPSet {
	if icmp.TypeCode == nil {
		return AllICMPSet()
	}
	icmpType := int64(icmp.TypeCode.Type)
	if icmp.TypeCode.Code == nil {
		return NewICMPSet(icmpType, icmpType, int64(netp.MinICMPCode), int64(netp.MaxICMPCode))
	}
	icmpCode := int64(*icmp.TypeCode.Code)
	return NewICMPSet(icmpType, icmpType, icmpCode, icmpCode)
}

func EmptyICMPSet() *ICMPSet {
	return &ICMPSet{props: ds.NewProductLeft[*TypeSet, *CodeSet]()}
}

func AllICMPSet() *ICMPSet {
	return icmpPropsPathLeft(
		AllICMPTypes(),
		AllICMPCodes(),
	)
}

func AllICMPCodes() *CodeSet {
	return interval.New(int64(netp.MinICMPCode), int64(netp.MaxICMPCode)).ToSet()
}

func AllICMPTypes() *TypeSet {
	return interval.New(int64(netp.MinICMPType), int64(netp.MaxICMPType)).ToSet()
}

var allICMP = AllICMPSet()

func (c *ICMPSet) IsAll() bool {
	return c.Equal(allICMP)
}

func getICMPCubeStr(cube ds.Pair[*TypeSet, *CodeSet]) string {
	if cube.Right.Equal(AllICMPCodes()) {
		return fmt.Sprintf("ICMP type: %s", cube.Left.String())
	}
	return fmt.Sprintf("ICMP type: %s code: %s", cube.Left.String(), cube.Right.String())
}

func (c *ICMPSet) String() string {
	if c.IsAll() {
		return string(netp.ProtocolStringICMP)
	}
	cubes := c.Partitions()
	var resStrings = make([]string, len(cubes))
	for i, cube := range cubes {
		resStrings[i] = getICMPCubeStr(cube)
	}
	sort.Strings(resStrings)
	return strings.Join(resStrings, " | ")
}
