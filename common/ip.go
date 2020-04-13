package common

import (
	"errors"
	"net"
)

func GetIntranceIp() (string, error) {
	addressd, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addressd {
		//检查Ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("获取本地IP地址异常")
}
