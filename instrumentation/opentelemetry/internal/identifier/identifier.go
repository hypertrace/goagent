package identifier

import (
	uuid "github.com/satori/go.uuid"
	"go.opentelemetry.io/otel/attribute"
)

var instanceId, _ = uuid.NewV4()
var ServiceInstanceIDAttr = attribute.StringValue(instanceId.String())

const ServiceInstanceIDKey = "service.instance.id"

var ServiceInstanceKeyValue = attribute.KeyValue{Key: ServiceInstanceIDKey, Value: ServiceInstanceIDAttr}
