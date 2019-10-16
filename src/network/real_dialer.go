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
	fmt.Printf("Dialing %v", ip)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, constants.KADEMLIA_PORT))
	if err != nil {
		errCh <- err
		fmt.Printf(err.Error())
		return
	}
	fmt.Printf("Worlds")
	ch <- conn
	return
}

func timeout(errCh chan error) {
	time.Sleep(30 * time.Second)
	errCh <- errors.New("Connection timed out")
}

func (r *RealDialer) DialIP(ip string) (io.ReadWriter, error) {
	ch := make(chan io.ReadWriter)
	errCh := make(chan error)
	go dial(ip, ch, errCh)
	go timeout(errCh)

	select {
  case err := <-errCh:
    fmt.Printf(err.Error())
		return nil, err
  case conn := <-ch:
		return conn, nil
	}
}
