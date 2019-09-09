package network

import (
	"net"
)

type RealAddress struct {
	addr net.Addr
}

func (a *RealAddress) IP() IP {
	var ip net.IP

	switch v := a.addr.(type) {
	case *net.IPNet:
		ip = v.IP
	}
	return MkRealIP(ip)
}
