package network

import (
	"os/exec"
	"strings"
	"testing"
)

func TestIntegrationLinux(t *testing.T) {
	command := exec.Command("hostname", "-I")
	output, _ := command.Output()
	expected_ip := strings.Split(string(output), " ")[0]
	var networkUtils NetworkUtils = &RealNetworkUtils{}
	actual_ip, err := networkUtils.GetIP()
	if err != nil {
		t.Errorf(err.Error())
	}
	if expected_ip != actual_ip {
		t.Errorf(
			"Invalid ip, expected %s to equal %s",
			expected_ip,
			actual_ip,
		)
	}
}
