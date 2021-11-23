package arpscan

import (
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/moezakura/auto-unlock/pkg/times"
	"golang.org/x/xerrors"
)

type ArpScan struct {
	macAddress map[string]*mac
	lock       sync.Mutex
	isVerbose  bool
}

func NewArpScan(isVerbose bool) *ArpScan {
	return &ArpScan{
		macAddress: make(map[string]*mac),
		isVerbose:  isVerbose,
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
			nm[mac] = v
		}
		a.lock.Unlock()

		for _, m := range mas {
			if strings.TrimSpace(m) == "" {
				continue
			}

			if _, ok := nm[m]; ok {
				if a.isVerbose {
					log.Printf("Updated mac address: %s (%s -> %s)", m,
						times.ToYMDHIS(nm[m].LastTime), times.ToYMDHIS(now))
				}

				if nm[m].LastTime.Add(2 * time.Minute).Before(now) {
					if a.isVerbose {
						log.Printf("Refound mac address: %s (found: %s, refound: %s)", m,
							times.ToYMDHIS(nm[m].FoundAt), times.ToYMDHIS(now))
					}
					nm[m].FoundAt = now
				}

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
	cmd := exec.Command("sudo", "arp-scan", "-lqx")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	defer func() {
		_ = stdout.Close()
		_ = stderr.Close()
	}()

	if err != nil {
		o, _ := io.ReadAll(stdout)
		e, _ := io.ReadAll(stderr)

		return nil, xerrors.Errorf("failed to exec arp-scan: res: stdout: %s, stderr: %s,  err:%w",
			string(o), string(e), err)
	}

	o, _ := io.ReadAll(stdout)
	rs := string(o)
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
	if v.FoundAt.Add(30 * time.Second).Before(now) {
		log.Printf("%s is too old: %s sec", mac, now.Sub(v.FoundAt).String())
		return false
	}

	return ok
}
