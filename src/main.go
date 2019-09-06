package main

import (
	networkutils "utils/network"
	"fmt"
)

func main() {
	ip, err := networkutils.GetIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)
}
