package gomavlib

import (
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
		addrs, err := intf.Addrs()
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

// EndpointUDPBroadcast sets up a endpoint that works with UDP broadcast packets.
type EndpointUDPBroadcast struct {
	// the broadcast address to which sending outgoing frames, example: 192.168.5.255:5600
	BroadcastAddress string
	// (optional) the listening address. if empty, it will be computed
	// from the broadcast address.
	LocalAddress string
}

type endpointUDPBroadcast struct {
	conf          EndpointUDPBroadcast
	pc            net.PacketConn
	broadcastAddr net.Addr

	terminate chan struct{}
}

func (conf EndpointUDPBroadcast) init() (Endpoint, error) {
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

	pc, err := net.ListenPacket("udp4", conf.LocalAddress)
	if err != nil {
		return nil, err
	}

	iport, _ := strconv.Atoi(port)

	t := &endpointUDPBroadcast{
		conf:          conf,
		pc:            pc,
		broadcastAddr: &net.UDPAddr{IP: broadcastIP, Port: iport},
		terminate:     make(chan struct{}),
	}
	return t, nil
}

func (t *endpointUDPBroadcast) isEndpoint() {}

func (t *endpointUDPBroadcast) Conf() EndpointConf {
	return t.conf
}

func (t *endpointUDPBroadcast) Label() string {
	return fmt.Sprintf("udp:%s", t.broadcastAddr)
}

func (t *endpointUDPBroadcast) Close() error {
	close(t.terminate)
	t.pc.Close()
	return nil
}

func (t *endpointUDPBroadcast) Read(buf []byte) (int, error) {
	// read WITHOUT deadline. Long periods without packets are normal since
	// we're not directly connected to someone.
	n, _, err := t.pc.ReadFrom(buf)
	// wait termination, do not report errors
	if err != nil {
		<-t.terminate
		return 0, errorTerminated
	}

	return n, nil
}

func (t *endpointUDPBroadcast) Write(buf []byte) (int, error) {
	err := t.pc.SetWriteDeadline(time.Now().Add(netWriteTimeout))
	if err != nil {
		return 0, err
	}
	return t.pc.WriteTo(buf, t.broadcastAddr)
}
