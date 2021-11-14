package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/moezakura/auto-unlock/pkg/arpscan"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
	"github.com/moezakura/auto-unlock/pkg/api"
	"github.com/moezakura/auto-unlock/pkg/config"
	"github.com/moezakura/auto-unlock/pkg/soundmeter"
	"github.com/moezakura/auto-unlock/pkg/timedb"
)

var (
	lastUnlockedAt = time.Now()
)

func main() {
	home := os.Getenv("HOME")
	isVerbose := false
	configPath := ""
	flag.BoolVar(&isVerbose, "v", false, "")
	flag.StringVar(&configPath, "c", home+"/.auto-unlock.yaml", "")
	flag.Parse()

	cfg := &config.Config{}
	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(configPath),
		flags.NewBackend(),
	)
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("config: %#v\n", cfg)

	s := soundmeter.NewSoundMeter()
	as := arpscan.NewArpScan()
	td := timedb.NewTimeDB()
	client := api.NewApi(cfg)

	go func() {
		err := s.Exec()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		as.Run(1*time.Second, 5*time.Second)
	}()

	lt := time.Now().UnixMilli()
	for {
		l := s.GetLine()
		if l == "" {
			continue
		}

		n, err := strconv.Atoi(l)
		if err != nil {
			log.Printf("failed to parse value: %+v", err)
			continue
		}

		td.AddIntWithLife(n, 2*time.Second)

		avg := td.GetAVGByAllInt()
		if isVerbose {
			log.Printf("%s, %dms, avg: %f\n", l, time.Now().UnixMilli()-lt, avg)
		}

		if avg > float64(cfg.SoundLevel) {
			now := time.Now()
			if lastUnlockedAt.Add(3 * time.Second).After(now) {
				continue
			}

			macOk := false
			for _, m := range cfg.HostMacAddress {
				if as.Exist(m) {
					macOk = true
					break
				}
			}

			if !macOk {
				log.Printf("not found mac address")
				continue
			}

			lastUnlockedAt = now
			go func() {
				err := client.UnlockWithAll()
				if err != nil {
					log.Printf("failed to unlock by client: %v", err)
				}
			}()
		}

		lt = time.Now().UnixMilli()
	}
}
