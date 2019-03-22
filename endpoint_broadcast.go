package gomavlib

import (
	"fmt"
	"net"
	"reflect"
	"time"
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

// EndpointUdpBroadcast sets up a endpoint that works through UDP broadcast packets.
type EndpointUdpBroadcast struct {
	// the broadcast address to which sending outgoing frames, example: 192.168.5.255:5600
	BroadcastAddress string
	// (optional) the listening address. if empty, it will be computed
	// from the broadcast address.
	LocalAddress string
}

type endpointUdpBroadcast struct {
	conf          EndpointUdpBroadcast
	packetConn    net.PacketConn
	broadcastAddr net.Addr
	terminate     chan struct{}
}

func (conf EndpointUdpBroadcast) init() (endpoint, error) {
	_, port, err := net.SplitHostPort(conf.BroadcastAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid broadcast address")
	}
	broadcastAddr, err := net.ResolveUDPAddr("udp4", conf.BroadcastAddress)
	if err != nil {
		return nil, err
	}

	if conf.LocalAddress == "" {
		ip := ipByBroadcastIp(broadcastAddr.IP.To4())
		if ip == nil {
			return nil, fmt.Errorf("cannot find local address associated with given broadcast address")
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

	t := &endpointUdpBroadcast{
		conf:          conf,
		packetConn:    packetConn,
		broadcastAddr: broadcastAddr,
		terminate:     make(chan struct{}, 1),
	}
	return t, nil
}

func (t *endpointUdpBroadcast) Desc() string {
	return fmt.Sprintf("udp %s", t.broadcastAddr)
}

func (t *endpointUdpBroadcast) Close() error {
	t.terminate <- struct{}{}
	t.packetConn.Close()
	return nil
}

func (t *endpointUdpBroadcast) Read(buf []byte) (int, error) {
	// read WITHOUT deadline. Long periods without packets are normal since
	// we're not directly connected to someone.
	n, _, err := t.packetConn.ReadFrom(buf)

	// wait termination, do not report errors
	if err != nil {
		<-t.terminate
		return 0, errorTerminated
	}

	return n, nil
}

func (t *endpointUdpBroadcast) Write(buf []byte) (int, error) {
	err := t.packetConn.SetWriteDeadline(time.Now().Add(netWriteTimeout))
	if err != nil {
		return 0, err
	}
	return t.packetConn.WriteTo(buf, t.broadcastAddr)
}
