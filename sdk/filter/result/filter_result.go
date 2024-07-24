package result

type KeyValueString struct {
	key   string
	value string
}

type Decorations struct {
	RequestHeaderInjections []KeyValueString
}

type FilterResult struct {
	Block              bool
	ResponseStatusCode int32
	ResponseMessage    string
	Decorations        Decorations
}
