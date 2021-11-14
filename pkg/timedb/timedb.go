package timedb

import (
	"sync"
	"time"
)

type TimeDB struct {
	values []*value
	lock   sync.Mutex
}

func NewTimeDB() *TimeDB {
	d := &TimeDB{
		values: make([]*value, 0),
	}

	go func() {
		t := time.NewTicker(3 * time.Second)
		defer t.Stop()
		for {
			d.execGC()
			<-t.C
		}
	}()

	return d
}

func (d *TimeDB) execGC() {
	d.lock.Lock()
	defer d.lock.Unlock()

	now := time.Now()
	nv := make([]*value, 0, len(d.values))
	for _, v := range d.values {
		if now.After(v.Deadline) {
			continue
		}
		nv = append(nv, v)
	}
	d.values = nv
}

func (d *TimeDB) AddIntWithLife(v int, life time.Duration) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.values = append(d.values, &value{
		RawValue: v,
		Deadline: time.Now().Add(life),
	})
}

func (d *TimeDB) GetAVGByAllInt() float64 {
	d.execGC()

	d.lock.Lock()
	defer d.lock.Unlock()

	if len(d.values) == 0 {
		return 0
	}

	sum := 0
	for _, v := range d.values {
		sum += v.RawValue.(int)
	}
	return float64(sum / len(d.values))
}
