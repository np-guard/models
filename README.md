# models
A collection of Golang packages with models for cartesian products and network resources

## Packages and Types
* **ds** - A collection of generic data structures.
  * Interfaces:
    * `Sized` (IsEmpty, Size)
    * `Comparable` (Equal, Copy)
    * `Hashable` (Comparable, Hash)
    * `Set` (Hashable, Sized, IsSubset, Union, Intersect, Substract)
    * `Product[A, B]` - A x B (Partitions, NumPartitions, Left and Right projections, Swap)
    * `TripleSet[S1, S2, S3]` - S1 x S2 x S3; associativity-agnostic (Partitions)
  * Concrete types:
    * `Pair` - A simple generic pair
    * `Triple` - A simple generic triple
    * `HashMap` - A generic map for mapping any Hashable key to any Comparable.
    * `HashSet` - A generic `Set` for storing any Hashable.
    * `MultiMap` - A map for mapping any Hashable key to a set of Hashable values.
    * `ProductLeft` - A `Product` of two sets, implemented using a map where each key-values pair represents the cartesian product of the two sets.
    * `LeftTripleSet`, `RightTripleSet`, `OuterTripleSet` - `TripleSet` implementations.
    * `DisjointSum` - A sum type for two tagged sets.
* **interval** - Interval-related data structures.
    * `Interval` - A simple interval data structure.
    * `IntervalSet` - A set of numbers, implements using intervals.
* **netp** - Various structs and functions representing and handling common network protocols (TCP, UDP, ICMP).
  * `ICMP` - describing type and code values for ICMP packets.
  * `TCPUDP` - describing port and protocol values for TCP and UDP packets.
  * `Protocol` - an interface for protocol values.
  * `AnyProtocol` - a protocol value that matches any protocol.
* **netset** - Sets of network-related tuples: IP addresses x ports x protocols, etc.
  * `PortSet` - A set of ports. Implemented using an IntervalSet.
  * `ProtocolSet` - Whether the protocol is TCP or UDP. Implemented using IntervalSet.
  * `TCPUDPSet` - `TripleSet[*ProtocolSet, *PortSet, *PortSet]`.
  * `RFCICMPSet` - accurately tracking set of ICMP types and code pairs. Implemented using a bitset.
  * `TypeSet` - ICMP types set. Implemented using an IntervalSet.
  * `CodeSet` ICMP codes set. Implemented using an IntervalSet.
  * `ICMPSet` - ICMP types and code pairs, implemented as `Product[*TypeSet, *CodeSet]`.
  * `TransportSet` - either ICMPSet or TCPUDP set. Implemented as `Disjoint[*TCPUDPSet, *ICMPSet]`.
  * `IPBlock` - A set of IP addresses. Implemented using IntervalSet.
  * `EndpointsTrafficSet` - `TripleSet[*IPBlock, *IPBlock, *TransportSet]`.
* **spec** - A collection of structs for defining required connectivity. Automatically generated from a JSON schema (see below).

## Code generation
`spec_schema.json` is the JSON schema for the input to VPC-synthesis. The data model in `pkg/spec` is auto-generated from this file using the below procedure.

Install [atombender/go-jsonschema](https://github.com/atombender/go-jsonschema)
(important: **not** [xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema))

```commandline
go get github.com/atombender/go-jsonschema/...
go install github.com/atombender/go-jsonschema@latest
```

Then run

```commandline
make generate
```

The result is written into [pkg/spec/data_model.go](pkg/spec/data_model.go).
