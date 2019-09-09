package network

import (
	"net"
)

type RealInterface struct {
	net net.Interface
}

func (ifs *RealInterface) FlagLoopback() net.Flags {
	return ifs.net.Flags & net.FlagLoopback
}

func (ifs *RealInterface) FlagUp() net.Flags {
	return ifs.net.Flags & net.FlagUp
}

func (ifs *RealInterface) Addrs() ([]Address, error) {
	rawAddrs, err := ifs.net.Addrs()
	addrs := []Address{}
	for _, v := range rawAddrs {
		addrs = append(addrs, &RealAddress{addr: v})
	}
	return addrs, err
}

func MkRealInterface(iface net.Interface) Interface {
	return &RealInterface{net: iface}
}
