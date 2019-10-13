package network

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
	"utils/constants"
)

type RealDialer struct{}

func dial(ip string, ch chan io.ReadWriter, errCh chan error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, constants.KADEMLIA_PORT))
	ch <- conn
	errCh <- err
	return
}

func timeout(errCh chan error) {
	time.Sleep(2 * time.Second)
	errCh <- errors.New("Connection timed out")
}

func (r *RealDialer) DialIP(ip string) (io.ReadWriter, error) {
	ch := make(chan io.ReadWriter)
	errCh := make(chan error)
	go dial(ip, ch, errCh)
	go timeout(errCh)

	var err error
	var conn io.ReadWriter
	select {
	case err = <-errCh:
	case conn = <-ch:
		err = <-errCh
	}
	return conn, err
}
