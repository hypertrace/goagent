package identifier

import (
	uuid "github.com/satori/go.uuid"
	"go.opentelemetry.io/otel/attribute"
)

var ServiceInstanceIDAttr = attribute.StringValue(uuid.NewV4().String())

const ServiceInstanceIDKey = "service.instance.id"

var ServiceInstanceKeyValue = attribute.KeyValue{Key: ServiceInstanceIDKey, Value: ServiceInstanceIDAttr}
