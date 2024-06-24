# models
A collection of Golang packages with models for connectivity and network resources

## Packages
* **interval** - A canonical representation of a set of intervals defined over the integers
* **ipblock** - A canonical representation of a set of IP ranges. Currently limited to IPv4
* **hypercube** - A canonical representation of a set of n-dimensional hypercubes. All dimensions are defined over the integers.
* **netp** - Various structs for representing and handling common network protocols (TCP, UDP, ICMP)
* **connection** - A canonical representation of a set of connections. E.g., for representing all protocols/ports/codes permitted by a given firewall, given a specific source and destination.
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
