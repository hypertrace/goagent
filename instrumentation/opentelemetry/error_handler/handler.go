package error_handler

import (
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// type assertion
var _ otel.ErrorHandler = (*ErrorHandler)(nil)

type ErrorHandler struct {
	logger *zap.Logger
}

func (e *ErrorHandler) Handle(err error) {
	e.logger.Error("opentelemetry error", zap.Error(err))
}

func Init(logger *zap.Logger) {
	handler := &ErrorHandler{logger: logger}
	otel.SetErrorHandler(handler)
}
