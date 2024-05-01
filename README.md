# models
Models for connectivity and network resources
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
