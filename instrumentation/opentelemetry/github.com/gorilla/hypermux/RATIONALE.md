# Rationale

## Why create an own middleware instead of extending the one in gorila/mux/otelmux?

**gorilla/mux/otelmux** [requires you to pass a service name](https://github.com/open-telemetry/opentelemetry-go-contrib/blob/f284e28/instrumentation/github.com/gorilla/mux/otelmux/mux.go#L40) requiring to have access to the configuration from outside of the SDK. This is not only unnecessary but also requires us to leak loaded config into instrumentations. Ultimately, our wrapper handler allows us to inherit all the features and configuration without the need of leaking the config to instrumentation.
