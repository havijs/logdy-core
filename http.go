package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type MessageType int

const MessageTypeStdout MessageType = 1
const MessageTypeStderr MessageType = 2

type Message struct {
	Mtype       MessageType     `json:"message_type"`
	Content     string          `json:"content"`
	JsonContent json.RawMessage `json:"json_content"`
	IsJson      bool            `json:"is_json"`
	Ts          time.Time       `json:"ts"`
}

type Clients struct {
	mu                 sync.Mutex
	mainChan           <-chan Message
	clients            map[int]chan Message
	buffer             []Message
	currentlyConnected int
}

func (c *Clients) Start() {
	for {
		msg := <-c.mainChan
		c.mu.Lock()
		if c.currentlyConnected == 0 {
			logger.Debug("Received a log message but no client is connected, buffering message")
			c.buffer = append(c.buffer, msg)
		}

		for _, ch := range c.clients {

			ch <- msg
		}
		c.mu.Unlock()
	}
}

func (c *Clients) Join(id int) <-chan Message {
	c.mu.Lock()
	defer func() {
		logger.WithFields(logrus.Fields{
			"msg_count": len(c.buffer),
		}).Info("Flushing log messages buffer to a recently connected client")
		for _, msg := range c.buffer {
			c.clients[id] <- msg
		}

		c.buffer = []Message{}
	}()
	defer c.mu.Unlock()

	c.clients[id] = make(chan Message, 100)
	c.currentlyConnected++
	return c.clients[id]
}

func (c *Clients) Close(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.clients, id)
	c.currentlyConnected--
}

func handleHttp(msgs <-chan Message, httpPort string) {
	// Create a new WebSocket server.
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients := Clients{
		mu:                 sync.Mutex{},
		mainChan:           msgs,
		clients:            map[int]chan Message{},
		currentlyConnected: 0,
		buffer:             []Message{},
	}

	go clients.Start()

	cid := 0

	assets, _ := Assets()

	// Use the file system to serve static files
	fs := http.FileServer(http.FS(assets))
	http.Handle("/", http.StripPrefix("/", fs))

	// Listen for WebSocket connections on port 8080.
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to a WebSocket connection.
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		logger.Info("New Web UI client connected")

		cid++
		clientId := cid
		ch := clients.Join(cid)

		for {
			msg := <-ch
			bts, err := json.Marshal(msg)

			logger.WithFields(logrus.Fields{
				"msg": string(bts),
			}).Debug("Sending message through WebSocket")

			if err != nil {
				fmt.Printf("Received message %+v", msg)
				fmt.Println("Error while serializing message", err)
				continue
			}

			err = conn.WriteMessage(1, bts)

			if err != nil {
				logger.Error("Err", err)
				clients.Close(clientId)
				logger.Info("Closed client", clientId)
				break
			}
		}

	})

	logger.WithFields(logrus.Fields{
		"port": httpPort,
	}).Info("WebUI started, visit http://localhost:" + httpPort)

	http.ListenAndServe(":"+httpPort, nil)
}
