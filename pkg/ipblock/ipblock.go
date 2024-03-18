package ipblock

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/np-guard/models/pkg/interval"
)

const (
	// CidrAll represents the CIDR for all addresses "0.0.0.0/0"
	CidrAll = "0.0.0.0/0"

	// internal const  below
	ipByte         = 0xff
	ipShift0       = 24
	ipShift1       = 16
	ipShift2       = 8
	ipBase         = 10
	ipMask         = 0xffffffff
	maxIPv4Bits    = 32
	cidrSeparator  = "/"
	bitSize64      = 64
	commaSeparator = ", "
	dash           = "-"
)

// IPBlock captures a set of IP ranges
type IPBlock struct {
	ipRange *interval.CanonicalSet
}

// New returns a new IPBlock object
func New() *IPBlock {
	return &IPBlock{
		ipRange: interval.NewCanonicalIntervalSet(),
	}
}

// ToIPRanges returns a string of the ip ranges in the current IPBlock object
func (b *IPBlock) ToIPRanges() string {
	return strings.Join(b.toIPRangesList(), commaSeparator)
}

// toIPRange returns a string of the ip range of a single interval
func toIPRange(i interval.Interval) string {
	startIP := intToIP4(i.Start)
	endIP := intToIP4(i.End)
	return rangeIPstr(startIP, endIP)
}

// toIPRangesList: returns a list of the ip-ranges strings in the current IPBlock object
func (b *IPBlock) toIPRangesList() []string {
	IPRanges := make([]string, len(b.ipRange.IntervalSet))
	for index := range b.ipRange.IntervalSet {
		IPRanges[index] = toIPRange(b.ipRange.IntervalSet[index])
	}
	return IPRanges
}

// ContainedIn checks if this IP block is contained within another IP block.
func (b *IPBlock) ContainedIn(other *IPBlock) bool {
	if b == other {
		return true
	}
	return b.ipRange.ContainedIn(other.ipRange)
}

// Intersect returns a new IPBlock from intersection of this IPBlock with input IPBlock
func (b *IPBlock) Intersect(c *IPBlock) *IPBlock {
	if b == c {
		return b.Copy()
	}
	return &IPBlock{
		ipRange: b.ipRange.Intersect(c.ipRange),
	}
}

// Equal returns true if this IPBlock equals the input IPBlock
func (b *IPBlock) Equal(c *IPBlock) bool {
	if b == c {
		return true
	}
	return b.ipRange.Equal(c.ipRange)
}

// Subtract returns a new IPBlock from subtraction of input IPBlock from this IPBlock
func (b *IPBlock) Subtract(c *IPBlock) *IPBlock {
	if b == c {
		return New()
	}
	return &IPBlock{
		ipRange: b.ipRange.Subtract(c.ipRange),
	}
}

// Union returns a new IPBlock from union of input IPBlock with this IPBlock
func (b *IPBlock) Union(c *IPBlock) *IPBlock {
	if b == c {
		return b.Copy()
	}
	return &IPBlock{
		ipRange: b.ipRange.Union(c.ipRange),
	}
}

// Empty returns true if this IPBlock is empty
func (b *IPBlock) Empty() bool {
	return b.ipRange.IsEmpty()
}

func rangeIPstr(start, end string) string {
	return fmt.Sprintf("%v-%v", start, end)
}

// Copy returns a new copy of IPBlock object
func (b *IPBlock) Copy() *IPBlock {
	return &IPBlock{ipRange: b.ipRange.Copy()}
}

func (b *IPBlock) ipCount() int {
	res := 0
	for _, r := range b.ipRange.IntervalSet {
		res += int(r.End) - int(r.Start) + 1
	}
	return res
}

// Split returns a set of IpBlock objects, each with a single range of ips
func (b *IPBlock) Split() []*IPBlock {
	res := make([]*IPBlock, len(b.ipRange.IntervalSet))
	for index, ipr := range b.ipRange.IntervalSet {
		res[index] = &IPBlock{
			ipRange: interval.FromInterval(ipr.Start, ipr.End),
		}
	}
	return res
}

// intToIP4 returns a string of an ip address from an input integer ip value
func intToIP4(ipInt int64) string {
	var d [4]byte
	binary.BigEndian.PutUint32(d[:], uint32(ipInt))
	return net.IPv4(d[0], d[1], d[2], d[3]).String()
}

// DisjointIPBlocks returns an IPBlock of disjoint ip ranges from 2 input IPBlock objects
func DisjointIPBlocks(set1, set2 []*IPBlock) []*IPBlock {
	ipbList := []*IPBlock{}
	for _, ipb := range set1 {
		ipbList = append(ipbList, ipb.Copy())
	}
	for _, ipb := range set2 {
		ipbList = append(ipbList, ipb.Copy())
	}
	// sort ipbList by ip_count per ipblock
	sort.Slice(ipbList, func(i, j int) bool {
		return ipbList[i].ipCount() < ipbList[j].ipCount()
	})
	// making sure the resulting list does not contain overlapping ipBlocks
	blocksWithNoOverlaps := []*IPBlock{}
	for _, ipb := range ipbList {
		blocksWithNoOverlaps = addIntervalToList(ipb, blocksWithNoOverlaps)
	}

	res := blocksWithNoOverlaps
	if len(res) == 0 {
		newAll := GetCidrAll()
		res = append(res, newAll)
	}
	return res
}

// addIntervalToList is used for computation of DisjointIPBlocks
func addIntervalToList(ipbNew *IPBlock, ipbList []*IPBlock) []*IPBlock {
	toAdd := []*IPBlock{}
	for idx, ipb := range ipbList {
		if !ipb.ipRange.Overlaps(ipbNew.ipRange) {
			continue
		}
		intersection := ipb.Intersect(ipbNew)
		ipbNew = ipbNew.Subtract(intersection)
		if !ipb.Equal(intersection) {
			toAdd = append(toAdd, intersection)
			ipbList[idx] = ipbList[idx].Subtract(intersection)
		}
		if ipbNew.Empty() {
			break
		}
	}
	ipbList = append(ipbList, ipbNew.Split()...)
	ipbList = append(ipbList, toAdd...)
	return ipbList
}

// FromCidr returns a new IPBlock object from input CIDR string
func FromCidr(cidr string) (*IPBlock, error) {
	span, err := cidrToInterval(cidr)
	if err != nil {
		return nil, err
	}
	return &IPBlock{
		ipRange: interval.FromInterval(span.Start, span.End),
	}, nil
}

// PairCIDRsToIPBlocks returns two IPBlock objects from two input CIDR strings
func PairCIDRsToIPBlocks(cidr1, cidr2 string) (ipb1, ipb2 *IPBlock, err error) {
	ipb1, err1 := FromCidr(cidr1)
	ipb2, err2 := FromCidr(cidr2)
	if err1 != nil || err2 != nil {
		return nil, nil, errors.Join(err1, err2)
	}
	return ipb1, ipb2, nil
}

// FromCidrOrAddress returns a new IPBlock object from input string of CIDR or IP address
func FromCidrOrAddress(s string) (*IPBlock, error) {
	if strings.Contains(s, cidrSeparator) {
		return FromCidr(s)
	}
	return FromIPAddress(s)
}

// FromCidrList returns IPBlock object from multiple CIDRs given as list of strings
func FromCidrList(cidrsList []string) (*IPBlock, error) {
	res := New()
	for _, cidr := range cidrsList {
		block, err := FromCidr(cidr)
		if err != nil {
			return nil, err
		}
		res = res.Union(block)
	}
	return res, nil
}

// Except creates a new IP block that excludes the specified CIDRs from the current IP block
func (b *IPBlock) Except(exceptions ...string) (*IPBlock, error) {
	holes := interval.NewCanonicalIntervalSet()
	for i := range exceptions {
		intervalHole, err := cidrToInterval(exceptions[i])
		if err != nil {
			return nil, err
		}
		holes.AddInterval(intervalHole)
	}
	return &IPBlock{ipRange: b.ipRange.Subtract(holes)}, nil
}

func ipv4AddressToCidr(ipAddress string) string {
	return ipAddress + "/32"
}

// FromIPAddress returns an IPBlock object from input IP address string
func FromIPAddress(ipAddress string) (*IPBlock, error) {
	return FromCidr(ipv4AddressToCidr(ipAddress))
}

func cidrToIPRange(cidr string) (start, end int64, err error) {
	// convert string to IPNet struct
	_, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		return
	}

	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	startNum := binary.BigEndian.Uint32(ipv4Net.IP)
	// find the final address
	endNum := (startNum & mask) | (mask ^ ipMask)
	start = int64(startNum)
	end = int64(endNum)
	return
}

func cidrToInterval(cidr string) (interval.Interval, error) {
	start, end, err := cidrToIPRange(cidr)
	if err != nil {
		return interval.Interval{}, err
	}
	return interval.Interval{Start: start, End: end}, nil
}

// ToCidrList returns a list of CIDR strings for this IPBlock object
func (b *IPBlock) ToCidrList() []string {
	cidrList := []string{}
	for _, interval := range b.ipRange.IntervalSet {
		cidrList = append(cidrList, intervalToCidrList(interval.Start, interval.End)...)
	}
	return cidrList
}

// ToCidrListString returns a string with all CIDRs within the IPBlock object
func (b *IPBlock) ToCidrListString() string {
	return strings.Join(b.ToCidrList(), commaSeparator)
}

// ListToPrint: returns a uniform to print list s.t. each element contains either a single cidr or an ip range
func (b *IPBlock) ListToPrint() []string {
	cidrsIPRangesList := []string{}
	for _, interval := range b.ipRange.IntervalSet {
		cidr := intervalToCidrList(interval.Start, interval.End)
		if len(cidr) == 1 {
			cidrsIPRangesList = append(cidrsIPRangesList, cidr[0])
		} else {
			cidrsIPRangesList = append(cidrsIPRangesList, toIPRange(interval))
		}
	}
	return cidrsIPRangesList
}

// ToIPAdressString returns the IP Address string for this IPBlock
func (b *IPBlock) ToIPAddressString() string {
	if b.ipRange.IsSingleNumber() {
		return intToIP4(b.ipRange.IntervalSet[0].Start)
	}
	return ""
}

func intervalToCidrList(ipStart, ipEnd int64) []string {
	start := ipStart
	end := ipEnd
	res := []string{}
	for end >= start {
		maxSize := maxIPv4Bits
		for maxSize > 0 {
			s := maxSize - 1
			mask := int64(math.Round(math.Pow(2, maxIPv4Bits) - math.Pow(2, float64(maxIPv4Bits)-float64(s))))
			maskBase := start & mask
			if maskBase != start {
				break
			}
			maxSize--
		}
		x := math.Log(float64(end)-float64(start)+1) / math.Log(2)
		maxDiff := byte(maxIPv4Bits - math.Floor(x))
		if maxSize < int(maxDiff) {
			maxSize = int(maxDiff)
		}
		ip := intToIP4(start)
		res = append(res, fmt.Sprintf("%s/%d", ip, maxSize))
		start += int64(math.Pow(2, maxIPv4Bits-float64(maxSize)))
	}
	return res
}

// FromIPRangeStr returns IPBlock object from input IP range string (example: "169.255.0.0-172.15.255.255")
func FromIPRangeStr(ipRangeStr string) (*IPBlock, error) {
	ipAddresses := strings.Split(ipRangeStr, dash)
	if len(ipAddresses) != 2 {
		return nil, errors.New("unexpected ipRange str")
	}
	var startIP, endIP *IPBlock
	var err error
	if startIP, err = FromIPAddress(ipAddresses[0]); err != nil {
		return nil, err
	}
	if endIP, err = FromIPAddress(ipAddresses[1]); err != nil {
		return nil, err
	}
	startIPNum := startIP.ipRange.IntervalSet[0].Start
	endIPNum := endIP.ipRange.IntervalSet[0].Start
	return &IPBlock{
		ipRange: interval.FromInterval(startIPNum, endIPNum),
	}, nil
}

// GetCidrAll returns IPBlock object of the entire range 0.0.0.0/0
func GetCidrAll() *IPBlock {
	res, _ := FromCidr(CidrAll)
	return res
}

// PrefixLength returns the cidr's prefix length, assuming the ipBlock is exactly one cidr.
// Prefix length specifies the number of bits in the IP address that are to be used as the subnet mask.
func (b *IPBlock) PrefixLength() (int64, error) {
	cidrs := b.ToCidrList()
	if len(cidrs) != 1 {
		return 0, errors.New("prefixLength err: ipBlock is not a single cidr")
	}
	cidrStr := cidrs[0]
	lenStr := strings.Split(cidrStr, cidrSeparator)[1]
	return strconv.ParseInt(lenStr, ipBase, bitSize64)
}
