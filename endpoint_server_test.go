package gomavlib

import (
	"testing"
)

func TestEndpointServerTCP(t *testing.T) {
	doTest(t, EndpointTCPServer{"127.0.0.1:5601"}, EndpointTCPClient{"127.0.0.1:5601"})
}

func TestEndpointServerUDP(t *testing.T) {
	doTest(t, EndpointUDPServer{"127.0.0.1:5601"}, EndpointUDPClient{"127.0.0.1:5601"})
}
