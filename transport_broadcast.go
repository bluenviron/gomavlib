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

type transportUdpBroadcastChannel struct {
	net.PacketConn
	broadcastAddr net.Addr
}

func (t *transportUdpBroadcastChannel) Write(buf []byte) (int, error) {
	return t.WriteTo(buf, t.broadcastAddr)
}

func (t *transportUdpBroadcastChannel) SetWriteDeadline(ti time.Time) error {
	// causes a stack overflow, disabled
	return nil
}

func (*transportUdpBroadcastChannel) LocalAddr() net.Addr {
	// provided for netTimedConn, not used
	return nil
}

func (*transportUdpBroadcastChannel) RemoteAddr() net.Addr {
	// provided for netTimedConn, not used
	return nil
}

func (*transportUdpBroadcastChannel) Read(buf []byte) (int, error) {
	// provided for netTimedConn, not used
	return 0, nil
}

func (*transportUdpBroadcastChannel) SetDeadline(time.Time) error {
	// provided for netTimedConn, not used
	return nil
}

func (*transportUdpBroadcastChannel) SetReadDeadline(t time.Time) error {
	// provided for netTimedConn, not used
	return nil
}

// TransportUdpBroadcast sends and reads frames through UDP broadcast packets.
type TransportUdpBroadcast struct {
	// the broadcast address to which sending outgoing frames, example: 192.168.5.255:5600
	BroadcastAddr string
	// (optional) the address on which listening. if empty, it will be inferred
	// from the broadcast address.
	LocalAddr string
}

type transportUdpBroadcast struct {
	conf          TransportUdpBroadcast
	node          *Node
	terminate     chan struct{}
	broadcastAddr net.Addr
	packetConn    net.PacketConn
	tconn         *TransportChannel
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

	t := &transportUdpBroadcast{
		conf:          conf,
		node:          node,
		terminate:     make(chan struct{}),
		broadcastAddr: broadcastAddr,
		packetConn:    packetConn,
	}

	t.tconn = &TransportChannel{
		transport: t,
		writer: &netTimedConn{&transportUdpBroadcastChannel{
			t.packetConn,
			t.broadcastAddr,
		}},
	}
	t.node.addTransportChannel(t.tconn)

	return t, nil
}

func (t *transportUdpBroadcast) closePrematurely() {
	t.packetConn.Close()
}

func (t *transportUdpBroadcast) do() {
	listenDone := make(chan struct{})
	go func() {
		defer func() { listenDone <- struct{}{} }()
		defer t.node.removeTransportChannel(t.tconn)

		var buf [netBufferSize]byte
		for {
			n, _, err := t.packetConn.ReadFrom(buf[:])
			if err != nil {
				break
			}
			t.node.processBuffer(nil, buf[:n])
		}
	}()

	select {
	// unexpected error, wait for terminate
	case <-listenDone:
		t.packetConn.Close()
		<-t.node.terminate

	// terminated
	case <-t.node.terminate:
		t.packetConn.Close()
		<-listenDone
	}
}
