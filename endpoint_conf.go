package gomavlib

// EndpointConf is the interface implemented by all endpoint configurations.
type EndpointConf interface {
	init(*Node) (Endpoint, error)
}
