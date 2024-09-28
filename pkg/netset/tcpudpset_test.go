package netset_test

import (
	"fmt"
	"testing"

	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"
	"github.com/stretchr/testify/require"
)

func TestTCPUDPSetBasicFunctionality(t *testing.T) {
	all := netset.AllTCPUDPSet()
	empty := netset.EmptyTCPorUDPSet()
	tcp80 := netset.NewTCPorUDPSet(netp.ProtocolStringTCP, netp.MinPort, netp.MaxPort, 80, 80)
	tcpudp81 := netset.NewTCPorUDPSet(netp.ProtocolStringTCP, netp.MinPort, netp.MaxPort, 81, 81).Union(
		netset.NewTCPorUDPSet(netp.ProtocolStringUDP, netp.MinPort, netp.MaxPort, 81, 81),
	)
	tcp80WithSrcPorts := netset.NewTCPorUDPSet(netp.ProtocolStringTCP, 5000, 5300, 80, 80)
	allButTCP80 := all.Subtract(tcp80)

	// TODO: this obj creation should fail?
	// invalidObj := netset.NewTCPorUDPSet(netp.ProtocolStringTCP, 65539, 65539, 65539, 65539)

	fmt.Println(all)                   // TCP,UDP
	fmt.Println(empty)                 // "" (empty string)
	fmt.Println(tcp80)                 // TCP dst-ports: 80
	fmt.Println(allButTCP80)           // TCP dst-ports: 1-79,81-65535,UDP
	fmt.Println(tcpudp81)              // TCP,UDP dst-ports: 81
	fmt.Println(tcp80.Union(tcpudp81)) // TCP dst-ports: 80-81,UDP dst-ports: 81
	fmt.Println(tcpudp81.Union(tcp80)) // TCP dst-ports: 80-81,UDP dst-ports: 81
	fmt.Println(tcp80WithSrcPorts)     // TCP src-ports: 5000-5300 dst-ports: 80

	// IsAll, IsEmpty
	require.True(t, all.IsAll())
	require.False(t, all.IsEmpty())

	require.True(t, empty.IsEmpty())
	require.False(t, empty.IsAll())

	require.False(t, tcp80.IsAll())
	require.False(t, tcp80.IsEmpty())

	// Equal
	require.False(t, all.Equal(empty))
	require.False(t, empty.Equal(all))
	require.True(t, empty.Equal(empty))
	require.True(t, all.Equal(all))
	require.False(t, empty.Equal(tcp80))
	require.False(t, tcp80.Equal(empty))
	require.False(t, all.Equal(tcp80))
	require.False(t, tcp80.Equal(all))

	// IsSubset
	require.True(t, empty.IsSubset(all))
	require.False(t, all.IsSubset(empty))
	require.True(t, empty.IsSubset(empty))
	require.True(t, all.IsSubset(all))

	require.True(t, tcp80.IsSubset(all))
	require.False(t, tcp80.IsSubset(empty))
	require.True(t, tcp80.IsSubset(tcp80))

	// Subtract, Union, Intersect
	require.True(t, tcp80.Union(allButTCP80).Equal(all))
	require.True(t, tcp80.Intersect(allButTCP80).IsEmpty())

}
