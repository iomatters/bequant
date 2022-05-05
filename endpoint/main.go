package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/iomatters/config"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/spf13/pflag"
)

var (
	BuildVersion = "(build)"
	Interval     = time.Second * 30
)

var upgrader = websocket.Upgrader{}

func init() {
	_ = flag.Set("logtostderr", "true")
}

func main() {
	filename := pflag.String("config", "", "path to config file")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	defer glog.Flush()

	params, err := config.NewAppConfig(*filename)
	if err != nil {
		glog.Fatalln(err)
	}

	glog.Infof("Starting endpoint on *:8080, version %s", BuildVersion)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	hub := newHub("cryptocompare", params.DB.Host, params.DB.Port, params.DB.DBName, params.DB.User, params.DB.Password, ch)
	// Run hub to collect client messages
	go hub.run()

	http.HandleFunc("/price", func(w http.ResponseWriter, r *http.Request) {
		glog.Infoln("Requested", r.URL)

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			glog.Fatal(err)
			return
		}
		defer func() {
			glog.Infoln("Closing socket", ws.RemoteAddr())
			//delete(hub.clients, ws)
			hub.clients.Delete(ws)
			ws.Close()
		}()

		// Add client
		//hub.clients[ws] = true
		hub.clients.Store(ws, true)

		glog.Infoln("Connected", ws.RemoteAddr())

		hub.read(ws)

	})

	http.ListenAndServe(":8080", nil)
}
