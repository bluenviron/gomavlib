package gomavlib

import (
	"testing"
)

func TestEndpointBroadcast(t *testing.T) {
	doTest(t, EndpointUDPBroadcast{"127.255.255.255:5602", ":5601"},
		EndpointUDPBroadcast{"127.255.255.255:5601", ":5602"})
}
