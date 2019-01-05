package main

import (
	"flag"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Message is a wire message.
type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// Conn ...
type Conn struct {
	Conn *websocket.Conn
	Done chan struct{}
	*sync.Once
}

// Broker pairs clients together.
type Broker struct {
	conns chan Conn
}

// NewBroker creates a new Broker.
func NewBroker() *Broker {
	b := &Broker{
		conns: make(chan Conn),
	}
	go func() {
		for {
			x := <-b.conns
			log.WithField("x", x.Conn.RemoteAddr().String()).Info("pending pairing")
			y := <-b.conns
			go b.serve(x, y)
			go b.serve(y, x)
			log.WithFields(log.Fields{
				"x": x,
				"y": y,
			}).Info("paired connections")
		}
	}()
	return b
}

func (b *Broker) serve(x, y Conn) {
	for {
		_, msg, err := x.Conn.ReadMessage()
		if err != nil {
			log.WithError(err).WithField("conn", x).Error("could not decode next message")
			x.Do(func() { close(x.Done) })
			b.conns <- y // recycle y
			return
		}
		if err := y.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.WithError(err).WithField("conn", x).Error("could not encode next message")
			y.Do(func() { close(y.Done) })
			b.conns <- x // recycle x
			return
		}
	}
}

// Add adds a connection to the pool.
func (b *Broker) Add(ws *websocket.Conn) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		b.conns <- Conn{
			Conn: ws,
			Done: done,
			Once: new(sync.Once),
		}
	}()
	return done
}

func main() {
	var (
		addr = flag.String("addr", ":4998", "listen address")
		path = flag.String("path", "/cmd", "listen path")
	)
	flag.Parse()
	b := NewBroker()
	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}
	http.HandleFunc(*path, func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.WithError(err).Error("upgrade")
			return
		}
		done := b.Add(ws)
		<-done
	})
	log.Printf("handling %v on %v", *path, *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
