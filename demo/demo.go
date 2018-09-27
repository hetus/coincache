package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hetus/coincache"
)

func main() {
	cfg := &coincache.Config{
		Database: "cc-10s.db",
		Debug:    true,
		Interval: 10 * time.Second,
	}
	cc, err := coincache.New(cfg)
	if err != nil {
		log.Fatal("coincache new:", err)
	}
	defer cc.Close()

	count := 0
	cc.Subscribe(func(model *coincache.Model) {
		count++
	})

	go cc.Start()

	time.Sleep(30 * time.Second)
	cc.Stop()
	if cfg.Debug {
		fmt.Println("debug: models saved:", count)
	}

	models := make([]coincache.Model, 0)
	err = cc.All(&models)
	if err != nil {
		log.Fatal("error: all models:", err)
	}
	fmt.Println("debug: all models:", len(models))

	btc := make([]coincache.Model, 0)
	err = cc.AllByMarket("USDT-BTC", &btc)
	if err != nil {
		log.Fatal("error: all models by market:", err)
	}
	fmt.Println("debug: btc models:", len(btc))
}
