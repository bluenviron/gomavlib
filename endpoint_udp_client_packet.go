package gomavlib

import (
	"context"
	"net"
)

// EndpointUDPClientPacket is a UDP client endpoint that treats each incoming
// UDP datagram as an atomic unit.
//
// Unlike EndpointUDPClient (which uses a stream-oriented bufio reader), this
// endpoint sets PacketOriented = true on the resulting Channel.  After any
// parse error the frame reader discards the remaining bytes of the malformed
// datagram before attempting the next read, so a single bad datagram never
// contaminates parsing of the next one.
//
// Use this endpoint when the remote peer is a MAVLink device that sends
// well-formed datagrams but the link may occasionally deliver corrupted or
// out-of-order packets (e.g. a radio-link UDP tunnel in a lossy RF channel).
type EndpointUDPClientPacket struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (conf EndpointUDPClientPacket) init(node *Node) (Endpoint, error) {
	e := &endpointClientPacket{
		conf: conf,
		endpointClient: &endpointClient{
			node: node,
			conf: EndpointCustomClient{
				Connect: func(ctx context.Context) (net.Conn, error) {
					return (&net.Dialer{}).DialContext(ctx, "udp4", conf.Address)
				},
				Label: "udp:" + conf.Address,
			},
		},
	}
	err := e.endpointClient.initialize()
	return e, err
}

// endpointClientPacket wraps endpointClient and marks itself as
// packet-oriented so that channelProvider enables per-datagram recovery.
type endpointClientPacket struct {
	conf EndpointUDPClientPacket
	*endpointClient
}

// Conf returns the EndpointUDPClientPacket configuration.
func (e *endpointClientPacket) Conf() EndpointConf {
	return e.conf
}

// isPacketOrientedEndpoint implements packetOrientedEndpoint.
func (e *endpointClientPacket) isPacketOrientedEndpoint() bool {
	return true
}
