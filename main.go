package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Configpath string = "./conf.json"
	cfg        Config
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	log.SetOutput(gin.DefaultWriter)
	log.SetPrefix("[metrics] ")

	//1.get all process and listening ipaddr/port
	err := LoadConfig(&cfg, Configpath)
	log.Printf("loadconfig is successfully\n")

	if err != nil {
		log.Fatal("loadconfig fialed")
	}

	//2.slice []custome and each custime implment the func of the Desc and Collect
	allp := InitAllProcessMetric(cfg.Process)

	//3.register all custome process metric
	prometheus.MustRegister(allp...)

	//4.start http service
	r := gin.New()
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

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	listenaddr := fmt.Sprintf("%s:%s", cfg.Ipaddr, cfg.Port)

	if err := r.Run(listenaddr); err != nil {
		log.Fatal("http server start failed")
	}
}
