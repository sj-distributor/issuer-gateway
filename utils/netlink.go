package utils

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net"
)

func GetLocalId() string {
	// 获取本机的所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		logx.Errorf("无法获取网络接口: %s", err)
		return ""
	}

	// 遍历每个网络接口并获取它的IP地址
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			logx.Errorf("无法获取接口 %s 的IP地址：%v", iface.Name, err)
			continue
		}

		// 遍历该接口的IP地址列表
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				logx.Errorf("无法解析IP地址：%v", err)
				continue
			}

			// 排除IPv6地址和回环地址
			if ip.To4() != nil && !ip.IsLoopback() && iface.Name == "en1" {
				logx.Infof("接口 %s 的IPv4地址：%s", iface.Name, ip)
				return fmt.Sprintf("%s", ip)
			}
		}
	}
	return ""
}
