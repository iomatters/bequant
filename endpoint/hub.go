package main

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/iomatters/db"
	"github.com/iomatters/provider"
)

type errorMsg struct {
	Code    int    `toml:"code"`
	Message string `toml:"message"`
}

type request struct {
	ws  *websocket.Conn
	msg message
}

type message struct {
	Fsyms string `toml:"fsyms"`
	Tsyms string `toml:"tsyms"`
}

type hub struct {
	provider string
	host     string
	port     int
	dbname   string
	user     string
	password string
	//clients  sync.Map[*websocket.Conn]bool
	clients sync.Map
	request chan request
	stop    chan os.Signal
}

func newHub(provider string, host string, port int, dbname string, user string, password string, ch chan os.Signal) *hub {
	return &hub{
		provider: provider,
		host:     host,
		port:     port,
		dbname:   dbname,
		user:     user,
		password: password,
		//clients:  make(map[*websocket.Conn]bool),
		clients: sync.Map{},
		request: make(chan request),
		stop:    ch,
	}
}

func (h *hub) run() {
	for {
		//var wg sync.WaitGroup
		select {
		// Terminate signal
		case <-h.stop:
			//for ws := range h.clients.Range() {
			//	ws.Close()
			//}
			h.clients.Range(func(key, value interface{}) bool {
				v := key.(*websocket.Conn)
				v.Close()
				//*websocket.Conn(key.(*websocket.Conn)).Close()
				return true
			})

			glog.Infoln("Bye!")
			os.Exit(0)

		case request := <-h.request:
			//wg.Add(1)
			go func(id string, fsyms string, tsyms string) {
				//defer wg.Done()

				p, err := provider.NewProvider(id)
				if err != nil {
					glog.Fatalln(err)
				}

				json_data, err := p.Pull([]string{fsyms}, []string{tsyms})
				if err != nil {
					glog.Errorln(err)

					conn, err := db.OpenConn(h.host, h.port, h.dbname, h.user, h.password)
					if err != nil {
						glog.Errorln(err)
						// TODO, respond 500 ?
						notFound(request.ws)
						return
					}
					defer conn.Close()

					// Fetch last
					json_str, err := conn.ReadLast(h.provider, fsyms, tsyms)
					if err != nil {
						glog.Errorln(err)
						notFound(request.ws)
						return
					}

					var row map[string]interface{}
					if err := json.Unmarshal([]byte(*json_str), &row); err != nil {
						glog.Errorln(err)
					}

					if err := request.ws.WriteJSON(minify(fsyms, tsyms, row)); err != nil {
						glog.Errorln(err)
					}
					glog.Infof("Fetched from database, FSYMS=%s, TSYMS=%s", fsyms, tsyms)

					return

				}

				// Got price
				if err := request.ws.WriteJSON(minify(fsyms, tsyms, json_data.Raw[fsyms][tsyms])); err != nil {
					glog.Errorln(err)
				}
				glog.Infof("Collected from API, FSYMS=%s, TSYMS=%s", fsyms, tsyms)

			}("cryptocompare", request.msg.Fsyms, request.msg.Tsyms)

		}
		//wg.Wait()
	}
}

func (h *hub) read(client *websocket.Conn) {
	for {
		var msg message
		if err := client.ReadJSON(&msg); err != nil {
			glog.Infoln(err)
			//delete(h.clients, client)
			h.clients.Delete(client)

			break
		}

		// Send a request to hub
		h.request <- request{
			ws:  client,
			msg: msg,
		}
	}
}

func notFound(ws *websocket.Conn) {
	if err := ws.WriteJSON(errorMsg{
		Code:    404,
		Message: "Not found",
	}); err != nil {
		glog.Errorln(err)
	}
}

// TODO - need some prettyfy
func minify(fsyms string, tsyms string, row map[string]interface{}) interface{} {
	resp := map[string]map[string]map[string]map[string]interface{}{
		"RAW": {
			fsyms: {
				tsyms: make(map[string]interface{}),
			}},
	}
	resp["RAW"][fsyms][tsyms]["CHANGE24HOUR"] = row["CHANGE24HOUR"]
	resp["RAW"][fsyms][tsyms]["CHANGEPCT24HOUR"] = row["CHANGEPCT24HOUR"]
	resp["RAW"][fsyms][tsyms]["OPEN24HOUR"] = row["OPEN24HOUR"]
	resp["RAW"][fsyms][tsyms]["VOLUME24HOUR"] = row["VOLUME24HOUR"]
	resp["RAW"][fsyms][tsyms]["VOLUME24HOURTO"] = row["VOLUME24HOURTO"]
	resp["RAW"][fsyms][tsyms]["LOW24HOUR"] = row["LOW24HOUR"]
	resp["RAW"][fsyms][tsyms]["HIGH24HOUR"] = row["HIGH24HOUR"]
	resp["RAW"][fsyms][tsyms]["PRICE"] = row["PRICE"]
	resp["RAW"][fsyms][tsyms]["LASTUPDATE"] = row["LASTUPDATE"]
	resp["RAW"][fsyms][tsyms]["SUPPLY"] = row["SUPPLY"]
	resp["RAW"][fsyms][tsyms]["MKTCAP"] = row["MKTCAP"]
	return resp
}
