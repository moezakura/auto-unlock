package arpscan

import (
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/xerrors"
)

type ArpScan struct {
	macAddress map[string]*mac
	lock       sync.Mutex
}

func NewArpScan() *ArpScan {
	return &ArpScan{
		macAddress: make(map[string]*mac),
	}
}

func (a *ArpScan) Run(interval time.Duration, life time.Duration) {
	t := time.NewTicker(interval)
	for {
		<-t.C

		mas, err := a.exec()
		if err != nil {
			log.Printf("failed to exec: %+v", err)
			continue
		}

		now := time.Now()

		nm := make(map[string]*mac)
		a.lock.Lock()
		for mac, v := range a.macAddress {
			if !v.LastTime.Add(life).After(now) {
				continue
			}
			nm[mac] = v
		}
		a.lock.Unlock()

		for _, m := range mas {
			if _, ok := nm[m]; ok {
				nm[m].LastTime = now
				continue
			}
			mc := &mac{
				Address:  m,
				FoundAt:  now,
				LastTime: now,
			}
			nm[m] = mc
			log.Printf("found mac %s", m)
		}

		a.lock.Lock()
		a.macAddress = nm
		a.lock.Unlock()
	}
}

func (a *ArpScan) exec() ([]string, error) {
	rb, err := exec.Command("sudo", "arp-scan", "-l", "-q", "-x").Output()
	if err != nil {
		return nil, xerrors.Errorf("failed to exec arp-scan: %w", err)
	}

	rs := string(rb)
	rss := strings.Split(rs, "\n")
	ms := make([]string, len(rss))
	for _, l := range rss {
		l = strings.Trim(l, "\r\n\t ")
		ls := strings.Split(l, "\t")
		m := strings.Trim(ls[len(ls)-1], "\t\r ")
		if m == "" {
			continue
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (a *ArpScan) Exist(mac string) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	v, ok := a.macAddress[mac]
	if !ok {
		return false
	}

	now := time.Now()
	if !v.FoundAt.Add(30 * time.Second).After(now) {
		return false
	}

	return ok
}
