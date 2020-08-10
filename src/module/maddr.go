package module

import (
	"fmt"
	"learn/mywebcrawler/errors"
	"net"
	"strconv"
)

type mAddr struct {
	network string
	address string
}

func (maddr *mAddr) Network() string {
	return maddr.network
}

func (maddr *mAddr) String() string {
	return maddr.address
}

func NewAddr(network string, ip string, port uint64) (net.Addr, error) {
	if network != "http" && network != "https" {
		errMsg := fmt.Sprintf("illegal network for module address: %v", network)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	if parsedIP := net.ParseIP(ip); parsedIP == nil {
		errMsg := fmt.Sprintf("illegal IP for module address: %v", ip)
		return nil, errors.NewIllegalParameterError(errMsg)
	}

	return &mAddr{
		network: network,
		address: ip + ":" + strconv.Itoa(int(port)),
	}, nil
}
