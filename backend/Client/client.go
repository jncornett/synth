package main

import (
	"flag"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		addr     = flag.String("addr", "localhost:4998", "server endpoint")
		path     = flag.String("path", "/cmd", "server path")
		name     = flag.String("name", "", "identifier")
		interval = flag.Duration("ival", 2*time.Second, "interval")
	)
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: *path}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	log.WithField("ws", ws.RemoteAddr()).Info("connected")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.WithError(err).Error("read")
				break
			}
			log.WithField("message", string(msg)).Info("read message")
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range time.Tick(*interval) {
			_ = t
			_ = *name
			err := ws.WriteMessage(websocket.TextMessage, []byte(`{"note":"C"}`))
			// err := ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s,%v", *name, t)))
			if err != nil {
				log.WithError(err).Error("write")
				break
			}
		}
	}()
	wg.Wait()
}
