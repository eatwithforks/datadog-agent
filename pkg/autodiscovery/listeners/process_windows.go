// +build windows

package listeners

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// This should be done another way
func getProcessPorts(pid int32) ([]int32, error) {
	out, err := callNetstat()
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve used ports: %s", err)
	}

	return extractPorts(string(out), pid)
}

func callNetstat() ([]byte, error) {
	bin, err := exec.LookPath("netstat")
	if err != nil {
		return nil, fmt.Errorf("couldn't find netstat installed: %s", err)
	}
	cmd := exec.Command(bin, "-ano")

	return cmd.Output()
}

func extractPorts(raw string, pid int32) ([]int32, error) {
	ports := []int32{}

	for _, line := range strings.Split(raw, "\n") {
		fields := strings.Fields(line)

		if len(fields) > 1 && fields[len(fields)-1] == strconv.Itoa(int(pid)) {
			port, err := getPort(fields[1])
			if err == nil {
				ports = append(ports, int32(port))
			}
		}
	}

	return ports, nil
}

func getPort(addr string) (int, error) {
	fields := strings.Split(addr, ":")
	if len(fields) <= 1 {
		return 0, fmt.Errorf("wrong format for addr: %s", addr)
	}

	return strconv.Atoi(fields[len(fields)-1])
}
