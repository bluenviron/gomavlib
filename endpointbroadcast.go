package gomavlib

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
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

// EndpointUdpBroadcast sets up a endpoint that works with UDP broadcast packets.
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

	terminate chan struct{}
}

func (conf EndpointUdpBroadcast) init() (Endpoint, error) {
	ipString, port, err := net.SplitHostPort(conf.BroadcastAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid broadcast address")
	}
	broadcastIp := net.ParseIP(ipString)
	if broadcastIp == nil {
		return nil, fmt.Errorf("invalid IP")
	}
	broadcastIp = broadcastIp.To4()
	if broadcastIp == nil {
		return nil, fmt.Errorf("invalid IP")
	}

	if conf.LocalAddress == "" {
		localIp := ipByBroadcastIp(broadcastIp)
		if localIp == nil {
			return nil, fmt.Errorf("cannot find local address associated with given broadcast address")
		}
		conf.LocalAddress = fmt.Sprintf("%s:%s", localIp, port)

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

	iport, _ := strconv.Atoi(port)

	t := &endpointUdpBroadcast{
		conf:          conf,
		packetConn:    packetConn,
		broadcastAddr: &net.UDPAddr{IP: broadcastIp, Port: iport},
		terminate:     make(chan struct{}),
	}
	return t, nil
}

func (t *endpointUdpBroadcast) isEndpoint() {}

func (t *endpointUdpBroadcast) Conf() interface{} {
	return t.conf
}

func (t *endpointUdpBroadcast) Label() string {
	return fmt.Sprintf("udp:%s", t.broadcastAddr)
}

func (t *endpointUdpBroadcast) Close() error {
	close(t.terminate)
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
