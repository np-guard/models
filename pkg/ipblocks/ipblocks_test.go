package ipblocks_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/np-guard/models/pkg/ipblocks"
)

func TestOps(t *testing.T) {
	ipb1 := ipblocks.NewIPBlockFromCidrOrAddress("1.2.3.0/24")
	require.NotNil(t, ipb1)
	ipb2 := ipblocks.NewIPBlockFromCidrOrAddress("1.2.3.4")
	require.NotNil(t, ipb2)
	require.True(t, ipb2.IsIPAddress("1.2.3.4"))
	require.True(t, ipb2.ContainedIn(ipb1))
	require.False(t, ipb1.ContainedIn(ipb2))

	minus := ipb1.Subtract(ipb2)
	minusRanges := minus.ToIPRangesList()
	require.Len(t, minusRanges, 2)
	require.Equal(t, "1.2.3.0-1.2.3.3", minusRanges[0])
	require.Equal(t, "1.2.3.5-1.2.3.255", minusRanges[1])

	minus2, err := ipblocks.NewIPBlock(ipb1.ToCidrListString(), []string{ipb2.ToCidrListString()})
	require.Nil(t, err)
	require.Equal(t, minus.ToCidrListString(), minus2.ToCidrListString())

	intersect := ipb1.Intersection(ipb2)
	require.True(t, intersect.Equal(ipb2))

	union := intersect.Union(minus)
	require.True(t, union.Equal(ipb1))

	intersect2 := minus.Intersection(intersect)
	require.True(t, intersect2.Empty())
}

func TestConversions(t *testing.T) {
	iprange := "172.0.10.0-195.8.5.14"
	ipb1, err := ipblocks.IPBlockFromIPRangeStr(iprange)
	require.Nil(t, err)
	require.Equal(t, iprange, ipb1.ToIPRanges())

	cidrs := ipb1.ToCidrList()
	require.Len(t, cidrs, 26)

	ipb2 := ipblocks.NewIPBlockFromCidrList(cidrs)
	require.Equal(t, ipb1.ToCidrListString(), ipb2.ToCidrListString())

	toprint := ipb1.ListToPrint()
	require.Len(t, toprint, 1)
	require.Equal(t, iprange, toprint[0])

	require.Equal(t, "", ipb1.ToIPAddress())
}

func TestDisjointIPBlocks(t *testing.T) {
	allIPs := ipblocks.GetCidrAll()
	ipb := ipblocks.NewIPBlockFromCidrList([]string{"1.2.3.4/32", "172.0.0.0/8"})

	disjointBlocks := ipblocks.DisjointIPBlocks([]*ipblocks.IPBlock{allIPs}, []*ipblocks.IPBlock{ipb})
	require.Len(t, disjointBlocks, 5)
	require.Equal(t, "1.2.3.4", disjointBlocks[0].ToIPAddress()) // list is sorted by ip-block size
}

func TestIsAddressInSubnet(t *testing.T) {
	res, err := ipblocks.IsAddressInSubnet("1.2.3.4", "1.0.0.0/8")
	require.Nil(t, err)
	require.True(t, res)

	res, err = ipblocks.IsAddressInSubnet("1.2.3.4", "1.0.0.0/16")
	require.Nil(t, err)
	require.False(t, res)

	_, err = ipblocks.IsAddressInSubnet("1.2.3.4/30", "1.0.0.0/16")
	require.NotNil(t, err)
}

func TestPrefixLength(t *testing.T) {
	ipb := ipblocks.NewIPBlockFromCidrOrAddress("42.5.2.8/20")
	prefLen, err := ipb.PrefixLength()
	require.Nil(t, err)
	require.Equal(t, int64(20), prefLen)

	ipb = ipblocks.NewIPBlockFromCidrOrAddress("42.5.2.8")
	prefLen, err = ipb.PrefixLength()
	require.Nil(t, err)
	require.Equal(t, int64(32), prefLen)

	ipb = ipblocks.NewIPBlockFromCidrList([]string{"1.2.3.4/32", "172.0.0.0/8"})
	_, err = ipb.PrefixLength()
	require.NotNil(t, err)
}
