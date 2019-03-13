package gomavlib

import (
	"fmt"
	"net"
	"reflect"
)

// ipByBroadcastIp returns the ip of an interface associated with given broadcast ip
func ipByBroadcastIp(target net.IP) net.IP {
	if intfs, err := net.Interfaces(); err == nil {
		for _, intf := range intfs {
			if addrs, err := intf.Addrs(); err == nil {
				for _, addr := range addrs {
					if ipn, ok := addr.(*net.IPNet); ok {
						if ip := ipn.IP.To4(); ip != nil {
							broadcastIp := net.IP(make([]byte, 4))
							for i := range ip {
								broadcastIp[i] = ip[i] | ^ipn.Mask[i]
							}
							if reflect.DeepEqual(broadcastIp, target) == true {
								return ip
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// TransportUdpBroadcast reads and writes frames through UDP broadcast packets.
type TransportUdpBroadcast struct {
	// the broadcast address to which sending outgoing frames, example: 192.168.5.255:5600
	BroadcastAddr string
	// (optional) the address on which listening. if empty, it will be inferred
	// from the broadcast address.
	LocalAddr string
}

type transportUdpBroadcast struct {
	conf          TransportUdpBroadcast
	packetConn    net.PacketConn
	broadcastAddr net.Addr
}

func (conf TransportUdpBroadcast) init(node *Node) (transport, error) {
	_, port, err := net.SplitHostPort(conf.BroadcastAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid broadcast address")
	}
	broadcastAddr, err := net.ResolveUDPAddr("udp4", conf.BroadcastAddr)
	if err != nil {
		return nil, err
	}

	if conf.LocalAddr == "" {
		ip := ipByBroadcastIp(broadcastAddr.IP[:4])
		if ip == nil {
			return nil, fmt.Errorf("cannot find local address associated to given broadcast address")
		}
		conf.LocalAddr = fmt.Sprintf("%s:%s", ip, port)

	} else {
		_, _, err = net.SplitHostPort(conf.LocalAddr)
		if err != nil {
			return nil, fmt.Errorf("invalid local address")
		}
	}

	packetConn, err := net.ListenPacket("udp4", conf.LocalAddr)
	if err != nil {
		return nil, err
	}

	br := &transportUdpBroadcast{
		conf:          conf,
		packetConn:    packetConn,
		broadcastAddr: broadcastAddr,
	}

	tc := TransportCustom{
		ReadWriteCloser: br,
	}
	return tc.init(node)
}

func (t *transportUdpBroadcast) Close() error {
	return t.packetConn.Close()
}

func (t *transportUdpBroadcast) Read(buf []byte) (int, error) {
	n, _, err := t.packetConn.ReadFrom(buf)
	return n, err
}

func (t *transportUdpBroadcast) Write(buf []byte) (int, error) {
	return t.packetConn.WriteTo(buf, t.broadcastAddr)
}
