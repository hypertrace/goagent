package sdk

import "context"

// ReadbackSpan is a span whose attributes can be read. It is important
// to bear in mind that only those attributes being recorded through this
// interface will be available for read. This is, an existing span where
// attributes have been already recorded can be turned into a ReadbackSpan
// but the attributes won't be available.
type ReadbackSpan interface {
	Span
	GetAttribute(key string) interface{}
}

type StartReadbackSpan func(
	ctx context.Context,
	name string,
	opts *SpanOptions,
) (context.Context, ReadbackSpan, func())
