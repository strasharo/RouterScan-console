package routerscan

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Target struct {
	Ip   uint32
	Port uint16
}

func ParseTarget(str string) (*Target, error) {
	host := strings.Split(str, ":")
	if len(host) != 2 {
		return nil, fmt.Errorf("host %s should be in format ip:port", str)
	}
	port, err := strconv.ParseUint(host[1], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("cannot parse port %s: %s", host[1], err.Error())
	}
	ipInt, err := inetAton(host[0])
	if err != nil {
		return nil, err
	}
	return &Target{
		Ip:   ipInt,
		Port: uint16(port),
	}, nil
}

func inetAton(ip string) (uint32, error) {
	ipByte := net.ParseIP(ip).To4()
	if ipByte == nil {
		return 0, fmt.Errorf("cannot parse ip %s", ip)
	}
	var ipInt uint32
	for i := 0; i < len(ipByte); i++ {
		ipInt |= uint32(ipByte[i])
		if i < 3 {
			ipInt <<= 8
		}
	}
	return ipInt, nil
}
