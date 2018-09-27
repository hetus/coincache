package coincache

import (
	"fmt"
	"time"

	"github.com/Toorop/go-bittrex"
	"github.com/asdine/storm"
)

type CoinCache struct {
	cfg  *Config
	cli  *bittrex.Bittrex
	db   *storm.DB
	stop bool
	sub  func(model *Model)
}

func (cc *CoinCache) All(to interface{}) error {
	return cc.db.All(to)
}

func (cc *CoinCache) AllByMarket(market string, to interface{}) error {
	return cc.db.Find("Name", market, to)
}

func (cc *CoinCache) Close() error {
	return cc.db.Close()
}

func (cc *CoinCache) Start() {
	cc.stop = false
	if cc.cfg.Debug {
		fmt.Println("debug: starting")
	}

	for {
		if cc.cfg.Debug {
			fmt.Println("debug: interval")
		}
		if cc.stop {
			if cc.cfg.Debug {
				fmt.Println("debug: stopped")
			}
			break
		}

		markets, err := cc.cli.GetMarketSummaries()
		if err != nil {
			fmt.Println("error: get markets:", err)
		} else {
			for _, market := range markets {
				m := Model{
					Ask:    float(market.Ask),
					Bid:    float(market.Bid),
					High:   float(market.High),
					Last:   float(market.Last),
					Low:    float(market.Low),
					Name:   market.MarketName,
					Volume: float(market.Volume),
				}

				cc.sub(&m)

				err = cc.db.Save(&m)
				if err != nil {
					fmt.Println("error: save model:", err)
				}
			}
		}

		select {
		case <-time.After(cc.cfg.Interval):
			continue
		}
	}
}

func (cc *CoinCache) Stop() {
	if cc.cfg.Debug {
		fmt.Println("debug: stopping")
	}
	cc.stop = true
}

func (cc *CoinCache) Subscribe(fn func(model *Model)) {
	if cc.cfg.Debug {
		fmt.Println("debug: subscribed")
	}
	cc.sub = fn
}

func New(cfg *Config) (*CoinCache, error) {
	var err error
	cc := &CoinCache{
		cfg:  cfg,
		cli:  bittrex.New("KEY", "SECRET"),
		stop: false,
		sub:  func(model *Model) {},
	}

	cc.db, err = storm.Open(cfg.Database)
	if err != nil {
		return nil, err
	}
	if cfg.Debug {
		fmt.Println("debug: database connected")
	}

	err = cc.db.Init(&Model{})
	return cc, err
}
