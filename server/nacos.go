package server

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"mock-server/config"
	"net"
)

func RegisterWithNacos(c config.NacosConfig, serviceName string, port uint64) error {
	clientConfig := constant.ClientConfig{
		NamespaceId: c.NamespaceId,
	}
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: c.ServerAddr,
			Port:   c.ServerPort,
		},
	}

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfig,
		},
	)
	if err != nil {
		return err
	}

	ip, err := GetLocalIP()
	if err != nil {
		return err
	}

	// 注册服务
	b, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"env": "mock-server",
		},
	})
	if err != nil {
		return err
	}
	if b {
		log.Printf("Successfully registered service %s with Nacos.\n", serviceName)
	}

	return err
}

// GetLocalIP 返回本机的非环回 IPv4 地址
func GetLocalIP() (string, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addresses {
		// 检查网络地址是否是 IP 地址
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			// 确保 IP 地址是 IPv4 地址
			ip4 := ipNet.IP.To4()
			if ip4 != nil {
				return ip4.String(), nil
			}
		}
	}
	return "", errors.New("cannot find local IP address")
}
