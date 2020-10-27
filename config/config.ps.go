package config

type AgentConfig struct {
	// serviceName identifies the service/process running
	ServiceName *string      `json:"serviceName,omitempty"`
	Reporting   *Reporting   `json:"reporting,omitempty"`
	DataCapture *DataCapture `json:"dataCapture,omitempty"`
}

// GetServiceName returns the serviceName
func (x *AgentConfig) GetServiceName() string {
	if x.ServiceName == nil {
		return ""
	}
	return *x.ServiceName
}

// GetReporting returns the Reporting
func (x *AgentConfig) GetReporting() *Reporting {
	return x.Reporting
}

// GetDataCapture returns the DataCapture
func (x *AgentConfig) GetDataCapture() *DataCapture {
	return x.DataCapture
}

func (x *AgentConfig) loadFromEnv(prefix string, defaultValues *AgentConfig) {
	if x.ServiceName == nil {
		if val, ok := getStringEnv(prefix + "SERVICE_NAME"); ok {
			x.ServiceName = new(string)
			*x.ServiceName = val
		} else if defaultValues != nil && defaultValues.ServiceName != nil {
			x.ServiceName = new(string)
			*x.ServiceName = *defaultValues.ServiceName
		}
	}

	if x.Reporting == nil {
		x.Reporting = new(Reporting)
	}
	x.Reporting.loadFromEnv(prefix+"REPORTING_", defaultValues.Reporting)
	if x.DataCapture == nil {
		x.DataCapture = new(DataCapture)
	}
	x.DataCapture.loadFromEnv(prefix+"DATA_CAPTURE_", defaultValues.DataCapture)
}

type Reporting struct {
	// address represents the host for reporting the traes e.g. api.traceable.ai
	Address *string `json:"address,omitempty"`
	// isSecure permits connecting to the trace endpoint without a certificate
	IsSecure *bool `json:"isSecure,omitempty"`
}

// GetAddress returns the address
func (x *Reporting) GetAddress() string {
	if x.Address == nil {
		return ""
	}
	return *x.Address
}

// GetIsSecure returns the isSecure
func (x *Reporting) GetIsSecure() bool {
	if x.IsSecure == nil {
		return false
	}
	return *x.IsSecure
}

func (x *Reporting) loadFromEnv(prefix string, defaultValues *Reporting) {
	if x.Address == nil {
		if val, ok := getStringEnv(prefix + "ADDRESS"); ok {
			x.Address = new(string)
			*x.Address = val
		} else if defaultValues != nil && defaultValues.Address != nil {
			x.Address = new(string)
			*x.Address = *defaultValues.Address
		}
	}

	if x.IsSecure == nil {
		if val, ok := getBoolEnv(prefix + "IS_SECURE"); ok {
			x.IsSecure = new(bool)
			*x.IsSecure = val
		} else if defaultValues != nil && defaultValues.IsSecure != nil {
			x.IsSecure = new(bool)
			*x.IsSecure = *defaultValues.IsSecure
		}
	}

}

type Message struct {
	Request  *bool `json:"request,omitempty"`
	Response *bool `json:"response,omitempty"`
}

// GetRequest returns the request
func (x *Message) GetRequest() bool {
	if x.Request == nil {
		return false
	}
	return *x.Request
}

// GetResponse returns the response
func (x *Message) GetResponse() bool {
	if x.Response == nil {
		return false
	}
	return *x.Response
}

func (x *Message) loadFromEnv(prefix string, defaultValues *Message) {
	if x.Request == nil {
		if val, ok := getBoolEnv(prefix + "REQUEST"); ok {
			x.Request = new(bool)
			*x.Request = val
		} else if defaultValues != nil && defaultValues.Request != nil {
			x.Request = new(bool)
			*x.Request = *defaultValues.Request
		}
	}

	if x.Response == nil {
		if val, ok := getBoolEnv(prefix + "RESPONSE"); ok {
			x.Response = new(bool)
			*x.Response = val
		} else if defaultValues != nil && defaultValues.Response != nil {
			x.Response = new(bool)
			*x.Response = *defaultValues.Response
		}
	}

}

type DataCapture struct {
	// httpHeaders enables/disables the capture of the request/response headers in HTTP
	HTTPHeaders *Message `json:"httpHeaders,omitempty"`
	// httpBody enables/disables the capture of the request/response body in HTTP
	HTTPBody *Message `json:"httpBody,omitempty"`
	// rpcMetadata enables/disables the capture of the request/response metadata in RPC
	RPCMetadata *Message `json:"rpcMetadata,omitempty"`
	// rpcBody enables/disables the capture of the request/response body in RPC
	RPCBody *Message `json:"rpcBody,omitempty"`
}

// GetHTTPHeaders returns the HTTPHeaders
func (x *DataCapture) GetHTTPHeaders() *Message {
	return x.HTTPHeaders
}

// GetHTTPBody returns the HTTPBody
func (x *DataCapture) GetHTTPBody() *Message {
	return x.HTTPBody
}

// GetRPCMetadata returns the RPCMetadata
func (x *DataCapture) GetRPCMetadata() *Message {
	return x.RPCMetadata
}

// GetRPCBody returns the RPCBody
func (x *DataCapture) GetRPCBody() *Message {
	return x.RPCBody
}

func (x *DataCapture) loadFromEnv(prefix string, defaultValues *DataCapture) {
	if x.HTTPHeaders == nil {
		x.HTTPHeaders = new(Message)
	}
	x.HTTPHeaders.loadFromEnv(prefix+"HTTP_HEADERS_", defaultValues.HTTPHeaders)
	if x.HTTPBody == nil {
		x.HTTPBody = new(Message)
	}
	x.HTTPBody.loadFromEnv(prefix+"HTTP_BODY_", defaultValues.HTTPBody)
	if x.RPCMetadata == nil {
		x.RPCMetadata = new(Message)
	}
	x.RPCMetadata.loadFromEnv(prefix+"RPC_METADATA_", defaultValues.RPCMetadata)
	if x.RPCBody == nil {
		x.RPCBody = new(Message)
	}
	x.RPCBody.loadFromEnv(prefix+"RPC_BODY_", defaultValues.RPCBody)
}
