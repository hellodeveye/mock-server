package server

import (
	"context"
	"log"
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
		// 检查请求的 URL 是否有对应的处理函数
		_, pattern := mux.Handler(r)
		if r.URL.Path != pattern {
			// 如果没有找到匹配的路由，重定向到自定义的 404 页面
			//s := ctx.Value("server").(config.Server)
			w.Header().Set("env", "dev-1")
			redirectURL := "http://10.42.37.48:8080/" + r.URL.Path
			client := http.Client{}
			body, _ := r.GetBody()
			request, _ := http.NewRequest(r.Method, redirectURL, body)
			request.Header = r.Header
			do, err := client.Do(request)
			if err != nil {
				log.Println("请求失败")
				return
			}
			do.Write(w)
		} else {
			// 正常处理请求
			mux.ServeHTTP(w, r)
		}
	})
}
