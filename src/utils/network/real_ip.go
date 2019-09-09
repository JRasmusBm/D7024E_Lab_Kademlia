package network

import (
	"net"
)

type RealIP struct {
	ip net.IP
}

func (i *RealIP) String() string {
	return i.ip.String()
}

func (i *RealIP) IsLoopback() bool {
	return i.ip.IsLoopback()
}

func (i *RealIP) To4() IP {
	return &RealIP{i.ip.To4()}
}

func MkRealIP(ip net.IP) IP {
	return &RealIP{ip: ip}
}
