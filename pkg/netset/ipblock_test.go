/*
Copyright 2023- IBM Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package netset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/netset"
)

func TestOps(t *testing.T) {
	ipb1, err := netset.IPBlockFromCidrOrAddress("1.2.3.0/24")
	require.Nil(t, err)
	require.NotNil(t, ipb1)
	ipb2, err := netset.IPBlockFromCidrOrAddress("1.2.3.4")
	require.Nil(t, err)
	require.NotNil(t, ipb2)
	require.True(t, ipb2.IsSubset(ipb1))
	require.False(t, ipb1.IsSubset(ipb2))

	minus := ipb1.Subtract(ipb2)
	require.Equal(t, "1.2.3.0-1.2.3.3, 1.2.3.5-1.2.3.255", minus.ToIPRanges())
	require.Equal(t, "1.2.3.0", minus.FirstIPAddress())

	minus2, err := netset.IPBlockFromCidr(ipb1.ToCidrListString())
	require.Nil(t, err)
	minus2, err = minus2.ExceptCidrs(ipb2.ToCidrListString())
	require.Nil(t, err)
	require.Equal(t, minus.ToCidrListString(), minus2.ToCidrListString())

	intersect := ipb1.Intersect(ipb2)
	require.Equal(t, intersect, ipb2)

	union := intersect.Union(minus)
	require.Equal(t, union, ipb1)

	intersect2 := minus.Intersect(intersect)
	require.True(t, intersect2.IsEmpty())

	ipb3, err := ipb2.NextIP() // ipb3 = 1.2.3.5
	ipb4, _ := netset.IPBlockFromCidrOrAddress("1.2.3.5")
	require.Nil(t, err)
	require.Equal(t, ipb3, ipb4)

	ipb5, err := ipb3.PreviousIP() // ipb5 = 1.2.3.4
	require.Nil(t, err)
	require.Equal(t, ipb2, ipb5)

	ipb6, err := ipb1.NextIP() // ipb6 = 1.2.4.0
	ipb7, _ := netset.IPBlockFromCidrOrAddress("1.2.4.0")
	require.Nil(t, err)
	require.Equal(t, ipb6, ipb7)

	require.False(t, ipb1.IsSingleIPAddress())
	require.True(t, ipb2.IsSingleIPAddress())

	ipb8, err := ipb7.PreviousIP() // ipb8 = 1.2.3.255
	require.Nil(t, err)
	require.Equal(t, ipb8, ipb1.LastIPAddressObject())

	ipb9, err := netset.IPBlockFromIPRange(ipb4, ipb6)
	require.Nil(t, err)
	// ipb9 = 1.2.3.5-1.2.4.0. Equal to the union of:
	// 1.2.3.5/32, 1.2.3.6/31, 1.2.3.8/29, 1.2.3.16/28,
	// 1.2.3.32/27, 1.2.3.64/26, 1.2.3.128/25, 1.2.4.0/32
	require.Len(t, ipb9.SplitToCidrs(), 8)

	t1, err := ipb9.TouchingIPRanges(ipb2)
	require.Nil(t, err)
	require.True(t, t1)

	t2, err := ipb9.TouchingIPRanges(ipb7)
	require.Nil(t, err)
	require.False(t, t2)

	require.Equal(t, ipb7, ipb7.FirstIPAddressObject())

	require.Equal(t, ipb5.Compare(ipb6), -1)
	require.Equal(t, ipb2.Compare(ipb1), 1)
	require.Equal(t, ipb3.Compare(ipb4), 0)
}

func TestConversions(t *testing.T) {
	ipRange := "172.0.10.0-195.8.5.14"
	ipb1, err := netset.IPBlockFromIPRangeStr(ipRange)
	require.Nil(t, err)
	require.Equal(t, ipRange, ipb1.ToIPRanges())
	require.Equal(t, "172.0.10.0", ipb1.FirstIPAddress())

	cidrs := ipb1.ToCidrList()
	require.Len(t, cidrs, 26)

	ipb2, err := netset.IPBlockFromCidrList(cidrs)
	require.Nil(t, err)
	require.Equal(t, ipb1.ToCidrListString(), ipb2.ToCidrListString())

	toPrint := ipb1.ListToPrint()
	require.Len(t, toPrint, 1)
	require.Equal(t, ipRange, toPrint[0])

	require.Equal(t, "", ipb1.ToIPAddressString())

	_, err = ipb1.AsCidr()
	require.NotNil(t, err)

	cidr := "5.2.1.0/24"
	ipb3, _ := netset.IPBlockFromCidr(cidr)
	str, err := ipb3.AsCidr()
	require.Nil(t, err)
	require.Equal(t, str, cidr)
}

func TestDisjointIPBlocks(t *testing.T) {
	allIPs := netset.GetCidrAll()
	ipb, err := netset.IPBlockFromCidrList([]string{"1.2.3.4/32", "172.0.0.0/8"})
	require.Nil(t, err)

	disjointBlocks := netset.DisjointIPBlocks([]*netset.IPBlock{allIPs}, []*netset.IPBlock{ipb})
	require.Len(t, disjointBlocks, 5)
	require.Equal(t, "1.2.3.4", disjointBlocks[0].ToIPAddressString()) // list is sorted by ip-block size

	ipb2, err := netset.IPBlockFromCidrList([]string{"1.2.3.0/30"})
	require.Nil(t, err)
	ipb3, err := netset.IPBlockFromIPRangeStr("1.2.2.255-1.2.3.1")
	require.Nil(t, err)
	disjointBlocks = netset.DisjointIPBlocks([]*netset.IPBlock{ipb2}, []*netset.IPBlock{ipb3})
	require.Len(t, disjointBlocks, 3)
	require.Equal(t, "1.2.2.255", disjointBlocks[0].ToIPAddressString())
	require.Equal(t, "1.2.3.2/31", disjointBlocks[1].ToCidrListString())
	require.Equal(t, "1.2.3.0/31", disjointBlocks[2].ToCidrListString())
}

func TestPairCIDRsToIPBlocks(t *testing.T) {
	first, second, err := netset.PairCIDRsToIPBlocks("5.6.7.8/24", "4.9.2.1/32")
	require.Nil(t, err)
	require.Equal(t, "5.6.7.0/24", first.ListToPrint()[0])
	require.Equal(t, "4.9.2.1/32", second.ListToPrint()[0])

	intersect := first.Intersect(second)
	require.Equal(t, "", intersect.ToIPAddressString())
	require.Empty(t, intersect.ListToPrint())
	require.Empty(t, intersect.ToCidrListString())
}

func TestPrefixLength(t *testing.T) {
	ipb, err := netset.IPBlockFromCidrOrAddress("42.5.2.8/20")
	require.Nil(t, err)
	prefLen, err := ipb.PrefixLength()
	require.Nil(t, err)
	require.Equal(t, int64(20), prefLen)

	ipb, err = netset.IPBlockFromCidrOrAddress("42.5.2.8")
	require.Nil(t, err)
	prefLen, err = ipb.PrefixLength()
	require.Nil(t, err)
	require.Equal(t, int64(32), prefLen)

	ipb, err = netset.IPBlockFromCidrList([]string{"1.2.3.4/32", "172.0.0.0/8"})
	require.Nil(t, err)
	_, err = ipb.PrefixLength()
	require.NotNil(t, err)
}

func TestString(t *testing.T) {
	ipb, err := netset.IPBlockFromCidrOrAddress("42.5.2.8")
	require.Nil(t, err)
	require.Equal(t, "42.5.2.8", ipb.String())

	ipb, err = netset.IPBlockFromCidr("42.5.0.0/20")
	require.Nil(t, err)
	require.Equal(t, "42.5.0.0/20", ipb.String())

	ipb, err = netset.IPBlockFromCidrList([]string{"1.2.3.4/32", "172.0.0.0/8"})
	require.Nil(t, err)
	require.Equal(t, "1.2.3.4/32, 172.0.0.0/8", ipb.String())
}

func TestBadPath(t *testing.T) {
	_, err := netset.IPBlockFromCidr("not-a-cidr")
	require.NotNil(t, err)

	_, err = netset.IPBlockFromCidr("2.5.7.9/24")
	require.Nil(t, err)

	_, err = netset.NewIPBlock().ExceptCidrs("5.6.7.8/20", "not-a-cidr")
	require.NotNil(t, err)

	_, err = netset.IPBlockFromCidrList([]string{"1.2.3.4/20", "not-a-cidr"})
	require.NotNil(t, err)

	_, err = netset.IPBlockFromCidrList([]string{"1.2.3.4/20", "1.2.3.4/40"})
	require.NotNil(t, err)

	_, err = netset.IPBlockFromIPRangeStr("1.2.3.4")
	require.NotNil(t, err)

	_, err = netset.IPBlockFromIPRangeStr("prefix-1.2.3.4")
	require.NotNil(t, err)

	_, err = netset.IPBlockFromIPRangeStr("1.2.3.290-1.2.3.4")
	require.NotNil(t, err)

	_, err = netset.IPBlockFromIPRangeStr("1.2.3.4-suffix")
	require.NotNil(t, err)

	_, err = netset.IPBlockFromIPRangeStr("1.2.3.4-2.5.6.7/20")
	require.NotNil(t, err)

	_, _, err = netset.PairCIDRsToIPBlocks("1.2.3.4/40", "1.2.3.5/24")
	require.NotNil(t, err)

	_, _, err = netset.PairCIDRsToIPBlocks("1.2.3.4/20", "not-a-cidr")
	require.NotNil(t, err)
}
