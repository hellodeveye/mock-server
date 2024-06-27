package server

import (
	"context"
	"io"
	"log"
	"mock-server/config"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("请求方法: %s, 请求URL: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func NotFoundMiddleware(mux *http.ServeMux, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router := ctx.Value("router").(config.RouterConfig)
		if !router.Enabled {
			mux.ServeHTTP(w, r)
			return
		}
		// 检查请求的 URL 是否有对应的处理函数
		_, pattern := mux.Handler(r)
		if r.URL.Path != pattern {
			// 如果没有找到匹配的路由，重定向到自定义的 404 页面
			s := ctx.Value("server").(config.Server)
			redirectURL := router.GatewayAddr + s.Name + r.URL.Path
			client := http.Client{}
			request, _ := http.NewRequest(r.Method, redirectURL, r.Body)
			r.Header.Del("env")
			request.Header = r.Header
			response, err := client.Do(request)
			if err != nil {
				http.Error(w, "Failed to forward request", http.StatusInternalServerError)
				return
			}
			defer func() {
				_ = response.Body.Close()
			}()
			// 复制响应头
			for key, values := range response.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			// 设置状态码
			w.WriteHeader(response.StatusCode)
			// 复制响应体
			_, err = io.Copy(w, response.Body)
			if err != nil {
				log.Println("Failed to copy response body:", err)
			}
		} else {
			// 正常处理请求
			mux.ServeHTTP(w, r)
		}
	})
}
