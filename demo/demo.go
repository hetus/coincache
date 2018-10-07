package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hetus/coincache"
)

func main() {
	// Setup coincache configuration.
	cfg := &coincache.Config{
		Database: "cc-10s.db",     // Name of database file.
		Debug:    true,            // Run in debug mode.
		Interval: 6 * time.Second, // Interval to download market data.
	}
	// Setup a new coincache server using configuration above.
	cc, err := coincache.New(cfg)
	if err != nil {
		log.Fatal("coincache new:", err)
	}
	defer cc.Close() // Close internal connections.

	// The subscribe method acts as a middleware of the market
	// data that can be used for indicator calculations or etc.
	// In this demo we calculate a simple 10 point moving
	// average on the Last price as an example of what can be done.
	markets := make(map[string][]float64)
	cc.Subscribe(func(model *coincache.Model) {
		// Setup initial slice for market if not found.
		ticks, found := markets[model.Name]
		if !found {
			ticks = make([]float64, 0, 10)
		}

		// Limit our range of values to 10 entries.
		if len(ticks) == 10 {
			ticks = append(ticks[1:], model.Last)
		} else {
			ticks = append(ticks, model.Last)
		}

		// Stuff out new ticks slice back into the market
		// for later usage.  For simplicity we will not
		// be doing much in this demo except counting the
		// markets below.
		markets[model.Name] = ticks
	})

	go cc.Start() // Start the coincache server.

	// For demonstration purposes we will kill this process after
	// running for 1 minute.
	time.Sleep(1 * time.Minute)
	cc.Stop()
	if cfg.Debug {
		fmt.Println("debug: models saved:", len(markets))
	}

	// The following is an example of fetching all of the market
	// data stored in the database.  This can take some time to load
	// as it gets bigger.
	models := make([]coincache.Model, 0)
	err = cc.All(&models)
	if err != nil {
		log.Fatal("error: all models:", err)
	}
	fmt.Println("debug: all models length:", len(models))

	// You can also return all markets that match a given market.
	// The same note about load time applies as above.
	btc := make([]coincache.Model, 0)
	err = cc.AllByMarket("USDT-BTC", &btc)
	if err != nil {
		log.Fatal("error: all models by market:", err)
	}
	fmt.Println("debug: btc models:", len(btc))
}
