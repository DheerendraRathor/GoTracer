package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"strings"

	"runtime"

	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/net/constants"
	"github.com/DheerendraRathor/GoTracer/tracer"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	pongTimeout = 60 * time.Second
	pongWait    = 20 * time.Second
	pingMessage = "a"
)

type RenderingClient struct {
	Conn                *websocket.Conn
	Results             chan *models.Pixel
	IsTracingInProgress bool
	CloseChan           chan bool
	OperationId         string
	PongChannel         chan bool
	WaitGroup           sync.WaitGroup
}

func (c *RenderingClient) Initialize() {
	c.WaitGroup = sync.WaitGroup{}

	c.Conn.SetPongHandler(
		func(message string) error {
			c.PongChannel <- true
			return nil
		},
	)

	go pingThread(c.Conn, c.PongChannel)

	c.WaitGroup.Add(2)

	go c.ReadHandler()
	go c.ResultSender()

	c.WaitGroup.Wait()
}

var addr = flag.String("addr", "", "http service address")

var upgrader = websocket.Upgrader{}
var mutex = &sync.Mutex{}

var activeConn *websocket.Conn = nil

func sendPing(conn *websocket.Conn) error {
	err := conn.WriteControl(websocket.PingMessage, []byte(pingMessage), time.Now().Add(time.Hour))
	return err
}

func pingThread(conn *websocket.Conn, pongChan <-chan bool) {
	ticker := time.NewTicker(pongWait)
	timer := time.NewTimer(pongTimeout)

	for {
		select {
		case <-timer.C:
			conn.WriteControl(websocket.CloseNoStatusReceived, []byte("No pong received"), time.Now().Add(time.Hour))
			conn.Close()
			return
		case <-ticker.C:
			err := sendPing(conn)
			if err != nil {
				ticker.Stop()
				return
			}
		case <-pongChan:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(pongTimeout)
		}
	}
}

func (c *RenderingClient) ReadHandler() {
	closeChan := make(chan bool)

	defer func() {
		c.WaitGroup.Done()
		closeChan <- true
	}()

	var message messages.WebSocketMessage
	var renderReqMsg messages.RenderRequestMessage
	for {
		_, msgStr, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Unable to understand message: %s", err)
			break
		}

		json.Unmarshal(msgStr, &message)

		if message.Type == messages.RenderRequest {
			json.Unmarshal(msgStr, &renderReqMsg)
			operationId := renderReqMsg.OperationId
			if strings.TrimSpace(operationId) == "" {
				byteOpId, _ := uuid.NewRandom()
				operationId = byteOpId.String()
			}
			c.OperationId = operationId

			responseMessage := messages.RenderRequestResponseMessage{
				Type:        messages.RenderRequestResponse,
				OperationId: c.OperationId,
			}

			if !c.IsTracingInProgress {
				go goTracer.GoTrace(&renderReqMsg.Data, c.Results, closeChan)
				c.IsTracingInProgress = true
				responseMessage.Code = messages.RenderRequestAccepted
			} else {
				responseMessage.Code = messages.RenderRequestTracingAlreadyInProgress
			}

			c.Conn.WriteJSON(message)
		}
	}
}

func (c *RenderingClient) ResultSender() {
	defer c.WaitGroup.Done()

	message := messages.WebSocketMessage{
		Type: messages.PixelResult,
	}

	for pixel := range c.Results {
		message.OperationId = c.OperationId
		if pixel == nil {
			message.Type = messages.RenderingCompleted
			c.Conn.WriteJSON(message)
			break
		}

		message.Data = pixel
		c.Conn.WriteJSON(message)
	}
}

func agentHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a connection request")
	shouldConnect := true

	mutex.Lock()
	if activeConn != nil {
		log.Print("An active connection already exists")
		http.Error(w, "An active connection already exists", http.StatusConflict)
		shouldConnect = false
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		shouldConnect = false
	}
	activeConn = conn
	mutex.Unlock()

	defer func() {
		activeConn = nil
		conn.Close()
	}()

	if !shouldConnect {
		return
	}

	client := RenderingClient{
		Conn:                conn,
		PongChannel:         make(chan bool),
		Results:             make(chan *models.Pixel, 100),
		IsTracingInProgress: false,
		OperationId:         "",
	}

	client.Initialize()
}

func agentStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := messages.AgentStatus{
		Available: activeConn == nil,
		Cores:     runtime.NumCPU(),
	}

	jsonStatus, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonStatus)
}

func main() {
	flag.Parse()

	if *addr == "" {
		flag.Usage()
		log.Fatalf("Agent address \"%s\" is invalid.", *addr)
	}

	log.SetFlags(0)
	http.HandleFunc("/", agentHandler)
	http.HandleFunc("/status", agentStatusHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
