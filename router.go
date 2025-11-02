package main

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Interface_handler struct {
	Name    string
	Method  string
	Handler gin.HandlerFunc
}

var (
	All_interface_handlers []Interface_handler = []Interface_handler{
		{
			Name:    "/metrics",
			Method:  "GET",
			Handler: gin.WrapH(promhttp.Handler()),
		},
		{
			Name:    "/health",
			Method:  "GET",
			Handler: CheckHealth,
		},
	}
)

func InitRouter() (*gin.Engine, error) {
	r := gin.New()

	if r != nil {
		return nil, errors.New("initrouter failed")
	}

	r.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		// 保留 Gin 的默认颜色格式，仅在最前面添加 "[MY-APP]" 前缀
		return fmt.Sprintf("[MY-APP] %s - [%s] \"%s %s\" %d %s\n",
			params.ClientIP,
			params.TimeStamp.Format("2006-01-02 15:04:05"),
			params.Method,
			params.Path,
			params.StatusCode,
			params.Latency,
		)
	}))

	// 强制启用颜色（确保终端显示彩色）
	gin.ForceConsoleColor()
	return r, nil
}

func Register_handler(all_interface []Interface_handler, r *gin.Engine) {
	for _, single_handler := range all_interface {
		if single_handler.Method == "GET" {
			r.GET(single_handler.Method, single_handler.Handler)
		}

		if single_handler.Method == "POST" {
			r.POST(single_handler.Method, single_handler.Handler)
		}
	}
}
