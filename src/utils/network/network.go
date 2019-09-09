package network

import (
	"errors"
	"net"
)

type IP interface {
	String() string
	IsLoopback() bool
	To4() IP
}

type Address interface {
	IP() IP
}

type Interface interface {
	Addrs() ([]Address, error)
	FlagUp() net.Flags
	FlagLoopback() net.Flags
}

func extractIP(ifaces []Interface, err error) (string, error) {
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.FlagUp() == 0 {
			continue // interface down
		}
		if iface.FlagLoopback() != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip := addr.IP()
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func GetIP() (string, error) {
	rawIfces, err := net.Interfaces()
	ifaces := []Interface{}
	for _, v := range rawIfces {
		ifaces = append(ifaces, MkRealInterface(v))
	}
	return extractIP(ifaces, err)
}
