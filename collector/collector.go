package main

import (
	"os"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/iomatters/db"
	"github.com/iomatters/provider"
)

const (
	Interval = 10
)

type collector struct {
	provider string
	host     string
	port     int
	dbname   string
	user     string
	password string
	fsyms    []string
	tsyms    []string
	stop     chan os.Signal
}

func newCollector(provider string, host string, port int, dbname string, user string, password string, fsyms []string, tsyms []string, ch chan os.Signal) *collector {
	return &collector{
		provider: provider,
		host:     host,
		port:     port,
		dbname:   dbname,
		user:     user,
		password: password,
		fsyms:    fsyms,
		tsyms:    tsyms,
		stop:     ch,
	}
}

func (c *collector) run() {
	var wg sync.WaitGroup

	for {
		glog.Infof("Next cycle in [%d] seconds", Interval)
		select {
		case <-c.stop:
			glog.Infof("Bye!")
			return
		case <-time.After(time.Second * time.Duration(Interval)):
			wg.Add(1)
			go func(id string, fsyms []string, tsyms []string) {
				defer wg.Done()
				p, err := provider.NewProvider(id)
				if err != nil {
					glog.Fatalln(err)
				}
				json_data, err := p.Pull(fsyms, tsyms)
				if err != nil {
					glog.Error(err)
					return
				}
				glog.Infof("Collected FSYMS: %s, TSYMS: %s", fsyms, tsyms)

				conn, err := db.OpenConn(c.host, c.port, c.dbname, c.user, c.password)
				if err != nil {
					glog.Error(err)
					return
				}
				defer conn.Close()

				if err := conn.Write(id, json_data); err != nil {
					glog.Error(err)
					return
				}
				glog.Infof("Recorded FSYMS: %s, TSYMS: %s", fsyms, tsyms)

			}(c.provider, c.fsyms, c.tsyms)

		}
		wg.Wait()
	}
}
