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
	actual_ip, _ := GetIP()
	if expected_ip != actual_ip {
		t.Errorf(
			"Invalid ip, expected %s to equal %s",
			expected_ip,
			actual_ip,
		)
	}
}
