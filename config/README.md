# Agent Config Go

Agent config holds all the configuration settings for the Hypertrace Go Agent.

## Getting Started

```go
// loads the config from the config file tbd.json and env vars
cfg := config.Load()

// overrides statically the service name
cfg.ServiceName = config.String("myservice")
cfg.DataCapture.HTTPHeaders.Request = config.Bool(true)
```

Values can also be overriden by the environment variables, e.g. `HT_DATA_CAPTURE_HTTP_HEADERS_RESPONSE=false`.

The location for the config file can also be overriden by passing the path in `HT_CONFIG_FILE` environment variable.

## Default Values

All default values are defined in the [defaults.go](./defaults.go), everything else will be default to zero values.

## Generating config structs

When changing the proto definition, the config structs must be regenerated to use the new configuration fields. This can be done by:

```bash
cd ../
make generate-config
```
