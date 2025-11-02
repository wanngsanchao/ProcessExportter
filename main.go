package main

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	Configpath string = "./conf.json"
	cfg        Config
)

func main() {

	//1.get all config including http server and process
	if err := LoadConfig(&cfg, Configpath); err != nil {
		log.Fatal("loadconfig fialed")
	}

	log.Printf("loadconfig is successfully\n")

	//2. register all desc to global desc map
	allp := InitAllProcessMetric(cfg.Process)

	//3.register all custome process metric
	prometheus.MustRegister(allp...)

	//4. init the http router
	r, err := InitRouter()
	if err != nil {
		log.Fatal("init router failed\n")
	}

	log.Printf("init router successfull\n")

	//5. register the handler
	Register_handler(All_interface_handlers, r)
	log.Printf("register the handler successfully\n")

	//6.start http service
	listenaddr := fmt.Sprintf("%s:%s", cfg.Ipaddr, cfg.Port)
	if err := r.Run(listenaddr); err != nil {
		log.Fatal("http server start failed")
	}
}
