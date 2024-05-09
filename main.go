package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	// Nacos 配置
	serverConfigs := []constant.ServerConfig{
		{IpAddr: "192.168.2.128", Port: 31146},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
	}

	// 创建 Nacos 客户端
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 注册全局 HTTP 处理函数
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from port %s!", r.Host)
	})

	// 注册多个服务
	services := []struct {
		Port        uint64
		ServiceName string
	}{
		{8080, "billing-center"},
		{8081, "write-off-center"},
	}

	localIP, err := GetLocalIP()
	if err != nil {
		log.Fatalf("Failed to get local IP address: %v", err)
	}
	for _, service := range services {
		go func(srv struct {
			Port        uint64
			ServiceName string
		}) {
			// 服务地址
			addr := fmt.Sprintf(":%d", srv.Port)
			log.Printf("Starting %s on %s", srv.ServiceName, addr)

			// 启动 HTTP 服务
			go func() {
				if err := http.ListenAndServe(addr, nil); err != nil {
					log.Fatalf("Failed to start HTTP server on %s: %v", addr, err)
				}
			}()

			// 注册服务到 Nacos
			_, err = client.RegisterInstance(vo.RegisterInstanceParam{
				Ip:          localIP,
				Port:        srv.Port,
				ServiceName: srv.ServiceName,
				Weight:      1,
				Enable:      true,
				Healthy:     true,
				Ephemeral:   true,
				Metadata: map[string]string{
					"env": "mock-server",
				},
			})
			if err != nil {
				log.Fatalf("Failed to register service %s: %v", srv.ServiceName, err)
			}
		}(service)
	}

	// 阻塞主线程, 避免主程序退出
	select {}
}

// GetLocalIP 返回本机的非环回 IPv4 地址
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
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
