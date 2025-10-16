package gomavlib

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"
)

// ipByBroadcastIP returns the ip of an interface associated with given broadcast ip
func ipByBroadcastIP(target net.IP) net.IP {
	intfs, err := net.Interfaces()
	if err != nil {
		return nil
	}

	for _, intf := range intfs {
		var addrs []net.Addr
		addrs, err = intf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipn, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipn.IP.To4()
			if ip == nil {
				continue
			}

			broadcastIP := net.IP(make([]byte, 4))
			for i := range ip {
				broadcastIP[i] = ip[i] | ^ipn.Mask[i]
			}
			if reflect.DeepEqual(broadcastIP, target) {
				return ip
			}
		}
	}
	return nil
}

type wrappedPacketConn struct {
	pc            net.PacketConn
	writeTimeout  time.Duration
	broadcastAddr net.Addr
}

func (r *wrappedPacketConn) Read(p []byte) (int, error) {
	// read WITHOUT deadline. Long periods without packets are normal since
	// we're not directly connected to someone.
	n, _, err := r.pc.ReadFrom(p)
	return n, err
}

func (r *wrappedPacketConn) Write(p []byte) (int, error) {
	err := r.pc.SetWriteDeadline(time.Now().Add(r.writeTimeout))
	if err != nil {
		return 0, err
	}
	return r.pc.WriteTo(p, r.broadcastAddr)
}

func (r *wrappedPacketConn) Close() error {
	return r.pc.Close()
}

// EndpointUDPBroadcast sets up a endpoint that works with UDP broadcast packets.
type EndpointUDPBroadcast struct {
	// broadcast address to which sending outgoing frames, example: 192.168.5.255:5600
	BroadcastAddress string

	// (optional) listening address. if empty, it will be computed
	// from the broadcast address.
	LocalAddress string
}

func (conf EndpointUDPBroadcast) init(node *Node) (Endpoint, error) {
	ipString, port, err := net.SplitHostPort(conf.BroadcastAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid broadcast address")
	}

	broadcastIP := net.ParseIP(ipString)
	if broadcastIP == nil {
		return nil, fmt.Errorf("invalid IP")
	}

	broadcastIP = broadcastIP.To4()
	if broadcastIP == nil {
		return nil, fmt.Errorf("invalid IP")
	}

	if conf.LocalAddress == "" {
		localIP := ipByBroadcastIP(broadcastIP)
		if localIP == nil {
			return nil, fmt.Errorf("cannot find local address associated with given broadcast address")
		}
		conf.LocalAddress = fmt.Sprintf("%s:%s", localIP, port)
	} else {
		_, _, err = net.SplitHostPort(conf.LocalAddress)
		if err != nil {
			return nil, fmt.Errorf("invalid local address")
		}
	}

	iport, _ := strconv.Atoi(port)

	broadcastAddr := &net.UDPAddr{IP: broadcastIP, Port: iport}

	e := &endpointClient{
		node: node,
		conf: EndpointCustomClient{
			Connect: func(_ context.Context) (net.Conn, error) {
				pc, err2 := net.ListenPacket("udp4", conf.LocalAddress)
				if err2 != nil {
					return nil, err
				}

				return &rwcToConn{&wrappedPacketConn{
					pc:            pc,
					writeTimeout:  node.WriteTimeout,
					broadcastAddr: broadcastAddr,
				}}, nil
			},
			Label: fmt.Sprintf("udp:%s", broadcastAddr),
		},
	}
	err = e.initialize()
	return e, err
}
