package ipblock_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/ipblock"
)

func TestOps(t *testing.T) {
	ipb1, err := ipblock.FromCidrOrAddress("1.2.3.0/24")
	require.Nil(t, err)
	require.NotNil(t, ipb1)
	ipb2, err := ipblock.FromCidrOrAddress("1.2.3.4")
	require.Nil(t, err)
	require.NotNil(t, ipb2)
	require.True(t, ipb2.ContainedIn(ipb1))
	require.False(t, ipb1.ContainedIn(ipb2))

	minus := ipb1.Subtract(ipb2)
	require.Equal(t, "1.2.3.0-1.2.3.3, 1.2.3.5-1.2.3.255", minus.ToIPRanges())

	minus2, err := ipblock.FromCidr(ipb1.ToCidrListString())
	require.Nil(t, err)
	minus2, err = minus2.Except(ipb2.ToCidrListString())
	require.Nil(t, err)
	require.Equal(t, minus.ToCidrListString(), minus2.ToCidrListString())

	intersect := ipb1.Intersect(ipb2)
	require.True(t, intersect.Equal(ipb2))

	union := intersect.Union(minus)
	require.True(t, union.Equal(ipb1))

	intersect2 := minus.Intersect(intersect)
	require.True(t, intersect2.Empty())
}

func TestConversions(t *testing.T) {
	ipRange := "172.0.10.0-195.8.5.14"
	ipb1, err := ipblock.FromIPRangeStr(ipRange)
	require.Nil(t, err)
	require.Equal(t, ipRange, ipb1.ToIPRanges())

	cidrs := ipb1.ToCidrList()
	require.Len(t, cidrs, 26)

	ipb2, err := ipblock.FromCidrList(cidrs)
	require.Nil(t, err)
	require.Equal(t, ipb1.ToCidrListString(), ipb2.ToCidrListString())

	toPrint := ipb1.ListToPrint()
	require.Len(t, toPrint, 1)
	require.Equal(t, ipRange, toPrint[0])

	require.Equal(t, "", ipb1.ToIPAddressString())
}

func TestDisjointIPBlocks(t *testing.T) {
	allIPs := ipblock.GetCidrAll()
	ipb, err := ipblock.FromCidrList([]string{"1.2.3.4/32", "172.0.0.0/8"})
	require.Nil(t, err)

	disjointBlocks := ipblock.DisjointIPBlocks([]*ipblock.IPBlock{allIPs}, []*ipblock.IPBlock{ipb})
	require.Len(t, disjointBlocks, 5)
	require.Equal(t, "1.2.3.4", disjointBlocks[0].ToIPAddressString()) // list is sorted by ip-block size

	ipb2, err := ipblock.FromCidrList([]string{"1.2.3.0/30"})
	require.Nil(t, err)
	ipb3, err := ipblock.FromIPRangeStr("1.2.2.255-1.2.3.1")
	require.Nil(t, err)
	disjointBlocks = ipblock.DisjointIPBlocks([]*ipblock.IPBlock{ipb2}, []*ipblock.IPBlock{ipb3})
	require.Len(t, disjointBlocks, 3)
	require.Equal(t, "1.2.2.255", disjointBlocks[0].ToIPAddressString())
	require.Equal(t, "1.2.3.2/31", disjointBlocks[1].ToCidrListString())
	require.Equal(t, "1.2.3.0/31", disjointBlocks[2].ToCidrListString())
}

func TestPairCIDRsToIPBlocks(t *testing.T) {
	first, second, err := ipblock.PairCIDRsToIPBlocks("5.6.7.8/24", "4.9.2.1/32")
	require.Nil(t, err)
	require.Equal(t, "5.6.7.0/24", first.ListToPrint()[0])
	require.Equal(t, "4.9.2.1/32", second.ListToPrint()[0])

	intersect := first.Intersect(second)
	require.Equal(t, "", intersect.ToIPAddressString())
	require.Empty(t, intersect.ListToPrint())
	require.Empty(t, intersect.ToCidrListString())
}

func TestPrefixLength(t *testing.T) {
	ipb, err := ipblock.FromCidrOrAddress("42.5.2.8/20")
	require.Nil(t, err)
	prefLen, err := ipb.PrefixLength()
	require.Nil(t, err)
	require.Equal(t, int64(20), prefLen)

	ipb, err = ipblock.FromCidrOrAddress("42.5.2.8")
	require.Nil(t, err)
	prefLen, err = ipb.PrefixLength()
	require.Nil(t, err)
	require.Equal(t, int64(32), prefLen)

	ipb, err = ipblock.FromCidrList([]string{"1.2.3.4/32", "172.0.0.0/8"})
	require.Nil(t, err)
	_, err = ipb.PrefixLength()
	require.NotNil(t, err)
}

func TestBadPath(t *testing.T) {
	_, err := ipblock.FromCidr("not-a-cidr")
	require.NotNil(t, err)

	_, err = ipblock.FromCidr("2.5.7.9/24")
	require.Nil(t, err)

	_, err = ipblock.New().Except("5.6.7.8/20", "not-a-cidr")
	require.NotNil(t, err)

	_, err = ipblock.FromCidrList([]string{"1.2.3.4/20", "not-a-cidr"})
	require.NotNil(t, err)

	_, err = ipblock.FromCidrList([]string{"1.2.3.4/20", "1.2.3.4/40"})
	require.NotNil(t, err)

	_, err = ipblock.FromIPRangeStr("1.2.3.4")
	require.NotNil(t, err)

	_, err = ipblock.FromIPRangeStr("prefix-1.2.3.4")
	require.NotNil(t, err)

	_, err = ipblock.FromIPRangeStr("1.2.3.290-1.2.3.4")
	require.NotNil(t, err)

	_, err = ipblock.FromIPRangeStr("1.2.3.4-suffix")
	require.NotNil(t, err)

	_, err = ipblock.FromIPRangeStr("1.2.3.4-2.5.6.7/20")
	require.NotNil(t, err)

	_, _, err = ipblock.PairCIDRsToIPBlocks("1.2.3.4/40", "1.2.3.5/24")
	require.NotNil(t, err)

	_, _, err = ipblock.PairCIDRsToIPBlocks("1.2.3.4/20", "not-a-cidr")
	require.NotNil(t, err)
}
