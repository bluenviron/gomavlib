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

// TransportUdpBroadcast sets up a transport that works through UDP broadcast packets.
type TransportUdpBroadcast struct {
	// the broadcast address to which sending outgoing frames, example: 192.168.5.255:5600
	BroadcastAddress string
	// (optional) the listening. if empty, it will be computed
	// from the broadcast address.
	LocalAddress string
}

type transportUdpBroadcast struct {
	conf          TransportUdpBroadcast
	packetConn    net.PacketConn
	broadcastAddr net.Addr
	terminate     chan struct{}
}

func (conf TransportUdpBroadcast) init() (transport, error) {
	_, port, err := net.SplitHostPort(conf.BroadcastAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid broadcast address")
	}
	broadcastAddr, err := net.ResolveUDPAddr("udp4", conf.BroadcastAddress)
	if err != nil {
		return nil, err
	}

	if conf.LocalAddress == "" {
		ip := ipByBroadcastIp(broadcastAddr.IP[:4])
		if ip == nil {
			return nil, fmt.Errorf("cannot find local address associated to given broadcast address")
		}
		conf.LocalAddress = fmt.Sprintf("%s:%s", ip, port)

	} else {
		_, _, err = net.SplitHostPort(conf.LocalAddress)
		if err != nil {
			return nil, fmt.Errorf("invalid local address")
		}
	}

	packetConn, err := net.ListenPacket("udp4", conf.LocalAddress)
	if err != nil {
		return nil, err
	}

	t := &transportUdpBroadcast{
		conf:          conf,
		packetConn:    packetConn,
		broadcastAddr: broadcastAddr,
		terminate:     make(chan struct{}, 1),
	}
	return t, nil
}

func (*transportUdpBroadcast) isTransport() {
}

func (t *transportUdpBroadcast) Close() error {
	t.terminate <- struct{}{}
	t.packetConn.Close()
	return nil
}

func (t *transportUdpBroadcast) Read(buf []byte) (int, error) {
	// read WITHOUT deadline
	n, _, err := t.packetConn.ReadFrom(buf)

	// wait termination, do not report errors
	if err != nil {
		<-t.terminate
		return 0, errorTerminated
	}

	return n, nil
}

func (t *transportUdpBroadcast) Write(buf []byte) (int, error) {
	return t.packetConn.WriteTo(buf, t.broadcastAddr)
}
