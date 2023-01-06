package identifier

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

var ServiceInstanceIDAttr = attribute.StringValue(uuid.New().String())

const ServiceInstanceIDKey = "service.instance.id"

var ServiceInstanceKeyValue = attribute.KeyValue{Key: ServiceInstanceIDKey, Value: ServiceInstanceIDAttr}
