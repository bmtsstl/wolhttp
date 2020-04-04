package wol

import (
	"net"
)

var (
	defaultNetwork    = "udp"
	defaultLocalAddr  = (*net.UDPAddr)(nil)
	defaultRemoteAddr = &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 9,
	}
)

func WakeOnLAN(hw net.HardwareAddr) error {
	return WakeOnLANRequest(&Request{
		Network:      defaultNetwork,
		LocalAddr:    defaultLocalAddr,
		RemoteAddr:   defaultRemoteAddr,
		HardwareAddr: hw,
	})
}

type Request struct {
	Network               string
	LocalAddr, RemoteAddr *net.UDPAddr
	HardwareAddr          net.HardwareAddr
}

func WakeOnLANRequest(req *Request) (err error) {
	c, err := net.DialUDP(req.Network, req.LocalAddr, req.RemoteAddr)
	if err != nil {
		return err
	}
	defer func() {
		e := c.Close()
		if e != nil && err == nil {
			err = e
		}
	}()

	m := MagicPacket(req.HardwareAddr)
	_, err = c.Write(m)
	return
}

func MagicPacket(hw net.HardwareAddr) []byte {
	s := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	m := make([]byte, len(s)+16*len(hw))
	copy(m, s)
	for i := len(s); i < len(m); i += len(hw) {
		copy(m[i:], []byte(hw))
	}
	return m
}
