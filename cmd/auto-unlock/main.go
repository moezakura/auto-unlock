package main

import (
	"log"
	"strconv"
	"time"

	"github.com/moezakura/auto-unlock/pkg/timedb"

	"github.com/moezakura/auto-unlock/pkg/soundmeter"
)

var (
	soundLevels = make([]int, 10)
)

func main() {
	s := soundmeter.NewSoundMeter()
	td := timedb.NewTimeDB()

	go func() {
		err := s.Exec()
		if err != nil {
			panic(err)
		}
	}()

	//lt := time.Now().UnixMilli()
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
		// fmt.Printf("%s, %dms, avg: %f\n", l, time.Now().UnixMilli()-lt, avg)

		if avg > 8500 {
			log.Println("UNLOCK!!")
		}

		//lt = time.Now().UnixMilli()
	}
}
