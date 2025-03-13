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
	// broadcast address to which sending outgoing frames, example: 192.168.5.255:5600
	BroadcastAddress string

	// (optional) listening address. if empty, it will be computed
	// from the broadcast address.
	LocalAddress string
}

func (conf EndpointUDPBroadcast) init(node *Node) (Endpoint, error) {
	e := &endpointUDPBroadcast{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}

type endpointUDPBroadcast struct {
	node *Node
	conf EndpointUDPBroadcast

	pc            net.PacketConn
	broadcastAddr net.Addr

	terminate chan struct{}
}

func (e *endpointUDPBroadcast) initialize() error {
	ipString, port, err := net.SplitHostPort(e.conf.BroadcastAddress)
	if err != nil {
		return fmt.Errorf("invalid broadcast address")
	}
	broadcastIP := net.ParseIP(ipString)
	if broadcastIP == nil {
		return fmt.Errorf("invalid IP")
	}
	broadcastIP = broadcastIP.To4()
	if broadcastIP == nil {
		return fmt.Errorf("invalid IP")
	}

	if e.conf.LocalAddress == "" {
		localIP := ipByBroadcastIP(broadcastIP)
		if localIP == nil {
			return fmt.Errorf("cannot find local address associated with given broadcast address")
		}
		e.conf.LocalAddress = fmt.Sprintf("%s:%s", localIP, port)
	} else {
		_, _, err = net.SplitHostPort(e.conf.LocalAddress)
		if err != nil {
			return fmt.Errorf("invalid local address")
		}
	}

	e.pc, err = net.ListenPacket("udp4", e.conf.LocalAddress)
	if err != nil {
		return err
	}

	iport, _ := strconv.Atoi(port)

	e.broadcastAddr = &net.UDPAddr{IP: broadcastIP, Port: iport}
	e.terminate = make(chan struct{})

	return nil
}

func (e *endpointUDPBroadcast) isEndpoint() {}

func (e *endpointUDPBroadcast) Conf() EndpointConf {
	return e.conf
}

func (e *endpointUDPBroadcast) label() string {
	return fmt.Sprintf("udp:%s", e.broadcastAddr)
}

func (e *endpointUDPBroadcast) Close() error {
	close(e.terminate)
	e.pc.Close()
	return nil
}

func (e *endpointUDPBroadcast) Read(buf []byte) (int, error) {
	// read WITHOUT deadline. Long periods without packets are normal since
	// we're not directly connected to someone.
	n, _, err := e.pc.ReadFrom(buf)
	// wait termination, do not report errors
	if err != nil {
		<-e.terminate
		return 0, errTerminated
	}

	return n, nil
}

func (e *endpointUDPBroadcast) Write(buf []byte) (int, error) {
	err := e.pc.SetWriteDeadline(time.Now().Add(e.node.WriteTimeout))
	if err != nil {
		return 0, err
	}
	return e.pc.WriteTo(buf, e.broadcastAddr)
}
